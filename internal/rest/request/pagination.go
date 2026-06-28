package request

// PaginationReq represents pagination request
type PaginationReq[T any] struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
	Filter T      `query:"filter"`
}
