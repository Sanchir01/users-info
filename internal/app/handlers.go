package app

import (
	"github.com/Sanchir01/users-info/internal/feature/user"
	"log/slog"
)

type Handlers struct {
	UserHandler *user.Handler
}

func NewHandlers(services *Services, lg *slog.Logger) *Handlers {
	return &Handlers{
		UserHandler: user.NewHandler(services.UserService, lg),
	}
}
