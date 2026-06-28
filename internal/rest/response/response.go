package response

// BaseResponse represents standard API response
type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// PaginationResponse represents paginated API response
type PaginationResponse struct {
	Items      interface{} `json:"items"`
	TotalItems int64       `json:"totalItems"`
	TotalPages int64       `json:"totalPages"`
	Page       int64       `json:"page"`
	Limit      int64       `json:"limit"`
}
