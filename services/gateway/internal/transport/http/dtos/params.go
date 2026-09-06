package dtos

type PaginationParams struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=10"`
}
