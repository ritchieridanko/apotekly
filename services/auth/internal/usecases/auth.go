package usecases

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	"github.com/ritchieridanko/apotekly/services/auth/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/ritchieridanko/apotekly/services/shared/contract/events/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/infra/publisher"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/bcrypt"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthUsecase interface {
	SignUp(ctx context.Context, req *models.SignUpReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	SignIn(ctx context.Context, req *models.SignInReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	SignOut(ctx context.Context, refreshToken string) (err *ce.Error)
	IsEmailAvailable(ctx context.Context, email string) (available bool, err *ce.Error)
	RotateAuthToken(ctx context.Context, refreshToken string) (at *models.AuthToken, err *ce.Error)
	ResendVerification(ctx context.Context) (email string, err *ce.Error)
	VerifyEmail(ctx context.Context, req *models.VerifyEmailReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	ChangeEmail(ctx context.Context, req *models.ChangeEmailReq) (email string, err *ce.Error)
	ChangePassword(ctx context.Context, req *models.ChangePasswordReq) (err *ce.Error)
}

type authUsecase struct {
	appName           string
	emailChangeToken  time.Duration
	verificationToken time.Duration
	su                SessionUsecase
	ar                repositories.AuthRepository
	tr                repositories.TokenRepository
	transactor        *database.Transactor
	acp               *publisher.Publisher
	aecrp             *publisher.Publisher
	aevrp             *publisher.Publisher
	bcrypt            *bcrypt.BCrypt
	validator         *validator.Validator
	logger            *logger.Logger
}

func NewAuthUsecase(
	appName string,
	emailChangeToken time.Duration,
	verificationToken time.Duration,
	su SessionUsecase,
	ar repositories.AuthRepository,
	tr repositories.TokenRepository,
	tx *database.Transactor,
	acp *publisher.Publisher,
	aecrp *publisher.Publisher,
	aevrp *publisher.Publisher,
	b *bcrypt.BCrypt,
	v *validator.Validator,
	l *logger.Logger,
) AuthUsecase {
	return &authUsecase{
		appName:           appName,
		emailChangeToken:  emailChangeToken,
		verificationToken: verificationToken,
		su:                su,
		ar:                ar,
		tr:                tr,
		transactor:        tx,
		acp:               acp,
		aecrp:             aecrp,
		aevrp:             aevrp,
		bcrypt:            b,
		validator:         v,
		logger:            l,
	}
}

func (u *authUsecase) SignUp(ctx context.Context, req *models.SignUpReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.SignUp")
	defer span.End()

	// Data Normalization
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Data Validation
	if ok, why := u.validator.Email(email); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}
	if ok, why := u.validator.Password(req.Password); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}

	var a *models.Auth
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Email Availability Check
		available, err := u.ar.IsEmailAvailable(ctx, email)
		if err != nil {
			return err
		}
		if !available {
			return ce.NewError(ce.CodeEmailNotAvailable, ce.MsgEmailAlreadyRegistered, nil)
		}

		// Password Hashing
		hash, hashErr := u.bcrypt.Hash(req.Password)
		if hashErr != nil {
			return ce.NewError(ce.CodeBCryptHashingFailed, ce.MsgInternalServer, hashErr)
		}

		// Auth Creation
		a, err = u.ar.Create(
			ctx,
			&models.CreateAuth{
				Email:    email,
				Password: &hash,
				Role:     constants.RoleUser,
			},
		)
		return err
	})
	if err != nil {
		return nil, nil, err
	}

	authIDField := logger.NewField("auth_id", a.ID)

	// Session Creation
	// NOTE: Fail to create session does not fail SignUp usecase
	var session *string
	at, err := u.su.CreateSession(
		ctx,
		&models.CreateSessionReq{
			AuthID:          a.ID,
			Role:            a.Role,
			IsEmailVerified: a.IsEmailVerified(),
		},
	)
	if err != nil {
		u.logger.Warn(
			ctx,
			"created auth. failed to create session",
			err.Append(
				logger.NewField("error_code", err.Code()),
				logger.NewField("error", err.Unwrap()),
			).Fields()...,
		)
	} else {
		session = &at.RefreshToken.Token
	}

	// Verification Token Creation
	// NOTE: Fail to create verification token does not fail SignUp usecase
	token := utils.GenerateUUID().String()
	verificationToken := &token
	err = u.tr.CreateVerification(
		ctx,
		&models.CreateVerificationToken{
			AuthID:   a.ID,
			Token:    token,
			Duration: u.verificationToken,
		},
	)
	if err != nil {
		verificationToken = nil
		u.logger.Warn(
			ctx,
			"created auth. failed to create verification token",
			err.Append(
				authIDField,
				logger.NewField("error_code", err.Code()),
				logger.NewField("error", err.Unwrap()),
			).Fields()...,
		)
	}

	evtTopicField := logger.NewField("event_topic", u.acp.Topic())

	// auth.created Event Publishing
	// NOTE: Fail to publish event does not fail SignUp usecase
	pubErr := u.acp.Publish(
		ctx,
		"auth_"+strconv.FormatUint(a.ID, 10),
		&events.AuthCreated{
			Id:                utils.GenerateUUID().String(),
			AuthId:            a.ID,
			Email:             a.Email,
			Role:              a.Role,
			EmailVerifiedAt:   utils.ToTimestamp(a.EmailVerifiedAt),
			Session:           session,
			VerificationToken: verificationToken,
			CreatedAt:         timestamppb.New(time.Now().UTC()),
		},
	)
	if pubErr != nil {
		u.logger.Warn(
			ctx,
			"created auth. failed to publish event",
			authIDField,
			evtTopicField,
			logger.NewField("error_code", ce.CodeEventPublishingFailed),
			logger.NewField("error", pubErr),
		)
		return a, at, nil
	}

	u.logger.Info(
		ctx,
		"EVENT PUBLISHED",
		authIDField,
		evtTopicField,
	)

	return a, at, nil
}

