package roomba_api

import (
	"fmt"
)

type DriveRequest struct {
	ConnectionId uint64 `json:"connection_id"`
	Velocity     int16  `json:"velocity"`
	Radius       int16  `json:"radius"`
}
type DriveResponse struct{}

type DirectDriveRequest struct {
	ConnectionId uint64 `json:"connection_id"`
	Left         int16  `json:"left"`
	Right        int16  `json:"right"`
}
type DirectDriveResponse struct{}

func (server RoombaServer) Drive(req DriveRequest, resp *DriveResponse) error {
	conn, ok := server.Connections[req.ConnectionId]
	if !ok {
		return fmt.Errorf("connection not found: %d", req.ConnectionId)
	}

	if !(-500 <= req.Velocity && req.Velocity <= 500) {
		return fmt.Errorf("velocity should be in range [-500; 500], not %d", req.Velocity)
	}

	if !(-2000 <= req.Radius &&
		req.Radius <= 2000 ||
		req.Radius == ^0x7fff ||
		req.Radius == 0x7fff) {
		return fmt.Errorf(
			"radius should either be in range [-2000; 2000] or one of 32767, 32768, not %d",
			req.Radius,
		)
	}

	return conn.Roomba.Drive(req.Velocity, req.Radius)
}

func (server RoombaServer) DirectDrive(req DirectDriveRequest, resp *DirectDriveResponse) error {
	conn, ok := server.Connections[req.ConnectionId]
	if !ok {
		return fmt.Errorf("connection not found: %d", req.ConnectionId)
	}

	if !(-500 <= req.Left && req.Left <= 500) {
		return fmt.Errorf(
			"left velocity should be in range [-500; 500], not %d",
			req.Left,
		)
	}
	if !(-500 <= req.Right && req.Right <= 500) {
		return fmt.Errorf(
			"left velocity should be in range [-500; 500], not %d",
			req.Right,
		)
	}

	return conn.Roomba.DirectDrive(req.Left, req.Right)
}
