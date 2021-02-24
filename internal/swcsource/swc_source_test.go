package swcsource

import (
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

func mockServer(t *testing.T, expected []MessageResponse, spy *WsSpy) func(http.ResponseWriter, *http.Request) {
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

func TestStart(t *testing.T) {
	file, readError := os.ReadFile("testdata/LOGIN;123456.xml")
	if readError != nil {
		t.Fatalf("failed to read file %v", readError)
	}

	var spy = WsSpy{}
	handler := mockServer(t, []MessageResponse{{
		message:  "LOGIN;123456",
		response: string(file),
	}}, &spy)
	server := httptest.NewServer(http.HandlerFunc(handler))

	// Convert http://127.0.0.1 to ws://127.0.0.
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	source := SWCSource{Pin: "123456", wsURL: wsURL}
	_, err := source.Start()
	server.Close()

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if len(spy.messagesLog) != 1 {
		t.Errorf("expected one message in WS got %d", len(spy.messagesLog))
	}

}
