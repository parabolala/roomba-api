package roomba_api

import (
	"errors"
	"log"

	"github.com/xa4a/go-roomba"
	rt "github.com/xa4a/go-roomba/testing"
)

type Port struct {
	Name string `json:"name"`
}

type Connection struct {
	Port       Port
	Roomba     *roomba.Roomba
	NumClients uint8
}

type RoombaServer struct {
	Connections     map[string]Connection
	connectionsChan chan ConnectionRequest
}

func (server *RoombaServer) acquireConnection(port_name string) (conn Connection, err error) {
	conn, ok := server.Connections[port_name]

	if ok {
		conn.NumClients += 1
		return
	}

	conn = Connection{NumClients: 1}
	if port_name == DUMMY_PORT_NAME {
		conn.Roomba = rt.MakeTestRoomba()
	} else {
		conn.Roomba, err = roomba.MakeRoomba(port_name)
		if err != nil {
			return
		}
	}
	conn.Roomba.Start()
	conn.Roomba.Safe()
	conn.Port = Port{Name: port_name}

	server.Connections[port_name] = conn
	return
}

func (server *RoombaServer) releaseConnection(port_name string) error {
	conn, ok := server.Connections[port_name]

	if !ok {
		return errors.New("unknown connection: " + port_name)
	}

	// close connection

	conn.NumClients -= 1

	if conn.NumClients == 0 {
		delete(server.Connections, port_name)
	}
	return nil
}

type operation int

const (
	OPEN  operation = iota
	CLOSE operation = iota
)

type ConnectionRequest struct {
	Port      string
	Operation operation

	C chan<- ConnectionResponse
}

type ConnectionResponse struct {
	Error error
}

func (server *RoombaServer) manageConnections() {
	var req ConnectionRequest
	var resp ConnectionResponse
	for {
		resp = ConnectionResponse{}
		select {
		case req = <-server.connectionsChan:
			if req.Operation == OPEN {
				log.Printf("Acquiring connection to " + req.Port)
				_, err := server.acquireConnection(req.Port)
				if err != nil {
					log.Println("Acquiring connection to " + req.Port +
						"failed: " + err.Error())

					resp.Error = err
				} else {
					log.Printf("Acquiring connection to " + req.Port +
						" success")
				}
			} else if req.Operation == CLOSE {
				log.Printf("Releasing connection %d", req.Port)
				resp.Error = server.releaseConnection(req.Port)
			} else {
				resp.Error = errors.New(
					"connectionrequest with unknown operation")
			}
			req.C <- resp
		}
	}
}

func (server *RoombaServer) GetConnection(port_name string) error {
	resp_chan := make(chan ConnectionResponse)
	log.Printf("Requesting connection to %s", port_name)
	req := ConnectionRequest{Port: port_name, Operation: OPEN, C: resp_chan}
	server.connectionsChan <- req
	resp := <-resp_chan
	return resp.Error
}

func (server *RoombaServer) CloseConnection(port_name string) error {
	resp_chan := make(chan ConnectionResponse)
	log.Printf("Closing connection %d", port_name)
	req := ConnectionRequest{Port: port_name, Operation: CLOSE, C: resp_chan}
	server.connectionsChan <- req
	resp := <-resp_chan
	return resp.Error
}

func MakeServer() (s RoombaServer) {
	s.Connections = make(map[string]Connection)
	s.connectionsChan = make(chan ConnectionRequest)
	go s.manageConnections()
	return
}
