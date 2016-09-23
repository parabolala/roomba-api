package roomba_api

import (
	"fmt"
)

type SensorRequest struct {
	Port     string `json:"port_name"`
	PacketId byte   `json:"packet_id"`
}
type SensorResponse struct {
	Value []byte `json:"value"`
}

type SensorListRequest struct {
	Port      string `json:"port_name"`
	PacketIds []byte `json:"packet_ids"`
}
type SensorListResponse struct {
	Values [][]byte `json:"values"`
}

func (server RoombaServer) Sensor(req SensorRequest, resp *SensorResponse) error {
	conn, ok := server.Connections[req.Port]
	if !ok {
		return fmt.Errorf("connection not found: %s", req.Port)
	}

	sensor_data, err := conn.Roomba.Sensors(req.PacketId)

	if err != nil {
		return fmt.Errorf("error reading sensor packet: " + err.Error())
	}

	resp.Value = sensor_data
	return nil
}

func (server RoombaServer) SensorList(req *SensorListRequest, resp *SensorListResponse) error {
	conn, ok := server.Connections[req.Port]
	if !ok {
		return fmt.Errorf("connection not found: %s", req.Port)
	}

	sensor_data, err := conn.Roomba.QueryList(req.PacketIds)

	if err != nil {
		return fmt.Errorf("error reading sensors data: " + err.Error())
	}

	resp.Values = sensor_data
	return nil
}
