package hackernews

import (
	"os"
	"time"

	"github.com/goccy/go-yaml"
	"go.uber.org/fx"
)

var configModule = fx.Provide(parseConfig)

const __filePath = "./hackernews/config.yaml"

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
	HNApiMaxRps       int           `yaml:"hnApiMaxRps"`
	Timeout           time.Duration `yaml:"timeout"`
	Sleep             time.Duration `yaml:"sleep"`
	StoryQueueSize    int           `yaml:"storyQueueSize"`
	StoriesSubsetSize int           `yaml:"storiesSubsetSize"`
	TempDir           string        `yaml:"tempDir"`
	PostSummaryDir    string        `yaml:"postSummaryDir"`
	PostsDir          string        `yaml:"postsDir"`
	SummarisePrompt   string        `yaml:"summarisePrompt"`
	TagPrompt         string        `yaml:"tagPrompt"`
	Tags              []string      `yaml:"tags"`
}
