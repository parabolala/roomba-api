package roomba_api

import (
	"github.com/ant0ine/go-json-rest"
	"net/http"
	"strconv"
)

type SensorRequest struct {
	PacketId byte `json:"packet_id"`
}

type SensorResponse struct {
	Status
	Value []byte `json:"value"`
}

func (server *RoombaServer) GetSensor(w *rest.ResponseWriter, req *rest.Request) {
	conn, err := server.getConnOrWriteError(w, req)
	if err != nil {
		return
	}

	packet_id_str := req.PathParam("packet_id")

	packet_id, err := strconv.ParseUint(packet_id_str, 10, 8)
	if err != nil {
		Error(w, "bad packet id requested: "+packet_id_str,
			http.StatusNotFound)
		return
	}

	sensor_data, err := conn.Roomba.Sensors(byte(packet_id))

	if err != nil {
		Error(w, "error reading sensor packet: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	response := SensorResponse{Status: Status{"ok"},
		Value: sensor_data}
	w.WriteJson(&response)
}
