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
	GetPharmacyByID(ctx context.Context, pharmacyID string) (p *models.Pharmacy, err *ce.Error)
	GetAllPharmacies(ctx context.Context, req *models.GetAllPharmaciesReq) (pss []models.PharmacySummary, total int64, err *ce.Error)
	UpdatePharmacy(ctx context.Context, req *models.UpdatePharmacyReq) (p *models.Pharmacy, err *ce.Error)
	UpdateProfilePicture(ctx context.Context, profilePictureURL string) (p *models.Pharmacy, err *ce.Error)
	UpdateProfileBanner(ctx context.Context, profileBannerURL string) (p *models.Pharmacy, err *ce.Error)
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

func (c *pharmacyClient) GetPharmacyByID(ctx context.Context, pharmacyID string) (*models.Pharmacy, *ce.Error) {
	resp, err := c.client.GetPharmacyByID(
		ctx,
		&apis.GetPharmacyByIDRequest{
			PharmacyId: pharmacyID,
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

func (c *pharmacyClient) GetAllPharmacies(ctx context.Context, req *models.GetAllPharmaciesReq) ([]models.PharmacySummary, int64, *ce.Error) {
	resp, err := c.client.GetAllPharmacies(
		ctx,
		&apis.GetAllPharmaciesRequest{
			Search:    req.Search,
			RadiusM:   req.RadiusM,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,

			ByLocation:  c.toSorter(req.ByLocation),
			ByCreatedAt: c.toSorter(req.ByCreatedAt),
			ByUpdatedAt: c.toSorter(req.ByUpdatedAt),

			Page:     req.Page,
			PageSize: req.PageSize,
		},
	)
	if err != nil {
		return nil, 0, ce.ToError(
			err,
		).Append(
			pharmacyServiceField,
		)
	}

	pss := make([]models.PharmacySummary, 0, len(resp.GetPharmacies()))
	for _, ps := range resp.GetPharmacies() {
		if ps == nil {
			continue
		}
		pss = append(pss, *c.toPharmacySummary(ps))
	}

	return pss, resp.GetTotal(), nil
}

func (c *pharmacyClient) UpdatePharmacy(ctx context.Context, req *models.UpdatePharmacyReq) (*models.Pharmacy, *ce.Error) {
	resp, err := c.client.UpdatePharmacy(
		ctx,
		&apis.UpdatePharmacyRequest{
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

func (c *pharmacyClient) UpdateProfilePicture(ctx context.Context, profilePictureURL string) (*models.Pharmacy, *ce.Error) {
	resp, err := c.client.UpdateProfilePicture(
		ctx,
		&apis.PharmacyUpdateProfilePictureRequest{
			ProfilePictureUrl: profilePictureURL,
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

func (c *pharmacyClient) UpdateProfileBanner(ctx context.Context, profileBannerURL string) (*models.Pharmacy, *ce.Error) {
	resp, err := c.client.UpdateProfileBanner(
		ctx,
		&apis.PharmacyUpdateProfileBannerRequest{
			ProfileBannerUrl: profileBannerURL,
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

func (c *pharmacyClient) toPharmacySummary(ps *apis.PharmacySummary) *models.PharmacySummary {
	if ps == nil {
		return nil
	}
	return &models.PharmacySummary{
		ID:             utils.ToUUID(ps.GetId()),
		Name:           ps.GetName(),
		LegalName:      ps.LegalName,
		OnlineHours:    utils.ToJSON(ps.GetOnlineHours()),
		ProfilePicture: ps.ProfilePicture,
		DistanceM:      ps.DistanceM,
		CreatedAt:      utils.ToTime(ps.GetCreatedAt()),
		UpdatedAt:      utils.ToTime(ps.GetUpdatedAt()),
	}
}

func (c *pharmacyClient) toSorter(s *string) *apis.Sorter {
	if s == nil {
		return nil
	}
	return &apis.Sorter{
		IsAsc: *s == "asc",
	}
}
