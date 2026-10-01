package ollama

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	configModule,
	gatewayModule,
	controllerModule,
)
