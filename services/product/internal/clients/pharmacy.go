package clients

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

var pharmacyServiceField logger.Field = logger.NewField("service", "pharmacy")

type PharmacyClient interface {
	GetActiveStatus(ctx context.Context, authID uint64) (active bool, err *ce.Error)
}

type pharmacyClient struct {
	client apis.PharmacyServiceClient
}

func NewPharmacyClient(c apis.PharmacyServiceClient) PharmacyClient {
	return &pharmacyClient{client: c}
}

func (c *pharmacyClient) GetActiveStatus(ctx context.Context, authID uint64) (bool, *ce.Error) {
	resp, err := c.client.GetActiveStatus(
		ctx,
		&apis.PharmacyGetActiveStatusRequest{
			AuthId: authID,
		},
	)
	if err != nil {
		return false, ce.ToError(
			err,
		).Append(
			pharmacyServiceField,
		)
	}
	return resp.GetIsActive(), nil
}