func (u *authUsecase) SignIn(ctx context.Context, req *models.SignInReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.SignIn")
	defer span.End()

	// Data Normalization
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Data Validation
	if ok, why := u.validator.Email(email); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}

	// Auth Fetching
	a, err := u.ar.GetByEmail(ctx, email)
	if err != nil && err.Code() == ce.CodeAuthNotFound {
		return nil, nil, ce.NewError(
			ce.CodeEmailNotRegistered,
			ce.MsgInvalidCredentials,
			err.Unwrap(),
		)
	}
	if err != nil {
		return nil, nil, err
	}

	authIDField := logger.NewField("auth_id", a.ID)

	// Password Validation
	if a.Password == nil {
		return nil, nil, ce.NewError(
			ce.CodeOAuthRegularSignIn,
			ce.MsgInvalidCredentials,
			errors.New("auth has no password"),
			authIDField,
		)
	}
	if err := u.bcrypt.Validate(*a.Password, req.Password); err != nil {
		return nil, nil, ce.NewError(
			ce.CodeWrongPassword,
			ce.MsgInvalidCredentials,
			err,
			authIDField,
		)
	}

	// Session Creation
	at, err := u.su.CreateSession(
		ctx,
		&models.CreateSessionReq{
			AuthID:          a.ID,
			Role:            a.Role,
			IsEmailVerified: a.IsEmailVerified(),
		},
	)

	return a, at, err
}

func (u *authUsecase) SignOut(ctx context.Context, refreshToken string) *ce.Error {
	// Data Normalization
	token := strings.TrimSpace(refreshToken)

	// Data Validation
	// NOTE: Empty token does not fail SignOut usecase
	if token == "" {
		return nil
	}

	// Session Revocation
	// NOTE: Invalid session does not fail SignOut usecase
	if err := u.su.RevokeSession(ctx, token); err != nil && err.Code() != ce.CodeSessionNotFound {
		return err
	}
	return nil
}

