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

	ws                  *websocket.Conn
	measurementsChannel chan api.Measurement
	errorsChannel       chan error
}

// NewSWCSource construct and validate an SWC source
func NewSWCSource(URL string, pollingInterval time.Duration) (*SWCSource, error) {
	source := SWCSource{
		WebSocketURL:        URL,
		measurementsChannel: make(chan api.Measurement, 10),
		errorsChannel:       make(chan error, 10),
	}

	if pollingInterval == 0 {
		source.PollIntervalMs = time.Minute
	}

	err := source.validateConfig()

	return &source, err
}

func (swc *SWCSource) validateConfig() error {
	if !strings.HasPrefix(swc.WebSocketURL, "ws:") {
		return errors.New("SWC Source URL must be a well formatted WebSocket URL")
	}

	return nil
}

// MeasurementsChannel satisfies Source interface
func (swc *SWCSource) MeasurementsChannel() <-chan api.Measurement {
	return swc.measurementsChannel
}

// ErrorsChannel satisfies Source interface
func (swc *SWCSource) ErrorsChannel() <-chan error {
	return swc.errorsChannel
}

// Start satisfies Source interface
func (swc *SWCSource) Start() {
	swc.measurementsChannel = make(chan api.Measurement)
	swc.errorsChannel = make(chan error)

	err := swc.validateConfig()
	if err != nil {
		swc.errorsChannel <- err
	}

	swc.goStart()
}

func (swc *SWCSource) goStart() {
	var err error

	err = swc.connect()
	if err != nil {
		swc.errorsChannel <- err
		return
	}

	err = swc.login()
	if err != nil {
		swc.errorsChannel <- err
		return
	}

	err = swc.getTemperatures()
	if err != nil {
		swc.errorsChannel <- err
		return
	}

	go swc.readMessages()
	go swc.poll()
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
			readError := fmt.Sprintf("Error while reading WS message %g, aborting", err)
			log.Println(readError)
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
			swc.measurementsChannel <- api.Measurement{
				MeasurementType: api.WaterTemperature,
				Timestamp:       time.Now().Unix(),
				Value:           data}
		} else {
			log.Printf("Failed to parse XML message %v", err)
		}
	}
}

func (swc *SWCSource) poll() {
	var err error
	for {
		time.Sleep(swc.PollIntervalMs)

		err = swc.ws.WriteMessage(websocket.TextMessage, []byte("REFRESH"))
		if err != nil {
			pollError := fmt.Sprintf("Error while polling for data %g, aborting", err)
			log.Println(pollError)
			swc.errorsChannel <- errors.New(pollError)
			return
		}
	}
}
