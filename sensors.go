package roomba_api

import (
	"fmt"
)

type SensorRequest struct {
	ConnectionId uint64 `json:"connection_id"`
	PacketId     byte   `json:"packet_id"`
}
type SensorResponse struct {
	Value []byte `json:"value"`
}

type SensorListRequest struct {
	ConnectionId uint64 `json:"connection_id"`
	PacketIds    []byte `json:"packet_ids"`
}
type SensorListResponse struct {
	Values [][]byte `json:"values"`
}

func (server RoombaServer) Sensor(req SensorRequest, resp *SensorResponse) error {
	conn, ok := server.Connections[req.ConnectionId]
	if !ok {
		return fmt.Errorf("connection not found: %d", req.ConnectionId)
	}

	sensor_data, err := conn.Roomba.Sensors(req.PacketId)

	if err != nil {
		return fmt.Errorf("error reading sensor packet: " + err.Error())
	}

	resp.Value = sensor_data
	return nil
}

func (server RoombaServer) SensorList(req *SensorListRequest, resp *SensorListResponse) error {
	conn, ok := server.Connections[req.ConnectionId]
	if !ok {
		return fmt.Errorf("connection not found: %d", req.ConnectionId)
	}

	sensor_data, err := conn.Roomba.QueryList(req.PacketIds)

	if err != nil {
		return fmt.Errorf("error reading sensors data: " + err.Error())
	}

	resp.Values = sensor_data
	return nil
}
