package clients

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/gateway/internal/models"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"google.golang.org/protobuf/types/known/emptypb"
)

var pharmacyServiceField logger.Field = logger.NewField("service", "pharmacy")

type PharmacyClient interface {
	CreatePharmacy(ctx context.Context, req *models.CreatePharmacyReq) (p *models.Pharmacy, err *ce.Error)
	GetMe(ctx context.Context) (p *models.Pharmacy, err *ce.Error)
}

type pharmacyClient struct {
	client apis.PharmacyServiceClient
}

func NewPharmacyClient(c apis.PharmacyServiceClient) PharmacyClient {
	return &pharmacyClient{client: c}
}

func (c *pharmacyClient) CreatePharmacy(ctx context.Context, req *models.CreatePharmacyReq) (*models.Pharmacy, *ce.Error) {
	resp, err := c.client.CreatePharmacy(
		ctx,
		&apis.CreatePharmacyRequest{
			Name:          req.Name,
			LegalName:     req.LegalName,
			Description:   req.Description,
			OnlineHours:   utils.ToByte(req.OnlineHours),
			Country:       req.Country,
			Subdivision_1: req.Subdivision1,
			Subdivision_2: req.Subdivision2,
			Subdivision_3: req.Subdivision3,
			Subdivision_4: req.Subdivision4,
			Street:        req.Street,
			PostalCode:    req.PostalCode,
			Latitude:      req.Latitude,
			Longitude:     req.Longitude,
			Email:         req.Email,
			Phone:         req.Phone,
			Website:       req.Website,
			Whatsapp:      req.Whatsapp,
		},
	)
	if err != nil {
		return nil, ce.ToError(
			err,
		).Append(
			pharmacyServiceField,
		)
	}
	return c.toPharmacy(resp.GetPharmacy()), nil
}

func (c *pharmacyClient) GetMe(ctx context.Context) (*models.Pharmacy, *ce.Error) {
	resp, err := c.client.GetMe(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, ce.ToError(
			err,
		).Append(
			pharmacyServiceField,
		)
	}
	return c.toPharmacy(resp.GetPharmacy()), nil
}

func (c *pharmacyClient) toPharmacy(p *apis.Pharmacy) *models.Pharmacy {
	if p == nil {
		return nil
	}
	return &models.Pharmacy{
		ID:             utils.ToUUID(p.GetId()),
		Name:           p.GetName(),
		LegalName:      p.LegalName,
		Description:    p.Description,
		Status:         p.GetStatus(),
		OnlineHours:    utils.ToJSON(p.GetOnlineHours()),
		Country:        p.GetCountry(),
		Subdivision1:   p.Subdivision_1,
		Subdivision2:   p.Subdivision_2,
		Subdivision3:   p.Subdivision_3,
		Subdivision4:   p.Subdivision_4,
		Street:         p.GetStreet(),
		PostalCode:     p.GetPostalCode(),
		Latitude:       p.GetLatitude(),
		Longitude:      p.GetLongitude(),
		Email:          p.Email,
		Phone:          p.Phone,
		Website:        p.Website,
		Whatsapp:       p.Whatsapp,
		ProfilePicture: p.ProfilePicture,
		ProfileBanner:  p.ProfileBanner,
		VerifiedAt:     utils.ToTime(p.GetVerifiedAt()),
		CreatedAt:      utils.ToTime(p.GetCreatedAt()),
		UpdatedAt:      utils.ToTime(p.GetUpdatedAt()),
	}
}
