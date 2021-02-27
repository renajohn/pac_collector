package swcsource

import (
	"testing"
)

type MockSessionFactory struct {
	nbCalled int
}

type MockSession struct {
}

func (ms *MockSession) StartSession() {
}

func (sf *MockSessionFactory) New(source *SWCSource) Session {
	sf.nbCalled++
	session := MockSession{}
	return &session
}

func TestStart(t *testing.T) {
	factory := MockSessionFactory{}
	source := swcSourceWithSessionFactory("ws:testurl", 1000, &factory)

	source.Start()

	if factory.nbCalled != 1 {
		t.Errorf("Expected 1 session to be created, got %d", factory.nbCalled)
	}
}
