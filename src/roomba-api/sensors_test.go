package roomba_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"roomba"
	"testing"
)

func TestSensorOk(t *testing.T) {
	server := MakeServer()

	handler := MakeHttpHandlerForServer(server)
	code, body := GetResponse(handler, "POST", "/ports/"+DUMMY_PORT_NAME)

	if code != 200 {
		t.Errorf("failed acquiring dummy connection")
	}

	port_resp := PortPostResponse{}
	json.Unmarshal(body, &port_resp)
	conn_id := port_resp.ConnectionId
	defer ClearDummyRoomba()
	url := fmt.Sprintf("/connection/%d/sensor", conn_id)

	output := []byte{42}
	MakeDummyRoomba().S.(*roomba.CloseableRWBuffer).WriteReadBuffer(output)

	code, body = GetResponse(handler, "GET",
		fmt.Sprintf("%s/%d", url, roomba.SENSOR_CLIFF_RIGHT))
	resp := SensorResponse{}
	json.Unmarshal(body, &resp)

	if resp.Status.Status != "ok" {
		t.Errorf("status != ok")
	}

	if len(resp.Value) != len(output) || (resp.Value[0] != output[0]) {
		t.Errorf("returned value (%v)!= expected output (%v)", resp.Value, output)
	}
}

func TestBadSensorURL(t *testing.T) {
	server := MakeServer()

	handler := MakeHttpHandlerForServer(server)
	code, body := GetResponse(handler, "POST", "/ports/"+DUMMY_PORT_NAME)

	if code != 200 {
		t.Errorf("failed acquiring dummy connection")
	}

	port_resp := PortPostResponse{}
	json.Unmarshal(body, &port_resp)
	conn_id := port_resp.ConnectionId
	defer ClearDummyRoomba()
	url := fmt.Sprintf("/connection/%d/sensor", conn_id)

	code, body = GetResponse(handler, "GET",
		fmt.Sprintf("%s/wrong_Sensor_code", url))
	resp := SensorResponse{}
	json.Unmarshal(body, &resp)

	if resp.Status.Status == "ok" {
		t.Errorf("status == ok")
	}

	if code != http.StatusNotFound {
		t.Errorf("code != 404")
	}
}
