package user

import (
	"time"

	"github.com/Sanchir01/users-info/internal/gender"
	api "github.com/Sanchir01/users-info/pkg/lib/api/response"
	"github.com/google/uuid"
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
type GetAllUsersResponse struct {
	api.Response
	Users        []*UserDB `json:"users"`
	Page         uint      `json:"page"`
	ItemsPerPage uint      `json:"items_per_page"`
}

type PaginationParams struct {
	Page     uint `json:"page" validate:"omitempty,gte=1"`
	PageSize uint `json:"page_size" validate:"omitempty,gte=1,lte=100"`
}
type UpdateUserRequest struct {
	Name       string `json:"name" validate:"required,min=1,max=100"`
	Surname    string `json:"surname" validate:"required,min=1,max=100"`
	Patronymic string `json:"patronymic,omitempty" validate:"omitempty,max=100"`
}
type UpdateUserResponse struct {
	api.Response
	Ok string `json:"ok" validate:"required"`
}
type DeleteUserResponse struct {
	api.Response
	Ok string `json:"ok" validate:"required"`
}
type UserDB struct {
	ID          uuid.UUID     `db:"id" json:"id"`
	Name        string        `db:"name" json:"name"`
	Surname     string        `db:"surname" json:"surname"`
	Patronymic  string        `db:"patronymic" json:"patronymic,omitempty"`
	CreatedAt   time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at" json:"updated_at"`
	Age         int           `db:"age" json:"age"`
	Gender      gender.Gender `db:"gender" json:"gender"`
	Nationality string        `db:"nationality" json:"nationality"`
	Version     int64         `db:"version" json:"version"`
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
