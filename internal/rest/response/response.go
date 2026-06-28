package response

// BaseResponse represents standard API response
type BaseResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// PaginationResp represents paginated API response
type PaginationResp[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []T    `json:"data"`
	Meta    Meta   `json:"meta"`
}

// Meta represents pagination metadata
type Meta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}
