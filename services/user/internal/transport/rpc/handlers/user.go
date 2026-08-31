package handlers

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/user/internal/models"
	"github.com/ritchieridanko/apotekly/services/user/internal/usecases"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	apis.UnimplementedUserServiceServer
	uu usecases.UserUsecase
}

func NewUserHandler(uu usecases.UserUsecase) *UserHandler {
	return &UserHandler{uu: uu}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *apis.CreateUserRequest) (*apis.CreateUserResponse, error) {
	u, err := h.uu.CreateUser(
		ctx,
		&models.CreateUserReq{
			Name:      req.GetName(),
			Sex:       req.Sex,
			Birthdate: utils.ToTime(req.Birthdate),
			Phone:     req.Phone,
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.CreateUserResponse{User: h.toUser(u)}, nil
}

func (h *UserHandler) GetMe(ctx context.Context, req *emptypb.Empty) (*apis.GetMeResponse, error) {
	u, err := h.uu.GetMe(ctx)
	if err != nil {
		return nil, err
	}
	return &apis.GetMeResponse{User: h.toUser(u)}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *apis.UpdateUserRequest) (*apis.UpdateUserResponse, error) {
	u, err := h.uu.UpdateUser(
		ctx,
		&models.UpdateUserReq{
			Name:      req.Name,
			Sex:       req.Sex,
			Birthdate: utils.ToTime(req.Birthdate),
			Phone:     req.Phone,
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.UpdateUserResponse{User: h.toUser(u)}, nil
}

func (h *UserHandler) toUser(u *models.User) *apis.User {
	if u == nil {
		return nil
	}
	return &apis.User{
		Id:             u.ID.String(),
		Name:           u.Name,
		Sex:            u.Sex,
		Birthdate:      utils.ToTimestamp(u.Birthdate),
		Phone:          u.Phone,
		ProfilePicture: u.ProfilePicture,
		ProfileBanner:  u.ProfileBanner,
	}
}
