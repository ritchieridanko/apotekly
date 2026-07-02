package usecases

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/notification/internal/channels"
	"github.com/ritchieridanko/apotekly/services/notification/internal/models"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"go.opentelemetry.io/otel"
)

type AuthUsecase interface {
	HandleAuthCreated(ctx context.Context, data *models.AuthCreatedEvt) (err *ce.Error)
}

type authUsecase struct {
	appName string
	ec      channels.EmailChannel
}

func NewAuthUsecase(appName string, ec channels.EmailChannel) AuthUsecase {
	return &authUsecase{appName: appName, ec: ec}
}

func (u *authUsecase) HandleAuthCreated(ctx context.Context, data *models.AuthCreatedEvt) *ce.Error {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.HandleAuthCreated")
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
