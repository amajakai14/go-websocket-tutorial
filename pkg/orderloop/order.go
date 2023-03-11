package orderloop

import (
	"github.com/amajakai14/go-websocket-tutorial/pkg/server"
)

type BuffetChannel struct {
	ChannelId int
	sockets []Socket
}
