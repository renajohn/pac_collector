package swcsource

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/renajohn/pac_collector/api"
)

// SWCSource interfaces with the SWC heat pump.
// https://www.alpha-innotec.ch/alpha-innotec/produits/pompes-a-chaleur/soleau/swc-82k3.html?L=2
type SWCSource struct {
	WebSocketURL   string
	PollIntervalMs time.Duration // defaults to 1 min

	ws              *websocket.Conn
	messagesChannel chan api.Measurement
}

// Start satisfies Source interface
func (swc *SWCSource) Start() (<-chan api.Measurement, error) {
	var err error

	err = swc.validateConfig()
	if err != nil {
		return nil, err
	}

	swc.messagesChannel = make(chan api.Measurement)
	return swc.messagesChannel, swc.goStart()
}

func (swc *SWCSource) validateConfig() error {
	if !strings.HasPrefix(swc.WebSocketURL, "ws:") {
		return errors.New("SWC Source URL should be a well formatted WebSocket URL")
	}

	if swc.PollIntervalMs == 0 {
		swc.PollIntervalMs = time.Minute
	}
	return nil
}

func (swc *SWCSource) goStart() error {
	var err error

	err = swc.connect()
	if err != nil {
		return err
	}

	err = swc.login()
	if err != nil {
		return err
	}

	err = swc.getTemperatures()
	if err != nil {
		return err
	}

	go swc.readMessages()
	go swc.poll()

	return nil
}

func (swc *SWCSource) connect() error {
	ws, _, err := websocket.DefaultDialer.Dial(swc.WebSocketURL, nil)
	if err != nil {
		return err
	}
	swc.ws = ws

	return nil
}

func (swc *SWCSource) login() error {
	return swc.ws.WriteMessage(websocket.TextMessage, []byte("LOGIN;000000"))
}

func (swc *SWCSource) getTemperatures() error {
	return swc.ws.WriteMessage(websocket.TextMessage, []byte("GET;0x46bd50"))
}

func (swc *SWCSource) readMessages() {
	for {
		// receive message
		messageType, message, err := swc.ws.ReadMessage()
		if err != nil {
			// handle error
			log.Printf("Error while reading WS message %v\n", err)
			return
		} else if messageType == websocket.TextMessage {
			swc.parseMessage(message)
		}
	}
}

func (swc *SWCSource) parseMessage(byteXML []byte) {
	if strings.HasPrefix(string(byteXML), "<values>") {
		values, err := parseXMLMeasurement(byteXML)

		if err == nil {
			data, _ := json.Marshal(values)
			swc.messagesChannel <- api.Measurement{
				MeasurementType: api.WaterTemperature,
				Timestamp:       time.Now().Unix(),
				Value:           data}
		} else {
			log.Printf("Failed to parse XML message %v", err)
		}
	}
}

func (swc *SWCSource) poll() error {
	var err error
	for {
		time.Sleep(swc.PollIntervalMs)

		err = swc.ws.WriteMessage(websocket.TextMessage, []byte("REFRESH"))
		if err != nil {
			fmt.Printf("Error while polling for data %g, aborting", err)
			return err
		}
	}
}
