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

// SWCSession interfaces with the SWC heat pump.
// https://www.alpha-innotec.ch/alpha-innotec/produits/pompes-a-chaleur/soleau/swc-82k3.html?L=2
type SWCSession struct {
	WebSocketURL   string
	PollIntervalMs time.Duration // defaults to 1 min

	MeasurementsChannel chan api.Measurement
	ErrorsChannel       chan error

	ws *websocket.Conn
}

// Session represent a Web Socket session. If the connection is broken, the session is destroyed
type Session interface {
	StartSession()
}

// NewSWCSession construct and validate an SWC Session
func NewSWCSession(URL string, pollingInterval time.Duration, measurementsChannel chan api.Measurement, errorsChannel chan error) (*SWCSession, error) {
	session := SWCSession{
		WebSocketURL:        URL,
		MeasurementsChannel: measurementsChannel,
		ErrorsChannel:       errorsChannel,
	}

	if pollingInterval == 0 {
		session.PollIntervalMs = time.Minute
	}

	err := session.validateConfig()

	return &session, err
}

func (swc *SWCSession) validateConfig() error {
	if !strings.HasPrefix(swc.WebSocketURL, "ws:") {
		return errors.New("SWC Session URL must be a well formatted WebSocket URL")
	}

	return nil
}

// StartSession satisfies Session interface
func (swc *SWCSession) StartSession() {
	var err error

	err = swc.connect()
	if err != nil {
		swc.ErrorsChannel <- err
		return
	}

	err = swc.login()
	if err != nil {
		swc.ErrorsChannel <- err
		return
	}

	err = swc.getTemperatures()
	if err != nil {
		swc.ErrorsChannel <- err
		return
	}

	go swc.readMessages()
	go swc.poll()
}

func (swc *SWCSession) connect() error {
	ws, _, err := websocket.DefaultDialer.Dial(swc.WebSocketURL, nil)
	if err != nil {
		return err
	}
	swc.ws = ws

	return nil
}

func (swc *SWCSession) login() error {
	return swc.ws.WriteMessage(websocket.TextMessage, []byte("LOGIN;000000"))
}

func (swc *SWCSession) getTemperatures() error {
	return swc.ws.WriteMessage(websocket.TextMessage, []byte("GET;0x46bd50"))
}

func (swc *SWCSession) readMessages() {
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

func (swc *SWCSession) parseMessage(byteXML []byte) {
	if strings.HasPrefix(string(byteXML), "<values>") {
		values, err := parseXMLMeasurement(byteXML)

		if err == nil {
			data, _ := json.Marshal(values)
			swc.MeasurementsChannel <- api.Measurement{
				MeasurementType: api.WaterTemperature,
				Timestamp:       time.Now().Unix(),
				Value:           data}
		} else {
			log.Printf("Failed to parse XML message %v", err)
		}
	}
}

func (swc *SWCSession) poll() {
	var err error
	for {
		time.Sleep(swc.PollIntervalMs)

		err = swc.ws.WriteMessage(websocket.TextMessage, []byte("REFRESH"))
		if err != nil {
			pollError := fmt.Sprintf("Error while polling for data %v, aborting", err)
			log.Println(pollError)
			swc.ErrorsChannel <- errors.New(pollError)
			return
		}
	}
}
