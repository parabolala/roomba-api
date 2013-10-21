package roomba_api

import (
	"github.com/ant0ine/go-json-rest"
	"log"
	"net/http"
	"strconv"
)

type PortGetResponse struct {
	Status
	Ports []Port `json:"ports"`
}

type PortPostResponse struct {
	Status
	Name         string `json:"name"`
	ConnectionId uint64 `json:"connection_id"`
}

func (server *RoombaServer) GetPorts(w *rest.ResponseWriter, req *rest.Request) {
	ports := make([]Port, 0, 3)
	all_ports, err := listAllPorts()
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ports = append(ports, Port{DUMMY_PORT_NAME, PORT_STATE_AVAILABLE})
	for _, port_filename := range all_ports {
		log.Println("Found port " + port_filename)

		_, in_use := server.PortsInUse[port_filename]
		var state string
		if in_use {
			state = PORT_STATE_IN_USE
		} else {
			state = PORT_STATE_AVAILABLE
		}

		ports = append(ports, Port{port_filename, state})
	}
	status := PortGetResponse{Status: Status{"ok"},
		Ports: ports}
	w.WriteJson(&status)
}

func (server *RoombaServer) PostPorts(w *rest.ResponseWriter, r *rest.Request) {
	requested_port_name := r.PathParam("name")
	var found bool

	all_ports, err := listAllPorts()
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if requested_port_name != DUMMY_PORT_NAME {
		for _, port_filename := range all_ports {
			if port_filename == requested_port_name {
				found = true
				break
			}
		}

		if !found {
			Error(w, "Not found", http.StatusNotFound)
			return
		}
	}

	// race condition here
	_, ok := server.PortsInUse[requested_port_name]

	if ok {
		Error(w, "port is already in use", http.StatusGone)
		return
	}

	conn_id, err := server.GetConnection(requested_port_name)

	if err != nil {
		Error(w, "failed getting a connection: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	resp := PortPostResponse{Status{"ok"}, requested_port_name, conn_id}
	w.WriteJson(resp)
}

func (server *RoombaServer) DeleteConnection(w *rest.ResponseWriter, r *rest.Request) {
	conn_id_str := r.PathParam("conn_id")

	conn_id, err := strconv.ParseUint(conn_id_str, 10, 32)

	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// race condition here
	_, ok := server.Connections[conn_id]
	if !ok {
		Error(w, "connection id %d not found", http.StatusNotFound)
		return
	}

	err = server.CloseConnection(conn_id)

	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(Status{"ok"})
}
