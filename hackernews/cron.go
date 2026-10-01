package hackernews

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/prachi-kushwaha-27/tldr/lib/gohn"
	"go.uber.org/fx"
)

var cronModule = fx.Provide(newCronJob)

type Cron interface {
	Run()
	Shutdown()
}

type cronParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Cfg       config
	HnClient  gohn.Client
	Pipeline  Pipeline
}

func newCronJob(p cronParams) Cron {
	job := &cron{
		cfg:      p.Cfg,
		hnClient: p.HnClient,
		pipeline: p.Pipeline,
		done:     make(chan struct{}),
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			job.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			job.Shutdown()
			return nil
		},
	})

	return job
}

type cron struct {
	cfg      config
	hnClient gohn.Client
	pipeline Pipeline

	done chan struct{}
}

func (c *cron) Run() {
	go c.processTopAndBestStories()
}

func (c *cron) Shutdown() {
	c.done <- struct{}{}
	close(c.done)
}

func (c *cron) processTopAndBestStories() {
	for {
		slog.Info("getting best stories")
		bestStories, err := c.hnClient.GetBestStories()
		if err != nil {
			slog.Error("failed getting best stories", slog.Any("error", err))
		} else {
			c.processStories(*bestStories)
		}

		// slog.Info("getting top stories")
		// topStories, err := c.hnClient.GetTopStories()
		// if err != nil {
		// 	slog.Error("failed getting top stories", slog.Any("error", err))
		// } else {
		// 	c.processStories(*topStories)
		// }

		slog.Info("top stories waiting", slog.String("for", c.cfg.Sleep.String()))
		select {
		case <-c.done:
			return
		case <-time.Tick(c.cfg.Sleep):
			continue
		}
	}
}

func (c *cron) processStories(stories gohn.Stories) {
	subsetSize := min(c.cfg.StoriesSubsetSize, len(stories))
	slog.Info(fmt.Sprintf("found %d stories, processing top %d", len(stories), subsetSize))
	subset := gohn.Stories([]int64(stories)[:subsetSize])
	for _, id := range subset {
		c.pipeline.Push(id)
	}
}