func (u *authUsecase) IsEmailAvailable(ctx context.Context, email string) (bool, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.IsEmailAvailable")
	defer span.End()

	// Data Normalization
	em := strings.ToLower(strings.TrimSpace(email))

	// Data Validation
	if ok, why := u.validator.Email(em); !ok {
		return false, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}

	// Email Availability Check
	return u.ar.IsEmailAvailable(ctx, em)
}

func (u *authUsecase) RotateAuthToken(ctx context.Context, refreshToken string) (*models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.RotateAuthToken")
	defer span.End()

	// Data Normalization
	token := strings.TrimSpace(refreshToken)

	// Data Validation
	if token == "" {
		return nil, ce.NewError(
			ce.CodeUnauthenticated,
			ce.MsgUnauthenticated,
			errors.New("refresh token is empty"),
		)
	}

	var at *models.AuthToken
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Session Fetching
		s, err := u.su.GetSession(ctx, token)
		if err != nil {
			return err
		}

		authIDField := logger.NewField("auth_id", s.AuthID)

		// Session Expiration Check
		if s.ExpiresAt.Before(time.Now().UTC()) {
			return ce.NewError(
				ce.CodeSessionExpired,
				ce.MsgSessionExpired,
				nil,
				authIDField,
			)
		}

		// Auth Fetching
		a, err := u.ar.GetByID(ctx, s.AuthID)
		if err != nil && err.Code() == ce.CodeAuthNotFound {
			return ce.NewError(
				ce.CodeAuthNotRegistered,
				ce.MsgInvalidCredentials,
				err.Unwrap(),
				authIDField,
			)
		}
		if err != nil {
			return err.Append(authIDField)
		}

		// Session Refresh
		at, err = u.su.RefreshSession(
			ctx,
			&models.RefreshSessionReq{
				AuthID:          a.ID,
				Role:            a.Role,
				IsEmailVerified: a.IsEmailVerified(),
				RefreshToken:    token,
			},
		)
		return err
	})

	return at, err
}

func (u *authUsecase) ResendVerification(ctx context.Context) (string, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.ResendVerification")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return "", ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)

	// Auth Fetching
	// NOTE: ResendVerification usecase is only allowed if not yet verified
	a, err := u.ar.GetByID(ctx, authCtx.AuthID)
	if err != nil && err.Code() == ce.CodeAuthNotFound {
		return "", ce.NewError(
			ce.CodeAuthNotRegistered,
			ce.MsgInvalidCredentials,
			err.Unwrap(),
			authIDField,
		)
	}
	if err != nil {
		return "", err.Append(authIDField)
	}
	if a.IsEmailVerified() {
		return "", ce.NewError(
			ce.CodeEmailAlreadyVerified,
			ce.MsgEmailAlreadyVerified,
			nil,
			authIDField,
		)
	}

	// Verification Token Creation
	token := utils.GenerateUUID().String()
	err = u.tr.CreateVerification(
		ctx,
		&models.CreateVerificationToken{
			AuthID:   a.ID,
			Token:    token,
			Duration: u.verificationToken,
		},
	)
	if err != nil {
		return "", err.Append(authIDField)
	}

	evtTopicField := logger.NewField("event_topic", u.aevrp.Topic())

	// auth.email.verification.requested Event Publishing
	id := utils.GenerateUUID().String()
	pubErr := u.aevrp.Publish(
		ctx,
		"auth_"+strconv.FormatUint(a.ID, 10)+"_"+id,
		&events.AuthEmailVerificationRequested{
			Id:        id,
			AuthId:    a.ID,
			Email:     a.Email,
			Role:      a.Role,
			Token:     token,
			CreatedAt: timestamppb.New(time.Now().UTC()),
		},
	)
	if pubErr != nil {
		return "", ce.NewError(
			ce.CodeEventPublishingFailed,
			ce.MsgInternalServer,
			pubErr,
			authIDField,
			evtTopicField,
		)
	}

	u.logger.Info(
		ctx,
		"EVENT PUBLISHED",
		authIDField,
		evtTopicField,
	)

	return a.Email, nil
}

