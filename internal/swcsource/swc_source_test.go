package swcsource

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type MessageResponse struct {
	message  string
	response string
}

type WsSpy struct {
	messagesLog []string
}

func assertValue(t *testing.T, expect string, got string) {
	t.Helper()
	if expect != got {
		t.Errorf("Expected value of %s got %s", expect, got)
	}
}

func generateHTTPHandler(t *testing.T, expected []MessageResponse, spy *WsSpy) func(http.ResponseWriter, *http.Request) {
	t.Helper()
	return func(response http.ResponseWriter, request *http.Request) {
		t.Helper()
		connection, err := upgrader.Upgrade(response, request, nil)
		if err != nil {
			return
		}
		defer connection.Close()
		for _, messageResponse := range expected {
			mt, message, err := connection.ReadMessage()
			if err != nil {
				t.Fatalf("ReadMessage failed %v", err)
				break
			}
			stringMessage := string(message)
			spy.messagesLog = append(spy.messagesLog, stringMessage)
			assertValue(t, messageResponse.message, stringMessage)

			err = connection.WriteMessage(mt, []byte(messageResponse.response))
			if err != nil {
				t.Fatalf("WriteMessage failed %v", err)
				break
			}
		}
	}
}

func readFixture(t *testing.T, fileName string) string {
	t.Helper()

	data, readError := os.ReadFile(fileName)
	if readError != nil {
		t.Fatalf("failed to read file %v", readError)
	}

	return string(data)
}

func toWs(url string) string {
	return "ws" + strings.TrimPrefix(url, "http")
}

func TestStart(t *testing.T) {
	loginXML := readFixture(t, "testdata/LOGIN;123456.xml")
	getXML := readFixture(t, "testdata/GET;0x46bd50.xml")
	refreshXML := readFixture(t, "testdata/GET;0x46bd50-REFRESH.xml")

	t.Run("WebSocket connection", func(t *testing.T) {
		var spy = WsSpy{}
		handler := generateHTTPHandler(t, []MessageResponse{{
			message:  "LOGIN;000000",
			response: loginXML,
		}}, &spy)

		server := httptest.NewServer(http.HandlerFunc(handler))
		source, err := NewSWCSource(toWs(server.URL), 1000)
		if err != nil {
			t.Errorf("Failed to create new SWC source: %v", source)
		}

		source.Start()

		server.Close()

		if len(spy.messagesLog) < 1 {
			t.Errorf("expected at least 1 message in WS got %d", len(spy.messagesLog))
		}
	})

	t.Run("Source should be polling for new measurements", func(t *testing.T) {
		var spy = WsSpy{}
		handler := generateHTTPHandler(t, []MessageResponse{{
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

		source, err := NewSWCSource(toWs(server.URL), 3)
		if err != nil {
			t.Errorf("Failed to create new SWC source: %v", source)
		}

		go source.Start()

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
			measurement := <-source.MeasurementsChannel()

			if string(expectedValue) != string(measurement.Value) {
				t.Errorf("Expected value of %v got %v", string(expectedValue), string(measurement.Value))
			}
		}

		if len(source.ErrorsChannel()) > 0 {
			t.Errorf("Expected error in SWC source: %v", <-source.ErrorsChannel())
		}
	})

	// Test to write: Connection fails
	t.Run("When connection drops, source should reconnect", func(t *testing.T) {

	})
}
