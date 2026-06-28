package request

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var (
	hasLowercase = regexp.MustCompile(`[a-z]`)
	hasUppercase = regexp.MustCompile(`[A-Z]`)
	hasNumber    = regexp.MustCompile(`[0-9]`)
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type VerifyEmailRequest struct {
	Token string `param:"token" json:"token"`
}

func BlacklistWithoutSpaceValidation(field string) validation.RuleFunc {
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

func StrongPassword(value interface{}) error {
	s, ok := value.(string)
	if !ok || s == "" {
		return nil
	}

	if len(s) < 6 {
		return errors.New("Password must be at least 6 characters")
	}
	if !hasUppercase.MatchString(s) {
		return errors.New("Password must contain at least 1 uppercase letter")
	}
	if !hasLowercase.MatchString(s) {
		return errors.New("Password must contain at least 1 lowercase letter")
	}
	if !hasNumber.MatchString(s) {
		return errors.New("Password must contain at least 1 number")
	}

	return nil
}

func (request RegisterRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Name, validation.Required, validation.Length(3, 100)),
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.Password, validation.Required, validation.Length(6, 255), validation.By(BlacklistWithoutSpaceValidation("password")), validation.By(StrongPassword)),
	)
}

func (request LoginRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.Password, validation.Required, validation.By(BlacklistWithoutSpaceValidation("password"))),
	)
}

func (request RefreshTokenRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.RefreshToken, validation.Required),
	)
}

func (request ForgotPasswordRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Email, validation.Required, is.Email),
	)
}

func (request ResetPasswordRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Token, validation.Required),
		validation.Field(&request.NewPassword, validation.Required, validation.Length(6, 255), validation.By(BlacklistWithoutSpaceValidation("newPassword")), validation.By(StrongPassword)),
	)
}

func (request VerifyEmailRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Token, validation.Required),
	)
}
