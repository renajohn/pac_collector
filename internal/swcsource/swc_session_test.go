package swcsource

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/renajohn/pac_collector/api"
)

func TestStartSession(t *testing.T) {
	loginXML := readFixture(t, "testdata/LOGIN;123456.xml")
	getXML := readFixture(t, "testdata/GET;0x46bd50.xml")
	refreshXML := readFixture(t, "testdata/GET;0x46bd50-REFRESH.xml")

	t.Run("WebSocket connection", func(t *testing.T) {
		t.Skip()
		var spy = wsSpy{}
		handler := generateHTTPHandler(t, []messageResponse{{
			message:  "LOGIN;000000",
			response: loginXML,
		}}, &spy)

		server := httptest.NewServer(http.HandlerFunc(handler))

		measurementsChannel := make(chan api.Measurement)
		errorsChannel := make(chan error)
		session, err := NewSWCSession(toWs(server.URL), 1000, measurementsChannel, errorsChannel)
		if err != nil {
			t.Errorf("Failed to create new SWC session: %v", session)
		}

		session.StartSession()

		server.Close()

		if len(spy.messagesLog) < 1 {
			t.Errorf("expected at least 1 message in WS got %d", len(spy.messagesLog))
		}
	})

	t.Run("session should be polling for new measurements", func(t *testing.T) {
		t.Skip()
		var spy = wsSpy{}
		handler := generateHTTPHandler(t, []messageResponse{{
			message:  "LOGIN;000000",
			response: loginXML,
		}, {
			message:  "GET;0x46bd50",
			response: getXML,
		}, {
			message:  "REFRESH",
			response: refreshXML,
		}, {
			message:  "REFRESH",
			response: refreshXML,
		}, {
			message:  "REFRESH",
			response: refreshXML,
		}}, &spy)

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		measurementsChannel := make(chan api.Measurement)
		errorsChannel := make(chan error)
		session, _ := NewSWCSession(toWs(server.URL), 3, measurementsChannel, errorsChannel)

		go session.StartSession()

		expectedValue, _ := json.Marshal(SWCMeasurement{
			HeatingOutboundTemperature:     33.8,
			HeatingInboundTemperature:      34.5,
			OutsideTemperature:             4.7,
			TankTemperature:                52.3,
			TargetTankTemperature:          52.0,
			DrillInboundTemperature:        11.0,
			DrillOutboundTemperature:       11.2,
			AmbiantIndoorTemperature:       21.1,
			AmbiantIndoorTargetTemperature: 21.0,
		})

		for index := 0; index < 3; index++ {
			measurement := <-session.MeasurementsChannel

			if string(expectedValue) != string(measurement.Value) {
				t.Errorf("Expected value of %v got %v", string(expectedValue), string(measurement.Value))
			}
		}

		if len(session.ErrorsChannel) > 0 {
			t.Errorf("Expected error in SWC session: %v", <-session.ErrorsChannel)
		}
	})

	t.Run("When connection drops, session should reconnect", func(t *testing.T) {
		var spy = wsSpy{}
		handler := generateHTTPHandler(t, []messageResponse{{
			message:  "LOGIN;000000",
			response: loginXML,
		}, {
			message:  "GET;0x46bd50",
			response: getXML,
		}, {
			message:  "REFRESH",
			response: refreshXML,
		}}, &spy)

		server := httptest.NewUnstartedServer(http.HandlerFunc(handler))
		go server.Start()
		defer server.Close()

		for len(server.URL) == 0 {
			time.Sleep(10)
		}

		measurementsChannel := make(chan api.Measurement)
		errorsChannel := make(chan error)
		session, _ := NewSWCSession(toWs(server.URL), 1000, measurementsChannel, errorsChannel)

		go session.StartSession()

		// wait for at least one measurement
		measurement := <-session.MeasurementsChannel
		if measurement.MeasurementType != api.WaterTemperature {
			t.Errorf("Expected measurement type %s, but got %s", api.WaterTemperature, measurement.MeasurementType)
		}

		server.CloseClientConnections()

		// test will fail if no errors is returned due to the test timeout
		<-session.ErrorsChannel
	})

	t.Run("When URL is not valid, an error should be generated", func(t *testing.T) {
		measurementsChannel := make(chan api.Measurement)
		errorsChannel := make(chan error)
		_, err := NewSWCSession("http://when-it-should-be-ws", 1000, measurementsChannel, errorsChannel)

		if err == nil {
			t.Error("An error was expected and none was returned")
		}
	})

	t.Run("When no polling interval is provided, use 1min", func(t *testing.T) {
		measurementsChannel := make(chan api.Measurement)
		errorsChannel := make(chan error)
		session, err := NewSWCSession("ws://when-it-should-be-ws", 0, measurementsChannel, errorsChannel)

		if err != nil {
			t.Errorf("No error was expected but got %v", err)
		}

		if session.PollIntervalMs != time.Millisecond*1000*60 {
			t.Errorf("Expected polling intervals of %dms and got %dms", time.Millisecond*1000*60, session.PollIntervalMs)
		}
	})
}
