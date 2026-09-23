package clients

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

var addressServiceField logger.Field = logger.NewField("service", "user.address")

type AddressClient interface {
	GetPrimaryLocation(ctx context.Context, authID uint64) (lat, lon float64, err *ce.Error)
}

type addressClient struct {
	client apis.AddressServiceClient
}

func NewAddressClient(c apis.AddressServiceClient) AddressClient {
	return &addressClient{client: c}
}

func (c *addressClient) GetPrimaryLocation(ctx context.Context, authID uint64) (float64, float64, *ce.Error) {
	resp, err := c.client.GetPrimaryLocation(
		ctx,
		&apis.AddressGetPrimaryLocationRequest{
			AuthId: authID,
		},
	)
	if err != nil {
		return 0, 0, ce.ToError(
			err,
		).Append(
			addressServiceField,
		)
	}
	return resp.GetLatitude(), resp.GetLongitude(), nil
}
