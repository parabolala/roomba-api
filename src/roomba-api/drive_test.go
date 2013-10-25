package roomba_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"roomba"
	"testing"
	//	"github.com/ant0ine/go-json-rest"
	"net/http"
)

func TestDriveOk(t *testing.T) {
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
	url := fmt.Sprintf("/connection/%d/control/drive", conn_id)

	data, _ := json.Marshal(DrivePutRequest{-200, 500})
	b := bytes.NewBuffer(data)
	code, body = GetResponseData(handler, "PUT", url, b)

	if code != 200 {
		t.Errorf("response code is not 200")
	}

	resp := DrivePutResponse{}
	json.Unmarshal(body, &resp)

	if resp.Status != "ok" {
		t.Errorf("status != ok")
	}
	expected := []byte{128, 131, 137, 255, 56, 1, 244}
	roomba.VerifyWritten(server.Connections[conn_id].Roomba, expected, t)

	// Special cases radius.
	for _, radius := range []int16{-2000, 2000, 0, 10, -10, 32767, -32768} {
		data, _ = json.Marshal(DrivePutRequest{-200, radius})
		b = bytes.NewBuffer(data)
		code, body = GetResponseData(handler, "PUT", url, b)

		if code != 200 {
			t.Errorf("response code is not 200 for radius %d", radius)
		}
	}
}

func TestDriveWrongConnId(t *testing.T) {
	server := MakeServer()

	handler := MakeHttpHandlerForServer(server)

	data, _ := json.Marshal(DrivePutRequest{-200, 500})
	b := bytes.NewBuffer(data)
	code, _ := GetResponseData(handler, "PUT", "/connection/42/control/drive", b)

	if code != http.StatusNotFound {
		t.Errorf("response code is not 404")
	}
}

func TestDriveWrongVelocityRadius(t *testing.T) {
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
	url := fmt.Sprintf("/connection/%d/control/drive", conn_id)

	for _, velocity := range []int16{-501, -1000, 501, 1001} {
		data, _ := json.Marshal(DrivePutRequest{velocity, 500})
		b := bytes.NewBuffer(data)
		code, _ = GetResponseData(handler, "PUT", url, b)

		if code != http.StatusBadRequest {
			t.Errorf("response code is not 400 for velocity %d", velocity)
		}
	}

	for _, radius := range []int16{-2001, -10000, 2001, 10000} {
		data, _ := json.Marshal(DrivePutRequest{315, radius})
		b := bytes.NewBuffer(data)
		code, _ = GetResponseData(handler, "PUT", url, b)

		if code != http.StatusBadRequest {
			t.Errorf("response code is not 400 for radius %d", radius)
		}
	}
}

func TestDirectDriveOk(t *testing.T) {
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
	url := fmt.Sprintf("/connection/%d/control/direct_drive", conn_id)

	data, _ := json.Marshal(DirectDrivePutRequest{127, 256})
	b := bytes.NewBuffer(data)
	code, body = GetResponseData(handler, "PUT", url, b)

	if code != 200 {
		t.Errorf("response code is not 200")
	}

	resp := DrivePutResponse{}
	json.Unmarshal(body, &resp)

	if resp.Status != "ok" {
		t.Errorf("status != ok")
	}
	expected := []byte{128, 131, 145, 0, 127, 1, 0}
	roomba.VerifyWritten(server.Connections[conn_id].Roomba, expected, t)
}

func TestDirectDriveWrongVelocity(t *testing.T) {
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
	url := fmt.Sprintf("/connection/%d/control/direct_drive", conn_id)

	for _, velocities := range [][2]int16{{-501, 500},
		{-1000, 501},
		{35, 1002}} {
		left := velocities[0]
		right := velocities[1]
		data, _ := json.Marshal(DirectDrivePutRequest{left, right})
		b := bytes.NewBuffer(data)
		code, _ = GetResponseData(handler, "PUT", url, b)

		if code != http.StatusBadRequest {
			t.Errorf("response code is not 400 for velocities %v", velocities)
		}
	}
}
