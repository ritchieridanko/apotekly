package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

type (
	httpResponse[T any] struct {
		Status   int                   `json:"status"`
		Message  string                `json:"message"`
		Data     T                     `json:"data,omitempty"`
		Metadata *httpResponseMetadata `json:"metadata,omitempty"`
	}

	httpResponseMetadata struct {
		RequestID string    `json:"request_id"`
		Page      *int      `json:"page,omitempty"`
		PageSize  *int      `json:"page_size,omitempty"`
		Total     *int      `json:"total,omitempty"`
		Timestamp time.Time `json:"timestamp"`
	}

	ResponseMetadata struct {
		Page     int
		PageSize int
		Total    int
	}
)

func SetHTTPResponse[T any](ctx *gin.Context, status int, message string, data T, rm *ResponseMetadata) {
	var page, pageSize, total *int
	if rm != nil {
		page, pageSize, total = &rm.Page, &rm.PageSize, &rm.Total
	}
	ctx.JSON(
		status,
		httpResponse[T]{
			Status:  status,
			Message: message,
			Data:    data,
			Metadata: &httpResponseMetadata{
				RequestID: CtxRequestID(ctx.Request.Context()),
				Page:      page,
				PageSize:  pageSize,
				Total:     total,
				Timestamp: time.Now().UTC(),
			},
		},
	)
}
