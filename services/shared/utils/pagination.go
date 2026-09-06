package utils

type (
	OffsetPagination struct {
		Page     int
		PageSize int
	}
)

func (p *OffsetPagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}
