package client

import (
	"Scholar/internal/pkg/client/model"
	"context"
)

type IWebClient interface {
	SendMessage(ctx context.Context, msgType model.MessageType, payload []byte) error
	GetId() string
	StartReceiver(ctx context.Context)
}
