package swcsource

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type messageResponse struct {
	message  string
	response string
}

type wsSpy struct {
	messagesLog []string
}

func assertValue(t *testing.T, expect string, got string) {
	t.Helper()
	if expect != got {
		t.Errorf("Expected value of %s got %s", expect, got)
	}
}

func generateHTTPHandler(t *testing.T, expected []messageResponse, spy *wsSpy) func(http.ResponseWriter, *http.Request) {
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
