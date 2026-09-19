package clients

import (
	"context"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

var pharmacyServiceField logger.Field = logger.NewField("service", "pharmacy")

type PharmacyClient interface {
	GetID(ctx context.Context, authID uint64) (pharmacyID uuid.UUID, err *ce.Error)
}

type pharmacyClient struct {
	client apis.PharmacyServiceClient
}

func NewPharmacyClient(c apis.PharmacyServiceClient) PharmacyClient {
	return &pharmacyClient{client: c}
}

func (c *pharmacyClient) GetID(ctx context.Context, authID uint64) (uuid.UUID, *ce.Error) {
	resp, err := c.client.GetID(
		ctx,
		&apis.PharmacyGetIDRequest{
			AuthId: authID,
		},
	)
	if err != nil {
		return uuid.Nil, ce.ToError(
			err,
		).Append(
			pharmacyServiceField,
		)
	}
	return utils.ToUUID(resp.GetPharmacyId()), nil
}
