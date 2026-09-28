package clients

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

var pharmacyServiceField logger.Field = logger.NewField("service", "pharmacy")

type PharmacyClient interface {
	GetStatus(ctx context.Context, authID uint64) (status string, err *ce.Error)
}

type pharmacyClient struct {
	client apis.PharmacyServiceClient
}

func NewPharmacyClient(c apis.PharmacyServiceClient) PharmacyClient {
	return &pharmacyClient{client: c}
}

func (c *pharmacyClient) GetStatus(ctx context.Context, authID uint64) (string, *ce.Error) {
	resp, err := c.client.GetStatus(
		ctx,
		&apis.PharmacyGetStatusRequest{
			AuthId: authID,
		},
	)
	if err != nil {
		return "", ce.ToError(
			err,
		).Append(
			pharmacyServiceField,
		)
	}
	return resp.GetStatus(), nil
}
