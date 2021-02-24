package swcsource

import (
	"errors"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/renajohn/pac_collector/api"
)

// SWCSource interfaces with the SWC heat pump.
// https://www.alpha-innotec.ch/alpha-innotec/produits/pompes-a-chaleur/soleau/swc-82k3.html?L=2
type SWCSource struct {
	Pin   string
	wsURL string
	ws    *websocket.Conn
}

// Start satisfies Source interface
func (swc *SWCSource) Start() (<-chan api.Measurement, error) {
	connectionError := swc.connect()

	return make(chan api.Measurement), connectionError
}

func (swc *SWCSource) connect() error {
	if len(swc.Pin) == 0 {
		return errors.New("No Pin provided for SWC source")
	}
	if !strings.HasPrefix(swc.wsURL, "ws:") {
		return errors.New("SWC Source URL should be a well formatted WebSocket URL")
	}

	ws, _, err := websocket.DefaultDialer.Dial(swc.wsURL, nil)
	if err != nil {
		return err
	}
	swc.ws = ws

	return swc.ws.WriteMessage(websocket.TextMessage, []byte("LOGIN;"+swc.Pin))
}
