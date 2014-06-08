package roomba_api

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/rpc"
	"net/rpc/jsonrpc"
	"testing"

	"code.google.com/p/go.net/websocket"
)

var testServerAddr string
var testHttpServer *httptest.Server
var testRoombaServer RoombaServer

func StartTestServer() {
	testRoombaServer = MakeServer()

	rpc_server := rpc.NewServer()
	rpc_server.Register(testRoombaServer)
	var WebSocketHandler = func(ws *websocket.Conn) {
		server_codec := jsonrpc.NewServerCodec(ws)
		rpc_server.ServeCodec(server_codec)
	}
	serveMux := http.NewServeMux()
	serveMux.Handle("/ws", websocket.Handler(WebSocketHandler))

	testHttpServer = httptest.NewServer(serveMux)
	testServerAddr = testHttpServer.Listener.Addr().String()
	log.Print("Test WebSocket server listening on ", testServerAddr)
}

func StopTestServer() {
	testHttpServer.CloseClientConnections()
	testHttpServer.Close()
	testHttpServer = nil
}

func newTestClientConfig(path string) *websocket.Config {
	config, _ := websocket.NewConfig(fmt.Sprintf("ws://%s%s", testServerAddr, path), "http://localhost")
	return config
}

func NewTestClient(t *testing.T) *rpc.Client {
	// websocket.Dial()
	tcp_client, err := net.Dial("tcp", testServerAddr)
	if err != nil {
		t.Fatal("dialing", err)
	}
	ws_client, err := websocket.NewClient(
		newTestClientConfig("/ws"),
		tcp_client)

	if err != nil {
		t.Errorf("WebSocket handshake error: %v", err)
		return nil
	}

	return jsonrpc.NewClient(ws_client)
}
