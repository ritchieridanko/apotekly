package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/clients"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/models"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/dtos"
	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

const (
	defaultPageSize int = 10
	maxPageSize     int = 100
)

type AddressHandler struct {
	uac clients.AddressClient
}

func NewAddressHandler(uac clients.AddressClient) *AddressHandler {
	return &AddressHandler{uac: uac}
}

func (h *AddressHandler) CreateAddress(ctx *gin.Context) {
	var payload dtos.CreateAddressRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err).Bind(ctx)
		return
	}

	authCtx := utils.CtxAuth(ctx.Request.Context())
	if authCtx == nil {
		ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		).Bind(
			ctx,
		)
		return
	}

	a, oldPrimary, err := h.uac.CreateAddress(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
		),
		&models.CreateAddressReq{
			Label:        payload.Label,
			Recipient:    payload.Recipient,
			Phone:        payload.Phone,
			Notes:        payload.Notes,
			IsPrimary:    payload.IsPrimary,
			Country:      payload.Country,
			Subdivision1: payload.Subdivision1,
			Subdivision2: payload.Subdivision2,
			Subdivision3: payload.Subdivision3,
			Subdivision4: payload.Subdivision4,
			Street:       payload.Street,
			PostalCode:   payload.PostalCode,
			Latitude:     payload.Latitude,
			Longitude:    payload.Longitude,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusCreated,
		"Address created successfully",
		dtos.CreateAddressResponse{
			Address:           h.toAddress(a),
			OldPrimaryAddress: h.toAddress(oldPrimary),
		},
		nil,
	)
}

func (h *AddressHandler) GetAllAddresses(ctx *gin.Context) {
	var params dtos.GetAllAddressesRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ce.NewError(ce.CodeInvalidParams, ce.MsgInvalidParams, err).Bind(ctx)
		return
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = defaultPageSize
	}
	if params.PageSize > maxPageSize {
		params.PageSize = maxPageSize
	}

	authCtx := utils.CtxAuth(ctx.Request.Context())
	if authCtx == nil {
		ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		).Bind(
			ctx,
		)
		return
	}

	as, total, err := h.uac.GetAllAddresses(
		utils.CtxWithMetadata(
			ctx.Request.Context(),
			constants.MDKeyAuthID,
			strconv.FormatUint(authCtx.AuthID, 10),
		),
		&models.GetAllAddressesReq{
			Page:     int32(params.Page),
			PageSize: int32(params.PageSize),
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	addresses := make([]dtos.Address, 0, len(as))
	for _, a := range as {
		addresses = append(addresses, *h.toAddress(&a))
	}

	utils.SetHTTPResponse(
		ctx,
		http.StatusOK,
		"Addresses retrieved successfully",
		dtos.GetAllAddressesResponse{
			Addresses: addresses,
		},
		&utils.ResponseMetadata{
			Page:     params.Page,
			PageSize: params.PageSize,
			Total:    total,
		},
	)
}

func (h *AddressHandler) toAddress(a *models.Address) *dtos.Address {
	if a == nil {
		return nil
	}
	return &dtos.Address{
		ID:           a.ID,
		Label:        a.Label,
		Recipient:    a.Recipient,
		Phone:        a.Phone,
		Notes:        a.Notes,
		IsPrimary:    a.IsPrimary,
		Country:      a.Country,
		Subdivision1: a.Subdivision1,
		Subdivision2: a.Subdivision2,
		Subdivision3: a.Subdivision3,
		Subdivision4: a.Subdivision4,
		Street:       a.Street,
		PostalCode:   a.PostalCode,
		Latitude:     a.Latitude,
		Longitude:    a.Longitude,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}
