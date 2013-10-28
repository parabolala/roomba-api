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

type GetSensorsResponse struct {
	Status
	Values [][]byte `json:"values"`
}

func (server *RoombaServer) GetSensors(w *rest.ResponseWriter, req *rest.Request) {
	conn, err := server.getConnOrWriteError(w, req)
	if err != nil {
		return
	}

	req.ParseForm()
	packet_id_strs := req.Form["packet_id"]

	packet_ids := make([]byte, 0, len(packet_id_strs))
	for _, packet_id_str := range packet_id_strs {
		packet_id, err := strconv.ParseUint(packet_id_str, 10, 8)
		if err != nil {
			Error(w, "bad packet id requested: "+packet_id_str,
				http.StatusNotFound)
			return
		}
		packet_ids = append(packet_ids, byte(packet_id))
	}

	sensor_data, err := conn.Roomba.QueryList(packet_ids)

	if err != nil {
		Error(w, "error reading sensors data: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	response := GetSensorsResponse{Status: Status{"ok"},
		Values: sensor_data}
	w.WriteJson(&response)
}
