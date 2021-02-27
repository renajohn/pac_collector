package swcsource

import (
	"errors"
	"strings"
	"time"

	"github.com/renajohn/pac_collector/api"
)

// SWCSource interfaces with the SWC heat pump.
type SWCSource struct {
	WebSocketURL   string
	PollIntervalMs time.Duration // defaults to 1 min

	sessionFactory      SWCSessionFactory
	currentSession      Session
	measurementsChannel chan api.Measurement
	errorsChannel       chan error
}

// SWCSessionFactory generate sessions with a layer of abstraction for better testing
type SWCSessionFactory interface {
	New(source *SWCSource) Session
}

type _SWCSessionFactoryImpl struct {
}

func (factory _SWCSessionFactoryImpl) New(source *SWCSource) Session {
	session, _ := NewSWCSession(source.WebSocketURL, source.PollIntervalMs, source.measurementsChannel, source.errorsChannel)
	return session
}

func (swc *SWCSource) validateConfig() error {
	if !strings.HasPrefix(swc.WebSocketURL, "ws:") {
		return errors.New("SWC Session URL must be a well formatted WebSocket URL")
	}

	return nil
}

// MeasurementsChannel satisfies Session interface
func (swc *SWCSource) MeasurementsChannel() <-chan api.Measurement {
	return swc.measurementsChannel
}

// ErrorsChannel satisfies Session interface
func (swc *SWCSource) ErrorsChannel() <-chan error {
	return swc.errorsChannel
}

func swcSourceWithSessionFactory(URL string, pollingInterval time.Duration, factory SWCSessionFactory) *SWCSource {
	source := SWCSource{
		WebSocketURL:        URL,
		measurementsChannel: make(chan api.Measurement, 10),
		errorsChannel:       make(chan error, 10),

		sessionFactory: factory,
	}

	if pollingInterval == 0 {
		source.PollIntervalMs = time.Minute
	}

	err := source.validateConfig()
	if err != nil {
		panic("Wrong input parameters, aborting!")
	}

	return &source
}

// Start satisfies the Source interface
func (swc *SWCSource) Start() {
	session := swc.sessionFactory.New(swc)
	swc.currentSession = session

	session.StartSession()
}
