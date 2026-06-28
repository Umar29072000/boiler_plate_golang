package request

type PaginationRequest struct {
	Page                  string `query:"page" json:"page"`
	Limit                 string `query:"limit" json:"limit"`
	Field                 string `query:"field" json:"field"`
	Sort                  string `query:"sort" json:"sort"`
	Search                string `query:"search" json:"search"`
	DisableCalculateTotal string `query:"disableCalculateTotal" json:"disableCalculateTotal"`
}
