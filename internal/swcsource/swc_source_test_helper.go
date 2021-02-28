package swcsource

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type wsSpy struct {
	messagesLog    []string
	messageChannel chan string

	stop bool
}

func makeSpy() *wsSpy {
	spy := wsSpy{
		messageChannel: make(chan string, 100),
	}

	return &spy
}

func assertValue(t *testing.T, expect string, got string) {
	t.Helper()
	if expect != got {
		t.Errorf("Expected value of %s got %s", expect, got)
	}
}

func generateHTTPHandler(t *testing.T, spy *wsSpy) func(http.ResponseWriter, *http.Request) {
	t.Helper()

	expected := map[string][]byte{
		"LOGIN;000000": readFixture(t, "testdata/LOGIN;123456.xml"),
		"GET;0x46bd50": readFixture(t, "testdata/GET;0x46bd50.xml"),
		"REFRESH":      readFixture(t, "testdata/GET;0x46bd50-REFRESH.xml"),
	}

	return func(response http.ResponseWriter, request *http.Request) {
		t.Helper()
		connection, err := upgrader.Upgrade(response, request, nil)
		if err != nil {
			return
		}
		defer connection.Close()
		for !spy.stop {
			mt, command, err := connection.ReadMessage()
			if err != nil {
				break
			}
			stringCommand := string(command)
			response := expected[stringCommand]
			spy.messagesLog = append(spy.messagesLog, stringCommand)

			err = connection.WriteMessage(mt, response)
			if err != nil {
				break
			}

			spy.messageChannel <- stringCommand
		}
	}
}

func readFixture(t *testing.T, fileName string) []byte {
	t.Helper()

	data, readError := os.ReadFile(fileName)
	if readError != nil {
		t.Fatalf("failed to read file %v", readError)
	}

	return data
}

func toWs(url string) string {
	return "ws" + strings.TrimPrefix(url, "http")
}
