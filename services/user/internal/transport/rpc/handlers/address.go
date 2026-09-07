package handlers

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/user/internal/models"
	"github.com/ritchieridanko/apotekly/services/user/internal/usecases"
)

type AddressHandler struct {
	apis.UnimplementedAddressServiceServer
	au usecases.AddressUsecase
}

func NewAddressHandler(au usecases.AddressUsecase) *AddressHandler {
	return &AddressHandler{au: au}
}

func (h *AddressHandler) CreateAddress(ctx context.Context, req *apis.CreateAddressRequest) (*apis.CreateAddressResponse, error) {
	a, err := h.au.CreateAddress(
		ctx,
		&models.CreateAddressReq{
			Label:        req.GetLabel(),
			Recipient:    req.GetRecipient(),
			Phone:        req.GetPhone(),
			Notes:        req.Notes,
			IsPrimary:    req.GetIsPrimary(),
			Country:      req.GetCountry(),
			Subdivision1: req.Subdivision_1,
			Subdivision2: req.Subdivision_2,
			Subdivision3: req.Subdivision_3,
			Subdivision4: req.Subdivision_4,
			Street:       req.GetStreet(),
			PostalCode:   req.GetPostalCode(),
			Latitude:     req.GetLatitude(),
			Longitude:    req.GetLongitude(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.CreateAddressResponse{Address: h.toAddress(a)}, nil
}

func (h *AddressHandler) GetAllAddresses(ctx context.Context, req *apis.GetAllAddressesRequest) (*apis.GetAllAddressesResponse, error) {
	as, total, err := h.au.GetAllAddresses(
		ctx,
		&models.GetAllAddressesReq{
			OffsetPagination: utils.OffsetPagination{
				Page:     int(req.GetPage()),
				PageSize: int(req.GetPageSize()),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	addresses := make([]*apis.Address, 0, len(as))
	for _, a := range as {
		addresses = append(addresses, h.toAddress(&a))
	}
	return &apis.GetAllAddressesResponse{
		Addresses: addresses,
		Total:     total,
	}, nil
}

func (h *AddressHandler) UpdateAddress(ctx context.Context, req *apis.UpdateAddressRequest) (*apis.UpdateAddressResponse, error) {
	a, err := h.au.UpdateAddress(
		ctx,
		&models.UpdateAddressReq{
			AddressID: req.GetAddressId(),

			Label:        req.Label,
			Recipient:    req.Recipient,
			Phone:        req.Phone,
			Notes:        req.Notes,
			Country:      req.Country,
			Subdivision1: req.Subdivision_1,
			Subdivision2: req.Subdivision_2,
			Subdivision3: req.Subdivision_3,
			Subdivision4: req.Subdivision_4,
			Street:       req.Street,
			PostalCode:   req.PostalCode,
			Latitude:     req.Latitude,
			Longitude:    req.Longitude,
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.UpdateAddressResponse{Address: h.toAddress(a)}, nil
}

func (h *AddressHandler) SetPrimaryAddress(ctx context.Context, req *apis.SetPrimaryAddressRequest) (*apis.SetPrimaryAddressResponse, error) {
	a, err := h.au.SetPrimaryAddress(ctx, req.GetAddressId())
	if err != nil {
		return nil, err
	}
	return &apis.SetPrimaryAddressResponse{Address: h.toAddress(a)}, nil
}

func (h *AddressHandler) toAddress(a *models.Address) *apis.Address {
	if a == nil {
		return nil
	}
	return &apis.Address{
		Id:            a.ID,
		Label:         utils.ToTitlecase(a.Label),
		Recipient:     a.Recipient,
		Phone:         a.Phone,
		Notes:         a.Notes,
		IsPrimary:     a.IsPrimary,
		Country:       utils.ToTitlecase(a.Country),
		Subdivision_1: utils.ToTitlecasePtr(a.Subdivision1),
		Subdivision_2: utils.ToTitlecasePtr(a.Subdivision2),
		Subdivision_3: utils.ToTitlecasePtr(a.Subdivision3),
		Subdivision_4: utils.ToTitlecasePtr(a.Subdivision4),
		Street:        a.Street,
		PostalCode:    a.PostalCode,
		Latitude:      a.Latitude,
		Longitude:     a.Longitude,
		CreatedAt:     utils.ToTimestamp(&a.CreatedAt),
		UpdatedAt:     utils.ToTimestamp(&a.UpdatedAt),
	}
}
