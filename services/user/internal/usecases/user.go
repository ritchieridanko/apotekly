package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
	"github.com/ritchieridanko/apotekly/services/user/internal/models"
	"github.com/ritchieridanko/apotekly/services/user/internal/repositories"
	"go.opentelemetry.io/otel"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req *models.CreateUserReq) (u *models.User, err *ce.Error)
	GetMe(ctx context.Context) (u *models.User, err *ce.Error)
	UpdateUser(ctx context.Context, req *models.UpdateUserReq) (u *models.User, err *ce.Error)
}

type userUsecase struct {
	appName   string
	ur        repositories.UserRepository
	validator *validator.Validator
	logger    *logger.Logger
}

func NewUserUsecase(appName string, ur repositories.UserRepository, v *validator.Validator, l *logger.Logger) UserUsecase {
	return &userUsecase{
		appName:   appName,
		ur:        ur,
		validator: v,
		logger:    l,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, req *models.CreateUserReq) (*models.User, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "user.usecase.CreateUser")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)

	// Data Normalization
	name := strings.TrimSpace(req.Name)
	sex := utils.ToLowerPtr(utils.TrimSpacePtr(req.Sex))
	phone := utils.TrimSpacePtr(req.Phone)

	// Data Validation
	if ok, why := u.validator.Name(name); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
	}
	if sex != nil {
		if ok, why := u.validator.Sex(*sex); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if req.Birthdate != nil {
		if ok, why := u.validator.Birthdate(*req.Birthdate); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if phone != nil {
		if ok, why := u.validator.Phone(*phone); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}

	// User Creation
	user, err := u.ur.Create(
		ctx,
		&models.CreateUser{
			ID:        utils.MustGenerateUUIDv7(),
			AuthID:    authCtx.AuthID,
			Name:      name,
			Sex:       sex,
			Birthdate: req.Birthdate,
			Phone:     phone,
		},
	)
	if err != nil {
		return nil, err.Append(authIDField)
	}

	return user, nil
}

func (u *userUsecase) GetMe(ctx context.Context) (*models.User, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "user.usecase.GetMe")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)

	// User Fetching
	user, err := u.ur.GetByAuthID(ctx, authCtx.AuthID)
	if err != nil {
		return nil, err.Append(authIDField)
	}
	return user, nil
}

func (u *userUsecase) UpdateUser(ctx context.Context, req *models.UpdateUserReq) (*models.User, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "user.usecase.UpdateUser")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)

	// Data Normalization
	name := utils.TrimSpacePtr(req.Name)
	sex := utils.ToLowerPtr(utils.TrimSpacePtr(req.Sex))
	phone := utils.TrimSpacePtr(req.Phone)

	// Data Validation
	if name != nil {
		if ok, why := u.validator.Name(*name); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if sex != nil {
		if ok, why := u.validator.Sex(*sex); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if req.Birthdate != nil {
		if ok, why := u.validator.Birthdate(*req.Birthdate); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}
	if phone != nil {
		if ok, why := u.validator.Phone(*phone); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authIDField)
		}
	}

	// User Update
	user, err := u.ur.Update(
		ctx,
		authCtx.AuthID,
		&models.UpdateUser{
			Name:      name,
			Sex:       sex,
			Birthdate: req.Birthdate,
			Phone:     phone,
		},
	)
	if err != nil {
		return nil, err.Append(authIDField)
	}

	return user, nil
}
