package roomba_api

import (
	"fmt"
	"github.com/ant0ine/go-json-rest"
	"net/http"
)

type DrivePutRequest struct {
	Velocity int16 `json:"velocity"`
	Radius   int16 `json:"radius"`
}

type DrivePutResponse Status

func (server *RoombaServer) PutDrive(w *rest.ResponseWriter, req *rest.Request) {
	conn, err := server.getConnOrWriteError(w, req)
	if err != nil {
		return
	}

	drive_req := DrivePutRequest{}
	err = req.DecodeJsonPayload(&drive_req)
	if err != nil {
		Error(w, "can't decode json request", http.StatusBadRequest)
		return
	}

	if !(-500 <= drive_req.Velocity &&
		drive_req.Velocity <= 500) {
		Error(w, fmt.Sprintf("velocity should be in range [-500; 500], not %d",
			drive_req.Velocity),
			http.StatusBadRequest)
		return
	}

	if !(-2000 <= drive_req.Radius &&
		drive_req.Radius <= 2000 ||
		drive_req.Radius == ^0x7fff ||
		drive_req.Radius == 0x7fff) {
		Error(w, fmt.Sprintf("radius should either be in range [-2000; 2000] or "+
			"be one of 32767, 32768, not %d", drive_req.Radius),
			http.StatusBadRequest)
		return
	}

	conn.Roomba.Drive(drive_req.Velocity, drive_req.Radius)

	w.WriteJson(&Status{"ok"})
	return
}

type DirectDrivePutRequest struct {
	Left  int16 `json:"left"`
	Right int16 `json:"right"`
}
type DirectDriveResponse Status

func (server *RoombaServer) PutDirectDrive(w *rest.ResponseWriter, req *rest.Request) {
	conn, err := server.getConnOrWriteError(w, req)
	if err != nil {
		return
	}

	drive_req := DirectDrivePutRequest{}
	err = req.DecodeJsonPayload(&drive_req)
	if err != nil {
		Error(w, "can't decode json request", http.StatusBadRequest)
		return
	}

	if !(-500 <= drive_req.Left &&
		drive_req.Left <= 500) {
		Error(w, fmt.Sprintf("left velocity should be in range [-500; 500], not %d",
			drive_req.Left),
			http.StatusBadRequest)
		return
	}
	if !(-500 <= drive_req.Right &&
		drive_req.Right <= 500) {
		Error(w, fmt.Sprintf("left velocity should be in range [-500; 500], not %d",
			drive_req.Right),
			http.StatusBadRequest)
		return
	}

	conn.Roomba.DirectDrive(drive_req.Left, drive_req.Right)

	w.WriteJson(&Status{"ok"})
	return
}
