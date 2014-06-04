package roomba_api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/xa4a/go-roomba"
	rt "github.com/xa4a/go-roomba/testing"
)

const PORT_STATE_AVAILABLE string = "available"
const PORT_STATE_IN_USE string = "in use"

type Port struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

type Status struct {
	Status string `json:"status"`
}

type ErrorStatus struct {
	Status
	Reason string `json:"reason"`
}

func Error(w rest.ResponseWriter, error string, code int) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	resp := ErrorStatus{Status: Status{"error"},
		Reason: error}
	err := w.WriteJson(resp)
	if err != nil {
		panic(err)
	}
}

type Connection struct {
	Id     uint64
	Port   Port
	Roomba *roomba.Roomba
}

type RoombaServer struct {
	Connections     map[uint64]Connection
	connectionsChan chan ConnectionRequest

	PortsInUse map[string]bool
	nextConnId uint64
}

func (server *RoombaServer) acquireConnection(port_name string) (conn Connection, err error) {
	_, ok := server.PortsInUse[port_name]

	if ok {
		err = errors.New("port is already in use: " + port_name)
		return
	}

	if port_name != DUMMY_PORT_NAME {
		server.PortsInUse[port_name] = true
	}

	conn.Id = server.nextConnId
	server.nextConnId++
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
	conn.Port = Port{Name: port_name, State: PORT_STATE_IN_USE}

	server.Connections[conn.Id] = conn
	return
}

func (server *RoombaServer) releaseConnection(conn_id uint64) error {
	conn, ok := server.Connections[conn_id]

	if !ok {
		return errors.New("unknown connection: " + fmt.Sprintf("%d", conn_id))
	}

	// close connection

	delete(server.Connections, conn_id)

	if conn.Port.Name != DUMMY_PORT_NAME {
		delete(server.PortsInUse, conn.Port.Name)
	}
	return nil
}

type ConnectionRequest struct {
	PortName     string
	ConnectionId uint64
	C            chan<- ConnectionResponse
}

type ConnectionResponse struct {
	ConnectionId uint64
	Error        error
}

func (server *RoombaServer) manageConnections() {
	var req ConnectionRequest
	var resp ConnectionResponse
	for {
		resp = ConnectionResponse{}
		select {
		case req = <-server.connectionsChan:
			if req.PortName != "" {
				log.Printf("Acquiring connection to " + req.PortName)
				conn, err := server.acquireConnection(req.PortName)
				if err != nil {
					log.Println("Acquiring connection to " + req.PortName +
						"failed: " + err.Error())

					resp.Error = err
				} else {
					log.Printf("Acquiring connection to "+req.PortName+
						" success: %d", conn.Id)
					resp.ConnectionId = conn.Id
				}
			} else if req.ConnectionId != 0 {
				log.Printf("Releasing connection %d", req.ConnectionId)
				resp.Error = server.releaseConnection(req.ConnectionId)
			} else {
				resp.Error = errors.New(
					"connectionrequest with neither PortName nor ConnectionId")
			}
			req.C <- resp
		}
	}
}

func (server *RoombaServer) GetConnection(port_name string) (uint64, error) {
	resp_chan := make(chan ConnectionResponse)
	log.Printf("Requesting connection to %s", port_name)
	req := ConnectionRequest{PortName: port_name, C: resp_chan}
	server.connectionsChan <- req
	resp := <-resp_chan
	return resp.ConnectionId, resp.Error
}

func (server *RoombaServer) CloseConnection(conn_id uint64) error {
	resp_chan := make(chan ConnectionResponse)
	log.Printf("Closing connection %d", conn_id)
	req := ConnectionRequest{ConnectionId: conn_id, C: resp_chan}
	server.connectionsChan <- req
	resp := <-resp_chan
	return resp.Error
}

func MakeServer() (s RoombaServer) {
	s.Connections = make(map[uint64]Connection)
	s.connectionsChan = make(chan ConnectionRequest)
	s.PortsInUse = make(map[string]bool)
	s.nextConnId = 1
	go s.manageConnections()
	return
}

func MakeHttpHandlerForServer(server RoombaServer) rest.ResourceHandler {
	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
	}
	handler.SetRoutes(
		rest.RouteObjectMethod("GET", "/ports", &server, "GetPorts"),
		rest.RouteObjectMethod("POST", "/ports/*name", &server, "PostPorts"),
		rest.RouteObjectMethod("DELETE", "/connection/:conn_id", &server, "DeleteConnection"),
		rest.RouteObjectMethod("PUT", "/connection/:conn_id/control/drive",
			&server, "PutDrive"),
		rest.RouteObjectMethod("PUT", "/connection/:conn_id/control/direct_drive",
			&server, "PutDirectDrive"),
		rest.RouteObjectMethod("GET", "/connection/:conn_id/sensor/list",
			&server, "GetSensors"),
		rest.RouteObjectMethod("GET", "/connection/:conn_id/sensor/:packet_id",
			&server, "GetSensor"),
	)
	return handler
}

func MakeHttpHandler() rest.ResourceHandler {
	server := MakeServer()
	return MakeHttpHandlerForServer(server)
}

func (server *RoombaServer) getConnOrWriteError(
	w rest.ResponseWriter, req *rest.Request) (conn Connection, err error) {
	conn_id_str := req.PathParam("conn_id")

	conn_id, err := strconv.ParseUint(conn_id_str, 10, 32)

	if err != nil {
		Error(w, "malformed connection id: "+err.Error(), http.StatusBadRequest)
		return
	}

	// race condition here
	conn, ok := server.Connections[conn_id]
	if !ok {
		Error(w, "connection not found", http.StatusNotFound)
		err = errors.New("connection not found: " + conn_id_str)
		return
	}
	return
}
