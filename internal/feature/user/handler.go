package user

import (
	"encoding/json"
	"log/slog"
	"net/http"

	api "github.com/Sanchir01/users-info/pkg/lib/api/response"
	"github.com/Sanchir01/users-info/pkg/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
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