func (u *authUsecase) VerifyEmail(ctx context.Context, req *models.VerifyEmailReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.VerifyEmail")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)

	// Data Normalization
	refreshToken := strings.TrimSpace(req.RefreshToken)
	verificationToken := strings.TrimSpace(req.VerificationToken)

	// Data Validation
	if refreshToken == "" {
		return nil, nil, ce.NewError(
			ce.CodeUnauthenticated,
			ce.MsgUnauthenticated,
			errors.New("refresh token is empty"),
			authIDField,
		)
	}
	if verificationToken == "" {
		return nil, nil, ce.NewError(
			ce.CodeInvalidPayload,
			"Verification token is required",
			nil,
			authIDField,
		)
	}

	// Verification Token Consumption
	authID, err := u.tr.UseVerification(ctx, verificationToken)
	if err != nil {
		return nil, nil, err.Append(authIDField)
	}

	// Verification Token Ownership Validation
	// NOTE: Invalid ownership re-creates the verification token
	if authID != authCtx.AuthID {
		err := u.tr.CreateVerification(
			ctx,
			&models.CreateVerificationToken{
				AuthID:   authID,
				Token:    verificationToken,
				Duration: u.verificationToken,
			},
		)
		if err != nil {
			return nil, nil, err.Append(authIDField)
		}
		return nil, nil, ce.NewError(
			ce.CodeTokenNotOwned,
			ce.MsgInvalidToken,
			nil,
			authIDField,
			logger.NewField("token_auth_id", authID),
		)
	}

	// Verification Update
	// NOTE: Fail to update verification status re-creates the verification token
	a, err := u.ar.SetVerified(ctx, authID)
	if err != nil && err.Code() == ce.CodeAuthNotFound {
		return nil, nil, ce.NewError(
			ce.CodeAuthNotRegistered,
			ce.MsgInvalidCredentials,
			err.Unwrap(),
			authIDField,
		)
	}
	if err != nil {
		createErr := u.tr.CreateVerification(
			ctx,
			&models.CreateVerificationToken{
				AuthID:   authID,
				Token:    verificationToken,
				Duration: u.verificationToken,
			},
		)
		if createErr != nil {
			return nil, nil, createErr.Append(authIDField)
		}
		return nil, nil, err.Append(authIDField)
	}

	// Session Refresh
	// NOTE: Fail to refresh session does not fail VerifyEmail usecase
	at, err := u.su.RefreshSession(
		ctx,
		&models.RefreshSessionReq{
			AuthID:          a.ID,
			Role:            a.Role,
			IsEmailVerified: a.IsEmailVerified(),
			RefreshToken:    refreshToken,
		},
	)
	if err != nil {
		u.logger.Warn(
			ctx,
			"verified email. failed to refresh session",
			err.Append(
				authIDField,
				logger.NewField("error_code", err.Code()),
				logger.NewField("error", err.Unwrap()),
			).Fields()...,
		)
		return a, nil, nil
	}

	return a, at, nil
}

