package schema

import "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/service"

type PaginationResponse struct {
	Page      uint32 `json:"page"`
	PageCount uint32 `json:"page_count"`
	PageSize  uint32 `json:"page_size"`
	Size      uint32 `json:"size"`
}

func NewPaginationResponse(result service.PaginationResult) PaginationResponse {
	return PaginationResponse{
		Page:      result.Page,
		PageCount: result.PageCount,
		PageSize:  result.PageSize,
		Size:      result.Size,
	}
}
