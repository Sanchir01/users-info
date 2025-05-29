package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	api "github.com/Sanchir01/users-info/pkg/lib/api/response"
	"github.com/Sanchir01/users-info/pkg/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=UserHandlers
type UserHandlers interface {
	GetAllUsers(ctx context.Context, page, pageSize uint, minAge, maxAge *int) ([]*UserDB, error)
	DeleteUserByID(ctx context.Context, id uuid.UUID) error
	UpdateUser(ctx context.Context, id uuid.UUID, name, surname, patronymic string) error
	CreateUserService(
		name, surname, patronymic string,
		ctx context.Context,
	) error
}
type Handler struct {
	service UserHandlers
	Log     *slog.Logger
}

func NewHandler(service UserHandlers, lg *slog.Logger) *Handler {
	return &Handler{
		service: service,
		Log:     lg,
	}
}

// @Tags user
// @Description create user
// @Accept json
// @Produce json
// @Param input body CreateUserRequest true "create body"
// @Success 200 {object}  CreateUserResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /users/create [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.CreateUser"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.JSON(w, r, api.Error("Ошибка при валидации тела"))
		return
	}
	log.Info("request body decoded", slog.Any("request", req))
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	if err := h.service.CreateUserService(req.Name, req.Surname, req.Patronymic, r.Context()); err != nil {
		log.Error("fail create user", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))

		return
	}
	log.Info("success create user")

	render.JSON(w, r, CreateUserResponse{
		Response: api.OK(),
		Ok:       "user created successfully",
	})
}

// @Tags user
// @Description get all users
// @Accept json
// @Produce json
// @Param page query int false "page number" default(1) minimum(1)
// @Param page_size query int false "items per page" default(10) minimum(1) maximum(100)
// @Param min_age query int false "minimum age filter"
// @Param max_age query int false "maximum age filter"
// @Success 200 {object}  GetAllUsersResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /users [get]
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.GetAllUsers"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	minAgeStr := r.URL.Query().Get("min_age")
	maxAgeStr := r.URL.Query().Get("max_age")

	var page uint = 1
	var pageSize uint = 10
	var minAge, maxAge *int

	if pageStr != "" {
		var pageInt int
		_, err := fmt.Sscanf(pageStr, "%d", &pageInt)
		if err == nil && pageInt > 0 {
			page = uint(pageInt)
		} else {
			log.Warn("invalid page parameter", slog.String("page", pageStr))
		}
	}

	if pageSizeStr != "" {
		var pageSizeInt int
		_, err := fmt.Sscanf(pageSizeStr, "%d", &pageSizeInt)
		if err == nil && pageSizeInt > 0 && pageSizeInt <= 100 {
			pageSize = uint(pageSizeInt)
		} else {
			log.Warn("invalid page_size parameter", slog.String("page_size", pageSizeStr))
		}
	}

	if minAgeStr != "" {
		var minAgeInt int
		_, err := fmt.Sscanf(minAgeStr, "%d", &minAgeInt)
		if err == nil && minAgeInt >= 0 {
			minAge = &minAgeInt
			log.Info("filtering by min age", slog.Int("min_age", minAgeInt))
		} else {
			log.Warn("invalid min_age parameter", slog.String("min_age", minAgeStr))
		}
	}

	if maxAgeStr != "" {
		var maxAgeInt int
		_, err := fmt.Sscanf(maxAgeStr, "%d", &maxAgeInt)
		if err == nil && maxAgeInt >= 0 {
			maxAge = &maxAgeInt
			log.Info("filtering by max age", slog.Int("max_age", maxAgeInt))
		} else {
			log.Warn("invalid max_age parameter", slog.String("max_age", maxAgeStr))
		}
	}

	users, err := h.service.GetAllUsers(r.Context(), page, pageSize, minAge, maxAge)
	if err != nil {
		log.Error("fail get all users", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	logParams := []any{
		slog.Uint64("page", uint64(page)),
		slog.Uint64("page_size", uint64(pageSize)),
	}

	if minAge != nil {
		logParams = append(logParams, slog.Int("min_age", *minAge))
	}

	if maxAge != nil {
		logParams = append(logParams, slog.Int("max_age", *maxAge))
	}

	log.Info("get all users success", logParams...)

	render.JSON(w, r, GetAllUsersResponse{
		Response:     api.OK(),
		Users:        users,
		Page:         page,
		ItemsPerPage: pageSize,
	})
}

// @Tags user
// @Description delete user by id
// @Param id path string true "user id"
// @Accept json
// @Produce json
// @Success 200 {object}  DeleteUserResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.DeleteUser"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	id := chi.URLParam(r, "id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		log.Error("invalid UUID format", sl.Err(err))
		render.JSON(w, r, api.Error("invalid UUID format"))
		return
	}

	err = h.service.DeleteUserByID(r.Context(), uuidID)
	if err != nil {
		log.Error("fail get user", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	log.Info("delete user success")

	render.JSON(w, r, DeleteUserResponse{
		Response: api.OK(),
		Ok:       "success deleted",
	})
}

// @Tags user
// @Description update user by id
// @Param id path string true "user id"
// @Accept json
// @Produce json
// @Param input body UpdateUserRequest true "update body"
// @Success 200 {object} UpdateUserResponse
// @Failure 400,404 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /users/{id} [patch]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.UpdateUser"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	id := chi.URLParam(r, "id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		log.Error("invalid UUID format", sl.Err(err))
		render.JSON(w, r, api.Error("invalid UUID format"))
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.JSON(w, r, api.Error("Ошибка при валидации тела"))
		return
	}

	log.Info("request body decoded", slog.Any("request", req))

	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}

	if err := h.service.UpdateUser(r.Context(), uuidID, req.Name, req.Surname, req.Patronymic); err != nil {
		log.Error("fail update user", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}

	log.Info("update user success")

	render.JSON(w, r, UpdateUserResponse{
		Response: api.OK(),
		Ok:       "user updated successfully",
	})
}
