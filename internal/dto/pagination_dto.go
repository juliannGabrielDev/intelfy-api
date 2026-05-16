package dto

type GenericPagination struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type PaginationMeta struct {
	TotalRecords int64 `json:"totalRecords"`
	CurrentPage  int   `json:"currentPage"`
	TotalPages   int   `json:"totalPages"`
	Limit        int   `json:"limit"`
}

type PaginatedResponse[T any] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
