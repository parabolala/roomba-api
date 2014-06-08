package roomba_api

import (
	"fmt"
	"log"
)

type GetPortsRequest struct{}
type GetPortsResponse struct {
	Ports []Port `json:"ports"`
}

type AcquireConnectionRequest struct {
	Name string `json:"name"`
}
type AcquireConnectionResponse struct {
	Name         string `json:"name"`
	ConnectionId uint64 `json:"connection_id"`
}

type ReleaseConnectionRequest struct {
	ConnectionId uint64 `json:"connection_id"`
}
type ReleaseConnectionResponse struct{}

func (server RoombaServer) GetPorts(req *GetPortsRequest, resp *GetPortsResponse) error {
	all_ports, err := listAllPorts()
	if err != nil {
		return err
	}
	resp.Ports = append(resp.Ports, Port{DUMMY_PORT_NAME, PORT_STATE_AVAILABLE})
	for _, port_filename := range all_ports {
		log.Println("Found port " + port_filename)

		_, in_use := server.PortsInUse[port_filename]
		var state string
		if in_use {
			state = PORT_STATE_IN_USE
		} else {
			state = PORT_STATE_AVAILABLE
		}

		resp.Ports = append(resp.Ports, Port{port_filename, state})
	}
	return nil
}

func (server RoombaServer) AcquireConnection(req *AcquireConnectionRequest, resp *AcquireConnectionResponse) error {
	requested_port_name := req.Name
	var found bool

	if requested_port_name != DUMMY_PORT_NAME {
		all_ports, err := listAllPorts()
		if err != nil {
			return err
		}

		for _, port_filename := range all_ports {
			if port_filename == requested_port_name {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("port %s not found", requested_port_name)
		}
	}

	// TODO: race condition here
	_, ok := server.PortsInUse[requested_port_name]

	if ok {
		return fmt.Errorf("port is already in use: %s", requested_port_name)
	}

	conn_id, err := server.GetConnection(requested_port_name)

	if err != nil {
		return fmt.Errorf("failed getting a connection: %s", err.Error())
	}

	resp.Name = requested_port_name
	resp.ConnectionId = conn_id
	return nil
}

func (server RoombaServer) ReleaseConnection(req *ReleaseConnectionRequest, resp *ReleaseConnectionResponse) error {
	conn, ok := server.Connections[req.ConnectionId]
	if !ok {
		return fmt.Errorf("connection not found: %d", conn)
	}

	err := server.CloseConnection(conn.Id)

	if err != nil {
		return err
	}
	return nil
}
