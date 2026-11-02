package common

// PaginationQuery parses common query parameters for listings.
type PaginationQuery struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// Normalize ensures safe page and limit boundaries.
func (p *PaginationQuery) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
}

// Offset returns the SQL OFFSET value.
func (p *PaginationQuery) Offset() int {
	return (p.Page - 1) * p.Limit
}

// PaginationMeta provides envelope metadata for paginated results.
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// NewPaginationMeta computes metadata from totals.
func NewPaginationMeta(page, limit int, totalItems int64) PaginationMeta {
	totalPages := 0
	if limit > 0 {
		totalPages = int((totalItems + int64(limit) - 1) / int64(limit))
	}
	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