func (u *authUsecase) ChangeEmail(ctx context.Context, req *models.ChangeEmailReq) (string, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.ChangeEmail")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return "", ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)

	// Data Normalization
	newEmail := strings.ToLower(req.NewEmail)

	// Data Validation
	if ok, why := u.validator.Email(newEmail); !ok {
		return "", ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}

	// Auth Fetching
	// NOTE: ChangeEmail usecase is not allowed for oauth account (password == nil)
	a, err := u.ar.GetByID(ctx, authCtx.AuthID)
	if err != nil && err.Code() == ce.CodeAuthNotFound {
		return "", ce.NewError(
			ce.CodeAuthNotRegistered,
			ce.MsgInvalidCredentials,
			err.Unwrap(),
			authIDField,
		)
	}
	if err != nil {
		return "", err.Append(authIDField)
	}
	if a.Password == nil {
		return "", ce.NewError(
			ce.CodeOAuthEmailChange,
			ce.MsgOAuthEmailChange,
			nil,
			authIDField,
		)
	}

	// Password Validation
	if err := u.bcrypt.Validate(*a.Password, req.Password); err != nil {
		return "", ce.NewError(
			ce.CodeWrongPassword,
			ce.MsgInvalidPassword,
			err,
			authIDField,
		)
	}

	// Email Availability Check
	available, err := u.ar.IsEmailAvailable(ctx, newEmail)
	if err != nil {
		return "", err.Append(authIDField)
	}
	if !available {
		return "", ce.NewError(
			ce.CodeEmailNotAvailable,
			ce.MsgEmailAlreadyRegistered,
			nil,
			authIDField,
		)
	}

	// Email Change Token Creation
	token := utils.GenerateUUID().String()
	err = u.tr.CreateEmailChange(
		ctx,
		&models.CreateEmailChange{
			AuthID:   a.ID,
			NewEmail: newEmail,
			Token:    token,
			Duration: u.emailChangeToken,
		},
	)
	if err != nil {
		return "", err.Append(authIDField)
	}

	evtTopicField := logger.NewField("event_topic", u.aecrp.Topic())

	// auth.email.change.requested Event Publishing
	// NOTE: Fail to publish event cancels email reservation
	id := utils.GenerateUUID().String()
	pubErr := u.aecrp.Publish(
		ctx,
		"auth_"+strconv.FormatUint(a.ID, 10)+"_"+id,
		&events.AuthEmailChangeRequested{
			Id:        id,
			AuthId:    a.ID,
			OldEmail:  a.Email,
			NewEmail:  newEmail,
			Role:      a.Role,
			Token:     token,
			CreatedAt: timestamppb.New(time.Now().UTC()),
		},
	)
	if pubErr != nil {
		if err := u.ar.UnreserveEmail(ctx, newEmail); err != nil {
			return "", err.Append(authIDField, evtTopicField)
		}
		return "", ce.NewError(
			ce.CodeEventPublishingFailed,
			ce.MsgInternalServer,
			pubErr,
			authIDField,
			evtTopicField,
		)
	}

	u.logger.Info(
		ctx,
		"EVENT PUBLISHED",
		authIDField,
		evtTopicField,
	)

	return newEmail, nil
}

func (u *authUsecase) ChangePassword(ctx context.Context, req *models.ChangePasswordReq) *ce.Error {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.ChangePassword")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)

	// Data Validation
	if ok, why := u.validator.Password(req.NewPassword); !ok {
		return ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}

	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Auth Fetching
		// NOTE: ChangePassword usecase is not allowed for oauth account (password == nil)
		a, err := u.ar.GetByID(ctx, authCtx.AuthID)
		if err != nil && err.Code() == ce.CodeAuthNotFound {
			return ce.NewError(ce.CodeAuthNotRegistered, ce.MsgInvalidCredentials, err.Unwrap())
		}
		if err != nil {
			return err
		}
		if a.Password == nil {
			return ce.NewError(ce.CodeOAuthPasswordChange, ce.MsgOAuthPasswordChange, nil)
		}

		// Old Password Validation
		if err := u.bcrypt.Validate(*a.Password, req.OldPassword); err != nil {
			return ce.NewError(ce.CodeWrongPassword, ce.MsgInvalidOldPassword, err)
		}

		// New Password Hashing
		hash, hashErr := u.bcrypt.Hash(req.NewPassword)
		if hashErr != nil {
			return ce.NewError(ce.CodeBCryptHashingFailed, ce.MsgInternalServer, hashErr)
		}

		// Password Update
		err = u.ar.UpdatePassword(ctx, a.ID, hash)
		if err != nil && err.Code() == ce.CodeAuthNotFound {
			return ce.NewError(ce.CodeAuthNotRegistered, ce.MsgInvalidCredentials, err.Unwrap())
		}
		return err
	})
	if err != nil {
		return err.Append(authIDField)
	}

	return nil
}
