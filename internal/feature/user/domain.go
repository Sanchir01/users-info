package user

import (
	"github.com/Sanchir01/users-info/internal/gender"
	api "github.com/Sanchir01/users-info/pkg/lib/api/response"
)

type CreateUserRequest struct {
	Name       string `json:"name" validate:"required,min=1,max=100"`
	Surname    string `json:"surname" validate:"required,min=1,max=100"`
	Patronymic string `json:"patronymic,omitempty" validate:"omitempty,max=100"`
}
type CreateUserResponse struct {
	api.Response
	Ok string `json:"ok" validate:"required"`
}
type NationalizeResponse struct {
	Name    string              `json:"name"`
	Country []CountryPrediction `json:"country"`
}

type CountryPrediction struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
type GenderizeResponse struct {
	Name        string        `json:"name" validate:"required"`
	Gender      gender.Gender `json:"gender" validate:"required"`
	Probability float64       `json:"probability" validate:"required,gte=0,lte=1"`
	Count       int           `json:"count" validate:"required,gte=0"`
}
type UserAgeResponse struct {
	Name  string `json:"name" validate:"required,min=1,max=100"`
	Age   int    `json:"age" validate:"required,gte=0,lte=120"`
	Count int    `json:"count" validate:"gte=0"`
}
