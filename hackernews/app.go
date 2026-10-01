package hackernews

import (
	"github.com/prachi-kushwaha-27/tldr/lib/gohn"
	"go.uber.org/fx"
)

var Module = fx.Options(
	configModule,
	hnClientModule,
	pipelineModule,
	cronModule,
	controllerModule,
)

var hnClientModule = fx.Provide(func(cfg config) gohn.Client {
	return gohn.NewClient(cfg.Timeout, cfg.HNApiMaxRps)
})
