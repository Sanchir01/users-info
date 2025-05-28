package app

import (
	"fmt"

	"log/slog"

	"github.com/Sanchir01/users-info/internal/config"
)

type Env struct {
	Logger   *slog.Logger
	DataBase *Database

	Config       *config.Config
	Handlers     *Handlers
	Repositories *Repositories
	Services     *Services
}

func NewEnv() (*Env, error) {
	cfg := config.InitConfig()
	fmt.Println(cfg.RedisDB)
	lg := setupLogger(cfg.Env)

	pgxdb, err := NewDataBases(cfg)
	if err != nil {
		lg.Error("pgx error connect", err.Error())
		return nil, err
	}

	repos := NewRepositories(pgxdb)
	servises := NewServices(repos, pgxdb)
	handlers := NewHandlers(servises, lg)

	env := Env{
		Logger:       lg,
		DataBase:     pgxdb,
		Config:       cfg,
		Handlers:     handlers,
		Services:     servises,
		Repositories: repos,
	}

	return &env, nil
}
