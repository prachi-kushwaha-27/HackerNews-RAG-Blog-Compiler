package app

import (
	"github.com/prachi-kushwaha-27/tldr/http/config"
	http "github.com/prachi-kushwaha-27/tldr/http/handler"
	"go.uber.org/fx"
)

var Module = fx.Options(
	config.Module,
	http.Module,
)
