package controller

import (
	"bme/internal/service"
	"bme/pkg/jwt"
	"bme/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthClaims struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type UserRegisterResponse struct {
	Tokens
}

type UserLoginRequest struct {
	AuthClaims
}

type UserLoginResponse struct {
	Tokens
}

type UserRegisterRequest struct {
	AuthClaims
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

type RefreshTokenRequest struct {
	RefreshToken string `header:"refresh_token" binding:"required"`
}

func tokensFromJwtTokens(tokens jwt.Tokens) Tokens {
	return Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}

func (req AuthClaims) validate() (*utils.ValidationError, error) {
	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		validationErr := &utils.ValidationError{
			Success: false,
			Message: "",
			Meta:    make([]map[string]any, 0),
		}

		for _, err = range err.(validator.ValidationErrors) {
			var fieldErr validator.FieldError
			errors.As(err, &fieldErr)

			serviceValidationErrorMeta := utils.ServiceValidationErrorMeta{
				Attribute: "",
				Message:   fieldErr.Error(),
				Error:     err,
			}

			switch fieldErr.Tag() {
			case "required":
				switch fieldErr.Field() {
				case "Username":
					serviceValidationErrorMeta.Message = "Username is required"
					serviceValidationErrorMeta.Attribute = "Username"
				case "Password":
					serviceValidationErrorMeta.Message = "Password is required"
					serviceValidationErrorMeta.Attribute = "Password"
				}

			case "min":
				switch fieldErr.Field() {
				case "Username":
					serviceValidationErrorMeta.Message = "Username minimum length required"
					serviceValidationErrorMeta.Attribute = "Username"
				case "Password":
					serviceValidationErrorMeta.Message = "Password minimum length required"
					serviceValidationErrorMeta.Attribute = "Password"
				}
			}

			validationErr.Meta = append(validationErr.Meta, serviceValidationErrorMeta.Map())
		}

		return validationErr, nil
	}

	return nil, nil
}

func (req UserLoginRequest) toSvc() service.AuthLoginRequest {
	return service.AuthLoginRequest{
		AuthClaims: req.AuthClaims.toSvc(),
	}
}

func (req AuthClaims) toSvc() service.AuthClaims {
	return service.AuthClaims{
		Username: req.Username,
		Password: req.Password,
	}
}

func (req UserRegisterRequest) toSvc() service.AuthRegisterRequest {
	return service.AuthRegisterRequest{
		AuthClaims:  req.AuthClaims.toSvc(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
	}
}
