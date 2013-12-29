package roomba_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"roomba"
	"strconv"
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

func TestSensorList(t *testing.T) {
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
	url_ := fmt.Sprintf("/connection/%d/sensor/list?", conn_id)

	requested_sensors := []byte{roomba.SENSOR_DISTANCE,
		roomba.SENSOR_WALL}
	expected_values := [][]byte{{10, 20},
		{35}}

	qs := url.Values{}
	for i, packet_data := range expected_values {
		qs.Add("packet_id", strconv.Itoa(int(requested_sensors[i])))
		MakeDummyRoomba().S.(*roomba.CloseableRWBuffer).WriteReadBuffer(packet_data)
	}

	url_ += qs.Encode()

	code, body = GetResponse(handler, "GET", url_)
	resp := GetSensorsResponse{}
	json.Unmarshal(body, &resp)

	if resp.Status.Status != "ok" {
		t.Errorf("status != ok")
	}

	if len(resp.Values) != len(expected_values) {
		t.Errorf("returned value len(%v) != expected output len(%v)", resp.Values, expected_values)
	}
	for i, packet_data := range resp.Values {
		if len(packet_data) != len(expected_values[i]) {
			t.Errorf("returned len for packet_id=%d data: %d != expected len %d", requested_sensors[i], len(resp.Values[i]), len(expected_values[i]))
		}
		for j, packet_byte := range resp.Values[i] {
			if packet_byte != expected_values[i][j] {
				t.Errorf("returned data for packet=%d %v!=expected %v in byte %d", requested_sensors[i], resp.Values[i], expected_values[i], j)
			}
		}
	}
}
