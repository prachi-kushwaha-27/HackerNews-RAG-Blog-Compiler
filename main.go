package main

import (
	"github.com/prachi-kushwaha-27/tldr/hackernews"
	httpapp "github.com/prachi-kushwaha-27/tldr/http/app"
	httphandler "github.com/prachi-kushwaha-27/tldr/http/handler"
	"github.com/prachi-kushwaha-27/tldr/ollama"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		ollama.Module,
		hackernews.Module,
		httpapp.Module,
		fx.Invoke(func(h httphandler.Handler) {}),
		fx.Invoke(func(h hackernews.Cron) {}),
	).Run()
}
