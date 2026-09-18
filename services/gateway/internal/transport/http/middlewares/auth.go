package middlewares

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/ritchieridanko/apotekly/services/shared/utils/jwt"
)

func Auth(j *jwt.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := strings.TrimSpace(ctx.GetHeader("Authorization"))
		if len(authorization) == 0 {
			ce.NewError(
				ce.CodeUnauthenticated,
				ce.MsgUnauthenticated,
				errors.New("access token missing from request header"),
			).Bind(
				ctx,
			)

			ctx.Abort()
			return
		}

		auth := strings.Split(authorization, " ")
		if len(auth) != 2 || strings.ToLower(auth[0]) != "bearer" {
			ce.NewError(
				ce.CodeUnauthenticated,
				ce.MsgUnauthenticated,
				errors.New("access token is malformed"),
				logger.NewField("access_token", authorization),
			).Bind(
				ctx,
			)

			ctx.Abort()
			return
		}

		claim, err := j.Parse(auth[1])
		if err != nil {
			switch {
			case
				errors.Is(err, ce.ErrInvalidJWTClaim),
				errors.Is(err, ce.ErrJWTExpired),
				errors.Is(err, ce.ErrJWTMalformed):
				ce.NewError(
					ce.CodeUnauthenticated,
					ce.MsgUnauthenticated,
					err,
				).Bind(
					ctx,
				)
			default:
				ce.NewError(
					ce.CodeUnknown,
					ce.MsgInternalServer,
					err,
				).Bind(
					ctx,
				)
			}

			ctx.Abort()
			return
		}

		var pharmacyID *uuid.UUID
		if claim.Role == constants.RolePharmacy && claim.PharmacyID != nil {
			id := utils.ToUUID(*claim.PharmacyID)
			pharmacyID = &id
		}

		ctx.Request = ctx.Request.WithContext(
			context.WithValue(
				ctx.Request.Context(),
				constants.CtxKeyAuth,
				&utils.AuthContext{
					AuthID:          claim.AuthID,
					Role:            claim.Role,
					IsEmailVerified: claim.IsEmailVerified,
					PharmacyID:      pharmacyID,
				},
			),
		)
		ctx.Next()
	}
}
