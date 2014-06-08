package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/xa4a/roomba-api"

	"code.google.com/p/go.net/websocket"
)

func main() {
	server := roomba_api.MakeServer()

	rpc.Register(server)

	// Gob over HTTP.
	rpc.HandleHTTP()

	// JSON-RPC over Websocket over HTTP.
	var WebSocketHandler = func(ws *websocket.Conn) {
		server_codec := jsonrpc.NewServerCodec(ws)
		rpc.DefaultServer.ServeCodec(server_codec)
	}
	http.Handle("/ws", websocket.Handler(WebSocketHandler))

	fmt.Println("Serving..")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("failed starting server: %s", err)
	}
}
