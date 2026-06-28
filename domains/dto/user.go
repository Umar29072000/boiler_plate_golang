package dto

// GetUserData represents data for getting users
type GetUserData struct {
	Page                  int    `json:"page"`
	Limit                 int    `json:"limit"`
	Field                 string `json:"field"`
	Sort                  string `json:"sort"`
	Search                string `json:"search"`
	DisableCalculateTotal string `json:"disableCalculateTotal"`
	ID                    string `json:"id"`
	Email                 string `json:"email"`
	PhoneNumber           string `json:"phoneNumber"`
}

// CreateUserData represents data for creating a user
type CreateUserData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UpdateUserData represents data for updating a user
type UpdateUserData struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
