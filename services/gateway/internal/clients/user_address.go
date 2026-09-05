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
	CreateAddress(ctx context.Context, req *models.CreateAddressReq) (a *models.Address, oldPrimary *models.Address, err *ce.Error)
}

type addressClient struct {
	client apis.AddressServiceClient
}

func NewAddressClient(c apis.AddressServiceClient) AddressClient {
	return &addressClient{client: c}
}

func (c *addressClient) CreateAddress(ctx context.Context, req *models.CreateAddressReq) (*models.Address, *models.Address, *ce.Error) {
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
		return nil, nil, ce.ToError(
			err,
		).Append(
			addressServiceField,
		)
	}
	return c.toAddress(resp.GetAddress()), c.toAddress(resp.GetOldPrimaryAddress()), nil
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
