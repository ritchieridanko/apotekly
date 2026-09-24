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

func (h *PharmacyHandler) GetID(ctx context.Context, req *apis.PharmacyGetIDRequest) (*apis.PharmacyGetIDResponse, error) {
	pharmacyID, err := h.pu.GetID(ctx, req.GetAuthId())
	if err != nil {
		return nil, err
	}
	return &apis.PharmacyGetIDResponse{PharmacyId: pharmacyID.String()}, nil
}

func (h *PharmacyHandler) GetPharmacyByID(ctx context.Context, req *apis.GetPharmacyByIDRequest) (*apis.GetPharmacyByIDResponse, error) {
	p, err := h.pu.GetPharmacyByID(ctx, utils.ToUUID(req.GetPharmacyId()))
	if err != nil {
		return nil, err
	}
	return &apis.GetPharmacyByIDResponse{Pharmacy: h.toPharmacy(p)}, nil
}

func (h *PharmacyHandler) GetAllPharmacies(ctx context.Context, req *apis.GetAllPharmaciesRequest) (*apis.GetAllPharmaciesResponse, error) {
	pss, total, err := h.pu.GetAllPharmacies(
		ctx,
		&models.GetAllPharmaciesReq{
			Search:    req.Search,
			RadiusM:   req.RadiusM,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,

			ByLocation: h.toSorter(req.GetByLocation()),
			DefaultSorters: utils.DefaultSorters{
				ByCreatedAt: h.toSorter(req.GetByCreatedAt()),
				ByUpdatedAt: h.toSorter(req.GetByUpdatedAt()),
			},

			OffsetPagination: utils.OffsetPagination{
				Page:     int(req.GetPage()),
				PageSize: int(req.GetPageSize()),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	pharmacies := make([]*apis.PharmacySummary, 0, len(pss))
	for _, ps := range pss {
		pharmacies = append(pharmacies, h.toPharmacySummary(&ps))
	}
	return &apis.GetAllPharmaciesResponse{
		Pharmacies: pharmacies,
		Total:      total,
	}, nil
}

func (h *PharmacyHandler) UpdatePharmacy(ctx context.Context, req *apis.UpdatePharmacyRequest) (*apis.UpdatePharmacyResponse, error) {
	p, err := h.pu.UpdatePharmacy(
		ctx,
		&models.UpdatePharmacyReq{
			Name:         req.Name,
			LegalName:    req.LegalName,
			Description:  req.Description,
			OnlineHours:  utils.ToJSON(req.GetOnlineHours()),
			Country:      req.Country,
			Subdivision1: req.Subdivision_1,
			Subdivision2: req.Subdivision_2,
			Subdivision3: req.Subdivision_3,
			Subdivision4: req.Subdivision_4,
			Street:       req.Street,
			PostalCode:   req.PostalCode,
			Latitude:     req.Latitude,
			Longitude:    req.Longitude,
			Email:        req.Email,
			Phone:        req.Phone,
			Website:      req.Website,
			Whatsapp:     req.Whatsapp,
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.UpdatePharmacyResponse{Pharmacy: h.toPharmacy(p)}, nil
}

func (h *PharmacyHandler) UpdateProfilePicture(ctx context.Context, req *apis.PharmacyUpdateProfilePictureRequest) (*apis.PharmacyUpdateProfilePictureResponse, error) {
	p, err := h.pu.UpdateProfilePicture(ctx, req.GetProfilePictureUrl())
	if err != nil {
		return nil, err
	}
	return &apis.PharmacyUpdateProfilePictureResponse{Pharmacy: h.toPharmacy(p)}, nil
}

func (h *PharmacyHandler) UpdateProfileBanner(ctx context.Context, req *apis.PharmacyUpdateProfileBannerRequest) (*apis.PharmacyUpdateProfileBannerResponse, error) {
	p, err := h.pu.UpdateProfileBanner(ctx, req.GetProfileBannerUrl())
	if err != nil {
		return nil, err
	}
	return &apis.PharmacyUpdateProfileBannerResponse{Pharmacy: h.toPharmacy(p)}, nil
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

func (h *PharmacyHandler) toPharmacySummary(ps *models.PharmacySummary) *apis.PharmacySummary {
	if ps == nil {
		return nil
	}
	return &apis.PharmacySummary{
		Id:             ps.ID.String(),
		Name:           ps.Name,
		LegalName:      ps.LegalName,
		OnlineHours:    utils.ToByte(ps.OnlineHours),
		ProfilePicture: ps.ProfilePicture,
		DistanceM:      ps.DistanceM,
		CreatedAt:      utils.ToTimestamp(&ps.CreatedAt),
		UpdatedAt:      utils.ToTimestamp(&ps.UpdatedAt),
	}
}

func (h *PharmacyHandler) toSorter(s *apis.Sorter) *utils.Sorter {
	if s == nil {
		return nil
	}
	return &utils.Sorter{
		IsAsc: s.GetIsAsc(),
	}
}
