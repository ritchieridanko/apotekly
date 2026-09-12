package handlers

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/models"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/usecases"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PharmacyHandler struct {
	apis.UnimplementedPharmacyServiceServer
	pu usecases.PharmacyUsecase
}

func NewPharmacyHandler(pu usecases.PharmacyUsecase) *PharmacyHandler {
	return &PharmacyHandler{pu: pu}
}

func (h *PharmacyHandler) CreatePharmacy(ctx context.Context, req *apis.CreatePharmacyRequest) (*apis.CreatePharmacyResponse, error) {
	p, err := h.pu.CreatePharmacy(
		ctx,
		&models.CreatePharmacyReq{
			Name:         req.GetName(),
			LegalName:    req.LegalName,
			Description:  req.Description,
			OnlineHours:  utils.ToJSON(req.GetOnlineHours()),
			Country:      req.GetCountry(),
			Subdivision1: req.Subdivision_1,
			Subdivision2: req.Subdivision_2,
			Subdivision3: req.Subdivision_3,
			Subdivision4: req.Subdivision_4,
			Street:       req.GetStreet(),
			PostalCode:   req.GetPostalCode(),
			Latitude:     req.GetLatitude(),
			Longitude:    req.GetLongitude(),
			Email:        req.Email,
			Phone:        req.Phone,
			Website:      req.Website,
			Whatsapp:     req.Whatsapp,
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.CreatePharmacyResponse{Pharmacy: h.toPharmacy(p)}, nil
}

func (h *PharmacyHandler) GetMe(ctx context.Context, req *emptypb.Empty) (*apis.PharmacyGetMeResponse, error) {
	p, err := h.pu.GetMe(ctx)
	if err != nil {
		return nil, err
	}
	return &apis.PharmacyGetMeResponse{Pharmacy: h.toPharmacy(p)}, nil
}

func (h *PharmacyHandler) toPharmacy(p *models.Pharmacy) *apis.Pharmacy {
	if p == nil {
		return nil
	}
	return &apis.Pharmacy{
		Id:             p.ID.String(),
		Name:           p.Name,
		LegalName:      p.LegalName,
		Description:    p.Description,
		Status:         p.Status,
		OnlineHours:    utils.ToByte(p.OnlineHours),
		Country:        utils.ToTitlecase(p.Country),
		Subdivision_1:  utils.ToTitlecasePtr(p.Subdivision1),
		Subdivision_2:  utils.ToTitlecasePtr(p.Subdivision2),
		Subdivision_3:  utils.ToTitlecasePtr(p.Subdivision3),
		Subdivision_4:  utils.ToTitlecasePtr(p.Subdivision4),
		Street:         p.Street,
		PostalCode:     p.PostalCode,
		Latitude:       p.Latitude,
		Longitude:      p.Longitude,
		Email:          p.Email,
		Phone:          p.Phone,
		Website:        p.Website,
		Whatsapp:       p.Whatsapp,
		ProfilePicture: p.ProfilePicture,
		ProfileBanner:  p.ProfileBanner,
		VerifiedAt:     utils.ToTimestamp(p.VerifiedAt),
		CreatedAt:      utils.ToTimestamp(&p.CreatedAt),
		UpdatedAt:      utils.ToTimestamp(&p.UpdatedAt),
	}
}
