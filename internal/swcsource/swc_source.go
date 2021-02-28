package swcsource

import (
	"fmt"
	"time"

	"github.com/renajohn/pac_collector/api"
)

// SWCSource interfaces with the SWC heat pump.
type SWCSource struct {
	WebSocketURL   string
	PollIntervalMs time.Duration // defaults to 1 min

	sessionFactory       SWCSessionFactory
	currentSession       Session
	measurementsChannel  chan api.Measurement
	sessionErrorsChannel chan error

	restartOnSessionFailure bool
}

// SWCSessionFactory generate sessions with a layer of abstraction for better testing
type SWCSessionFactory interface {
	New(source *SWCSource) Session
}

type _SWCSessionFactoryImpl struct {
}

func (factory _SWCSessionFactoryImpl) New(source *SWCSource) Session {
	session, err := newSWCSession(source.WebSocketURL, source.PollIntervalMs, source.measurementsChannel, source.sessionErrorsChannel)

	if err != nil {
		panic("Failed to create SWC session")
	}

	return session
}

// MeasurementsChannel satisfies Session interface
func (swc *SWCSource) MeasurementsChannel() <-chan api.Measurement {
	return swc.measurementsChannel
}

// NewSWCSource creates a new SWCSource
func NewSWCSource(URL string, pollingInterval time.Duration) *SWCSource {
	factory := _SWCSessionFactoryImpl{}

	return newSWCSourceWithSessionFactory(URL, pollingInterval, &factory)
}

func newSWCSourceWithSessionFactory(URL string, pollingInterval time.Duration, factory SWCSessionFactory) *SWCSource {
	source := SWCSource{
		WebSocketURL:         URL,
		measurementsChannel:  make(chan api.Measurement, 10),
		sessionErrorsChannel: make(chan error, 10),

		sessionFactory:          factory,
		restartOnSessionFailure: true,
		PollIntervalMs:          pollingInterval,
	}

	return &source
}

// Start satisfies the Source interface
func (swc *SWCSource) Start() {
	session := swc.sessionFactory.New(swc)
	swc.currentSession = session

	go session.StartSession()

	if swc.restartOnSessionFailure {
		swc.monitorSessionError()
	}
}

// The SWCSession will post an error to the session when ever something bad happens.
func (swc *SWCSource) monitorSessionError() {
	err := <-swc.sessionErrorsChannel
	if err != nil {
		fmt.Println(fmt.Sprintf("Session was terminited due to an error, restarting: %v", err))
		swc.Start()
	}
}
