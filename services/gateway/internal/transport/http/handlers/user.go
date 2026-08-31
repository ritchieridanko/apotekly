package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/clients"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/models"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/dtos"
	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type UserHandler struct {
	uc clients.UserClient
}

func NewUserHandler(uc clients.UserClient) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var payload dtos.CreateUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err).Bind(ctx)
		return
	}

	authCtx := utils.CtxAuth(ctx.Request.Context())
	if authCtx == nil {
		ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		).Bind(
			ctx,
		)
		return
	}

	u, err := h.uc.CreateUser(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
		),
		&models.CreateUserReq{
			Name:      payload.Name,
			Sex:       payload.Sex,
			Birthdate: payload.Birthdate,
			Phone:     payload.Phone,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusCreated,
		"User created successfully",
		dtos.CreateUserResponse{User: h.toUser(u)},
		nil,
	)
}

func (h *UserHandler) toUser(u *models.User) *dtos.User {
	if u == nil {
		return nil
	}
	return &dtos.User{
		ID:             u.ID.String(),
		Name:           u.Name,
		Sex:            u.Sex,
		Birthdate:      u.Birthdate,
		Phone:          u.Phone,
		ProfilePicture: u.ProfilePicture,
		ProfileBanner:  u.ProfileBanner,
	}
}
