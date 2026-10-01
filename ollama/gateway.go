package ollama

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/ollama/ollama/api"
	"go.uber.org/fx"
)

var gatewayModule = fx.Provide(newGateway)

type gatewayParams struct {
	fx.In

	Cfg config
}

type Gateway interface {
	Chat(context.Context, *chatRequest) (*chatResponse, error)
}

func newGateway(p gatewayParams) (Gateway, error) {
	u, err := url.Parse(p.Cfg.Host)
	if err != nil {
		return nil, err
	}
	client := api.NewClient(
		u,
		&http.Client{
			Timeout: p.Cfg.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:        p.Cfg.MaxIdleConns,
				IdleConnTimeout:     p.Cfg.IdleConnTimeout,
				MaxIdleConnsPerHost: p.Cfg.MaxIdleConnsPerHost,
			},
		},
	)
	return &gateway{
		cfg:    p.Cfg,
		client: client,
	}, nil
}

type chatRequest struct {
	Messages []*message
}

type message struct {
	Role    string
	Content string
	Images  [][]byte
}

type chatResponse struct {
	Message string
}

type gateway struct {
	cfg    config
	client *api.Client
}

func (g *gateway) Chat(ctx context.Context, req *chatRequest) (*chatResponse, error) {
	messages := []api.Message{}
	for _, m := range req.Messages {
		msg := api.Message{
			Role:    m.Role,
			Content: m.Content,
		}
		images := []api.ImageData{}
		for _, bytes := range m.Images {
			images = append(images, bytes)
		}
		msg.Images = images
		messages = append(messages, msg)
	}
	ollamaReq := &api.ChatRequest{
		Model:    g.cfg.Model,
		Messages: messages,
		Stream:   &g.cfg.Stream,
		Options: map[string]any{
			"num_ctx": g.cfg.ContextSize,
		},
	}
	responseTokens := []string{}
	err := g.client.Chat(ctx, ollamaReq, func(cr api.ChatResponse) error {
		responseTokens = append(responseTokens, cr.Message.Content)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &chatResponse{
		Message: strings.Join(responseTokens, ""),
	}, nil
}
