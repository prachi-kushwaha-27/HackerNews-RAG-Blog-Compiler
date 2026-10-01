package config

import (
	"os"

	"github.com/goccy/go-yaml"
	"go.uber.org/fx"
)

var Module = fx.Provide(ParseConfig)

const __filePath = "./http/config/config.yaml"

func ParseConfig() (Config, error) {
	cfg := Config{}
	bytes, err := os.ReadFile(__filePath)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(bytes, &cfg)
	return cfg, err
}

type Config struct {
	Server Server `yaml:"server"`
}

type Server struct {
	HtmlDir string `yaml:"htmlDir"`
	Http    struct {
		Address string `yaml:"address"`
	} `yaml:"http"`
}
