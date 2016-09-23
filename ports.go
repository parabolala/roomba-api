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
	Port string `json:"port_name"`
}
type AcquireConnectionResponse struct {
	Port string `json:"port_name"`
}

type ReleaseConnectionRequest struct {
	Port string `json:"port_name"`
}
type ReleaseConnectionResponse struct{}

func (server RoombaServer) GetPorts(req *GetPortsRequest, resp *GetPortsResponse) error {
	all_ports, err := listAllPorts()
	if err != nil {
		return err
	}
	resp.Ports = append(resp.Ports, Port{DUMMY_PORT_NAME})
	for _, port_filename := range all_ports {
		log.Println("Found port " + port_filename)

		resp.Ports = append(resp.Ports, Port{port_filename})
	}
	return nil
}

func (server RoombaServer) AcquireConnection(req *AcquireConnectionRequest, resp *AcquireConnectionResponse) error {
	requested_port_name := req.Port
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

	err := server.GetConnection(requested_port_name)

	if err != nil {
		return fmt.Errorf("failed getting a connection: %s", err.Error())
	}

	resp.Port = requested_port_name
	return nil
}

func (server RoombaServer) ReleaseConnection(req *ReleaseConnectionRequest, resp *ReleaseConnectionResponse) error {
	conn, ok := server.Connections[req.Port]
	if !ok {
		return fmt.Errorf("connection not found: %d", conn)
	}

	err := server.CloseConnection(conn.Port.Name)

	if err != nil {
		return err
	}
	return nil
}
