package ollama

import (
	"context"
	"encoding/base64"

	"go.uber.org/fx"
)

var controllerModule = fx.Provide(newController)

type Controller interface {
	Instruct(context.Context, *InstructRequest) (*InstructResponse, error)
}

type controllerParams struct {
	fx.In

	Model Gateway
}

func newController(p controllerParams) Controller {
	return &controller{
		model: p.Model,
	}
}

type controller struct {
	model Gateway
}

func (c *controller) Instruct(ctx context.Context, req *InstructRequest) (*InstructResponse, error) {
	images := [][]byte{}
	for _, b64 := range req.ImagesB64 {
		bytes, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return nil, err
		}
		images = append(images, bytes)
	}
	chatResp, err := c.model.Chat(ctx, &chatRequest{
		Messages: []*message{
			{
				Role:    roleSystem,
				Content: req.Instruction,
			},
			{
				Role:    roleUser,
				Content: req.Content,
				Images:  images,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &InstructResponse{Result: chatResp.Message}, nil
}
