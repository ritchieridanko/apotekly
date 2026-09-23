package utils

type (
	OffsetPagination struct {
		Page     int
		PageSize int
	}

	Sorter struct {
		IsAsc bool
	}

	DefaultSorters struct {
		ByCreatedAt *Sorter
		ByUpdatedAt *Sorter
	}
)

func (p *OffsetPagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}
