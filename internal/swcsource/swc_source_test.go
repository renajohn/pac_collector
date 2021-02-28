package swcsource

import (
	"errors"
	"testing"
	"time"
)

type MockSessionFactory struct {
	nbCalled int
}

func (sf *MockSessionFactory) New(source *SWCSource) Session {
	sf.nbCalled++
	session := MockSession{}
	return &session
}

type MockSession struct {
}

func (ms *MockSession) StartSession() {
}

func TestStart(t *testing.T) {

	t.Run("Happy case", func(t *testing.T) {
		factory := MockSessionFactory{}
		source := newSWCSourceWithSessionFactory("ws:testurl", 1000, &factory)
		source.restartOnSessionFailure = false

		source.Start()

		if factory.nbCalled != 1 {
			t.Errorf("Expected 1 session to be created, got %d", factory.nbCalled)
		}
	})

	t.Run("When session fails, a new one is created", func(t *testing.T) {
		factory := MockSessionFactory{}
		source := newSWCSourceWithSessionFactory("ws:testurl", 1000, &factory)

		go source.Start()

		// make sure the Start has been executed
		time.Sleep(100 * time.Microsecond)

		source.sessionErrorsChannel <- errors.New("Boom")

		// make sure the Start has been executed
		time.Sleep(100 * time.Microsecond)

		if factory.nbCalled != 2 {
			t.Errorf("Expected 2 session to be created, got %d", factory.nbCalled)
		}
	})

}
