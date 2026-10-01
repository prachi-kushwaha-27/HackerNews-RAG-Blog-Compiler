package hackernews

import (
	"context"
	"fmt"
	"slices"

	"go.uber.org/fx"
)

var controllerModule = fx.Provide(newController)

type Controller interface {
	GetSummary(ctx context.Context, date string, postId int64) (*Summary, error)
	ListPosts(ctx context.Context, date string, tags []string) (*Posts, error)
	ListAvailablePostsDates(ctx context.Context) (*PostDates, error)
}

type controllerParams struct {
	fx.In

	Cfg config
}

func newController(p controllerParams) Controller {
	return &controller{
		cfg: p.Cfg,
	}
}

type controller struct {
	cfg config
}

func (c *controller) GetSummary(ctx context.Context, date string, postId int64) (*Summary, error) {
	posts, err := readPosts(date, c.cfg)
	if err != nil {
		return nil, err
	}
	for _, post := range posts.Posts {
		if post.ID == postId {
			postSummaryBrr, err := readPostSummary(postId, c.cfg)
			if err != nil {
				return nil, err
			}
			postSummary := cleanSummary(string(postSummaryBrr))
			summary := Summary{
				ID:             post.ID,
				Title:          extractPostTitle(postSummary),
				Summary:        postSummary,
				DateSummarised: post.Datetime,
				URL:            post.URL,
				Tags:           post.Tags,
			}
			return &summary, nil
		}
	}
	return nil, fmt.Errorf("post not found for date: %s; ID: %d", date, postId)
}

func (c *controller) ListPosts(ctx context.Context, date string, tags []string) (*Posts, error) {
	posts, err := readPosts(date, c.cfg)
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return posts, nil
	}
	postsFinal := []Post{}
	for _, p := range posts.Posts {
		for _, t := range tags {
			if slices.Contains(p.Tags, t) {
				postsFinal = append(postsFinal, p)
				break
			}
		}
	}
	posts.Posts = postsFinal
	return posts, nil
}

func (c *controller) ListAvailablePostsDates(ctx context.Context) (*PostDates, error) {
	dates, err := listAvailablePostsDates(c.cfg)
	return &PostDates{
		Dates: dates,
	}, err
}
