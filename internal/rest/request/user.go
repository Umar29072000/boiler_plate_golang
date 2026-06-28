package request

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type GetUserRequest struct {
	PaginationRequest
	ID          string `query:"id" json:"id"`
	Email       string `query:"email" json:"email"`
	PhoneNumber string `query:"phoneNumber" json:"phoneNumber"`
}

type UpdateProfileRequest struct {
	Name string `json:"name"`
}

type DeleteUserRequest struct {
	ID string `param:"id" json:"id"`
}

func BlacklistValidation(field string) validation.RuleFunc {
	return func(value interface{}) error {
		val, ok := value.(string)

		if !ok {
			return errors.New("The " + field + " is not a string")
		}

		if val == "" {
			return nil
		}

		match, _ := regexp.MatchString(`^[^'"\\[\]<>\\{\\}]+$`, val)

		if !match {
			return errors.New("The " + field + " contains unsafe characters")
		}

		return nil
	}
}

func (request GetUserRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Page, is.Digit),
		validation.Field(&request.Limit, is.Digit),
		validation.Field(&request.Field, validation.In("id", "name", "email", "createdAt", "updatedAt")),
		validation.Field(&request.Sort, validation.In("asc", "desc")),
		validation.Field(&request.Search, validation.By(BlacklistValidation("search"))),
		validation.Field(&request.DisableCalculateTotal, validation.In("true", "false")),
		validation.Field(&request.ID, validation.By(BlacklistValidation("id"))),
		validation.Field(&request.Email, is.Email),
	)
}

func (request UpdateProfileRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Name, validation.Required, validation.Length(3, 100)),
	)
}

func (request DeleteUserRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.ID, validation.Required),
	)
}
