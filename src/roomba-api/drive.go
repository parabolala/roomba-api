package roomba_api

import (
	"github.com/ant0ine/go-json-rest"
	"net/http"
    "fmt"
    "log"
	"strconv"
)

type DrivePutRequest struct {
    Velocity int16 `json:"velocity"`
    Radius int16 `json:"radius"`
}

type DrivePutResponse Status

func (server *RoombaServer) PutDrive(w *rest.ResponseWriter, req *rest.Request) {
    log.Println("Got PUT request")
	conn_id_str := req.PathParam("conn_id")

	conn_id, err := strconv.ParseUint(conn_id_str, 10, 32)

	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// race condition here
	_, ok := server.Connections[conn_id]
	if !ok {
		Error(w, "not found", http.StatusNotFound)
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
        Error(w, fmt.Sprintf("radius should either be in range [-2000; 2000] or " +
                             "be one of 32767, 32768, not %d", drive_req.Radius),
              http.StatusBadRequest)
        return
    }

    server.Connections[conn_id].Roomba.Drive(drive_req.Velocity,
                                             drive_req.Radius)

    w.WriteJson(&Status{"ok"})
    return
}

