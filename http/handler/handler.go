package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/prachi-kushwaha-27/tldr/hackernews"
	"github.com/prachi-kushwaha-27/tldr/http/config"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewHttpHandler)

type Params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Config     config.Config
	HackerNews hackernews.Controller
}

type Handler interface {
	Run()
	Shutdown(context.Context)
}

func NewHttpHandler(p Params) Handler {
	mux := http.NewServeMux()
	addRoutes(mux, p.HackerNews, p.Config)

	httpHandler := addMiddlewares(mux)

	server := &http.Server{
		Addr:    p.Config.Server.Http.Address,
		Handler: httpHandler,
	}
	h := &handler{
		server: server,
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			h.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			h.Shutdown(ctx)
			return nil
		},
	})

	return h
}

type handler struct {
	server *http.Server
}

func (h *handler) Run() {
	go func() {
		slog.Info("http server listening", slog.String("address", h.server.Addr))
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed listening and serving", slog.Any("error", err))
		}
	}()
}

func (h *handler) Shutdown(ctx context.Context) {
	if err := h.server.Shutdown(ctx); err != nil {
		slog.Error("failed gracefully shutting down http server", slog.Any("error", err))
	} else {
		slog.Info("graceful shutdown success")
	}
}
