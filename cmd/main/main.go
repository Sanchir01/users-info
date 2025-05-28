package main

import (
	"context"
	"errors"
	"github.com/Sanchir01/users-info/internal/app"
	httpserver "github.com/Sanchir01/users-info/internal/server/http"
	httphandlers "github.com/Sanchir01/users-info/internal/server/http/handlers"
	"github.com/fatih/color"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title 🚀 Effective Mobile
// @version         1.0
// @description This is a sample server celler
// @termsOfService  http://swagger.io/terms/

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @contact.name GitHub
// @contact.url https://github.com/Sanchir01
func main() {
	env, err := app.NewEnv()
	if err != nil {
		panic(err)
	}
	serverrest := httpserver.NewHTTPServer(env.Config.Host, env.Config.Port,
		env.Config.Timeout, env.Config.IdleTimeout)
	prometheusserver := httpserver.NewHTTPServer(env.Config.Prometheus.Host, env.Config.Prometheus.Port, env.Config.Prometheus.Timeout,
		env.Config.Prometheus.IdleTimeout)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()
	green := color.New(color.FgGreen).SprintFunc()

	env.Logger.Info(green("🚀 Server started successfully!"),
		slog.String("time", time.Now().Format("2006-01-02 15:04:05")),
		slog.String("port", env.Config.Port),
	)

	go func() {
		if err := serverrest.Run(httphandlers.StartHTTTPHandlers(env.Handlers)); err != nil {
			if !errors.Is(err, context.Canceled) {
				env.Logger.Error("Listen server error", slog.String("error", err.Error()))
				return
			}
		}
	}()
	go func() {
		if err := prometheusserver.Run(httphandlers.StartPrometheusHandlers()); err != nil {
			if !errors.Is(err, context.Canceled) {
				env.Logger.Error("Listen prometheus server error", slog.String("error", err.Error()))
				return
			}
		}
	}()
	<-ctx.Done()
	if err := serverrest.Gracefull(ctx); err != nil {
		env.Logger.Error("server gracefull")
	}
	if err := env.DataBase.Close(); err != nil {
		env.Logger.Error("Close database", slog.String("error", err.Error()))
	}
}
