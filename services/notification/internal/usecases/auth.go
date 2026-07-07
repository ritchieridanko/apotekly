package usecases

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/notification/internal/channels"
	"github.com/ritchieridanko/apotekly/services/notification/internal/models"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"go.opentelemetry.io/otel"
)

type AuthUsecase interface {
	ProcessEventAC(ctx context.Context, data *models.EventAC) (err *ce.Error)
	ProcessEventAECR(ctx context.Context, data *models.EventAECR) (err *ce.Error)
	ProcessEventAEVR(ctx context.Context, data *models.EventAEVR) (err *ce.Error)
	ProcessEventAPRR(ctx context.Context, data *models.EventAPRR) (err *ce.Error)
}

type authUsecase struct {
	appName string
	ec      channels.EmailChannel
}

func NewAuthUsecase(appName string, ec channels.EmailChannel) AuthUsecase {
	return &authUsecase{appName: appName, ec: ec}
}

func (u *authUsecase) ProcessEventAC(ctx context.Context, data *models.EventAC) *ce.Error {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.ProcessEventAC")
	defer span.End()

	var token string
	if data.VerificationToken != nil {
		token = *data.VerificationToken
	}

	return u.ec.SendWelcome(
		ctx,
		&models.WelcomeEmail{
			Recipient:         data.Email,
			Role:              data.Role,
			IsEmailVerified:   data.EmailVerifiedAt != nil,
			VerificationToken: token,
		},
	)
}

func (u *authUsecase) ProcessEventAECR(ctx context.Context, data *models.EventAECR) *ce.Error {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.ProcessEventAECR")
	defer span.End()

	return u.ec.SendEmailChange(
		ctx,
		&models.EmailChangeEmail{
			EventID:  data.ID,
			OldEmail: data.OldEmail,
			NewEmail: data.NewEmail,
			Role:     data.Role,
			Token:    data.Token,
		},
	)
}

func (u *authUsecase) ProcessEventAEVR(ctx context.Context, data *models.EventAEVR) *ce.Error {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.ProcessEventAEVR")
	defer span.End()

	return u.ec.SendVerification(
		ctx,
		&models.VerificationEmail{
			Recipient: data.Email,
			Role:      data.Role,
			Token:     data.Token,
		},
	)
}

func (u *authUsecase) ProcessEventAPRR(ctx context.Context, data *models.EventAPRR) *ce.Error {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.ProcessEventAPRR")
	defer span.End()

	return u.ec.SendPasswordReset(
		ctx,
		&models.PasswordResetEmail{
			Recipient: data.Email,
			Role:      data.Role,
			Token:     data.Token,
		},
	)
}
