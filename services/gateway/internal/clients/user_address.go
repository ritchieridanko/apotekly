package clients

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/gateway/internal/models"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

var addressServiceField logger.Field = logger.NewField("service", "user.address")

type AddressClient interface {
	CreateAddress(ctx context.Context, req *models.CreateAddressReq) (a *models.Address, err *ce.Error)
	GetAllAddresses(ctx context.Context, req *models.GetAllAddressesReq) (as []models.Address, total int64, err *ce.Error)
	UpdateAddress(ctx context.Context, req *models.UpdateAddressReq) (a *models.Address, err *ce.Error)
	SetPrimaryAddress(ctx context.Context, addressID uint64) (a *models.Address, err *ce.Error)
}

type addressClient struct {
	client apis.AddressServiceClient
}

func NewAddressClient(c apis.AddressServiceClient) AddressClient {
	return &addressClient{client: c}
}

func (c *addressClient) CreateAddress(ctx context.Context, req *models.CreateAddressReq) (*models.Address, *ce.Error) {
	resp, err := c.client.CreateAddress(
		ctx,
		&apis.CreateAddressRequest{
			Label:         req.Label,
			Recipient:     req.Recipient,
			Phone:         req.Phone,
			Notes:         req.Notes,
			IsPrimary:     req.IsPrimary,
			Country:       req.Country,
			Subdivision_1: req.Subdivision1,
			Subdivision_2: req.Subdivision2,
			Subdivision_3: req.Subdivision3,
			Subdivision_4: req.Subdivision4,
			Street:        req.Street,
			PostalCode:    req.PostalCode,
			Latitude:      req.Latitude,
			Longitude:     req.Longitude,
		},
	)
	if err != nil {
		return nil, ce.ToError(
			err,
		).Append(
			addressServiceField,
		)
	}
	return c.toAddress(resp.GetAddress()), nil
}

func (c *addressClient) GetAllAddresses(ctx context.Context, req *models.GetAllAddressesReq) ([]models.Address, int64, *ce.Error) {
	resp, err := c.client.GetAllAddresses(
		ctx,
		&apis.GetAllAddressesRequest{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	)
	if err != nil {
		return nil, 0, ce.ToError(
			err,
		).Append(
			addressServiceField,
		)
	}

	as := make([]models.Address, 0, len(resp.GetAddresses()))
	for _, a := range resp.GetAddresses() {
		if a == nil {
			continue
		}
		as = append(as, *c.toAddress(a))
	}

	return as, resp.GetTotal(), nil
}

func (c *addressClient) UpdateAddress(ctx context.Context, req *models.UpdateAddressReq) (*models.Address, *ce.Error) {
	resp, err := c.client.UpdateAddress(
		ctx,
		&apis.UpdateAddressRequest{
			AddressId: req.AddressID,

			Label:         req.Label,
			Recipient:     req.Recipient,
			Phone:         req.Phone,
			Notes:         req.Notes,
			Country:       req.Country,
			Subdivision_1: req.Subdivision1,
			Subdivision_2: req.Subdivision2,
			Subdivision_3: req.Subdivision3,
			Subdivision_4: req.Subdivision4,
			Street:        req.Street,
			PostalCode:    req.PostalCode,
			Latitude:      req.Latitude,
			Longitude:     req.Longitude,
		},
	)
	if err != nil {
		return nil, ce.ToError(
			err,
		).Append(
			addressServiceField,
		)
	}
	return c.toAddress(resp.GetAddress()), nil
}

func (c *addressClient) SetPrimaryAddress(ctx context.Context, addressID uint64) (*models.Address, *ce.Error) {
	resp, err := c.client.SetPrimaryAddress(
		ctx,
		&apis.SetPrimaryAddressRequest{
			AddressId: addressID,
		},
	)
	if err != nil {
		return nil, ce.ToError(
			err,
		).Append(
			addressServiceField,
		)
	}
	return c.toAddress(resp.GetAddress()), nil
}

func (c *addressClient) toAddress(a *apis.Address) *models.Address {
	if a == nil {
		return nil
	}
	return &models.Address{
		ID:           a.GetId(),
		Label:        a.GetLabel(),
		Recipient:    a.GetRecipient(),
		Phone:        a.GetPhone(),
		Notes:        a.Notes,
		IsPrimary:    a.GetIsPrimary(),
		Country:      a.GetCountry(),
		Subdivision1: a.Subdivision_1,
		Subdivision2: a.Subdivision_2,
		Subdivision3: a.Subdivision_3,
		Subdivision4: a.Subdivision_4,
		Street:       a.GetStreet(),
		PostalCode:   a.GetPostalCode(),
		Latitude:     a.GetLatitude(),
		Longitude:    a.GetLongitude(),
		CreatedAt:    utils.ToTime(a.GetCreatedAt()),
		UpdatedAt:    utils.ToTime(a.GetUpdatedAt()),
	}
}
