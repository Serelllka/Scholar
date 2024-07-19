package client

import "Scholar/internal/pkg/client/model"

type IWebClient interface {
	SendMessage(msgType model.MessageType, payload []byte) error
	GetId() string
	StartReceiver()
}
