package ollama

import (
	"os"
	"time"

	"github.com/goccy/go-yaml"
	"go.uber.org/fx"
)

var configModule = fx.Provide(parseConfig)

const __filePath = "./ollama/config.yaml"

func parseConfig() (config, error) {
	cfg := config{}
	bytes, err := os.ReadFile(__filePath)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(bytes, &cfg)
	return cfg, err
}

type config struct {
	Host                string        `yaml:"host"`
	Model               string        `yaml:"model"`
	ContextSize         int           `yaml:"contextSize"`
	Timeout             time.Duration `yaml:"timeout"`
	Stream              bool          `yaml:"stream"`
	MaxIdleConns        int           `yaml:"maxIdleConns"`
	IdleConnTimeout     time.Duration `yaml:"idleConnTimeout"`
	MaxIdleConnsPerHost int           `yaml:"maxIdleConnsPerHost"`
}
