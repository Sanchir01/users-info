package user

import (
	"encoding/json"
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

type Handler struct {
	service *Service
	Log     *slog.Logger
}

func NewHandler(service *Service, lg *slog.Logger) *Handler {
	return &Handler{
		service: service,
		Log:     lg,
	}
}

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
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.GetAllUsers"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		log.Error("fail get all users", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	log.Info("get all users success")

	render.JSON(w, r, GetAllUsersResponse{
		Response: api.OK(),
		Users:    users,
	})
}

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
