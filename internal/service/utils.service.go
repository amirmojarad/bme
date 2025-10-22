package service

const (
	DefaultCurrentPage = 1
	DefaultPerPage     = 10
)

type PaginationRequest struct {
	CurrentPage int
	PerPage     int
}

type PaginationMeta struct {
	CurrentPage int
	PerPage     int
	Total       int
	TotalPages  int
}

func NewPaginationRequest(currentPage int, perPage int) *PaginationRequest {
	return &PaginationRequest{
		CurrentPage: currentPage,
		PerPage:     perPage,
	}
}

func (req PaginationRequest) PaginationMeta() *PaginationMeta {
	return &PaginationMeta{
		CurrentPage: req.CurrentPage,
		PerPage:     req.PerPage,
	}
}
