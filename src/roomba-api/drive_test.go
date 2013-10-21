package roomba_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"roomba"
	"testing"
	//	"github.com/ant0ine/go-json-rest"
	//    "net/http"
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

	data, _ := json.Marshal(DrivePutRequest{-200, 500})
	b := bytes.NewBuffer(data)
	code, body = GetResponseData(handler, "PUT",
		fmt.Sprintf("/connection/%d/control/drive", conn_id),
		b)

	if code != 200 {
		t.Errorf("response code is not 200")
	}

	resp := DrivePutResponse{}
	json.Unmarshal(body, &resp)

	if resp.Status != "ok" {
		t.Errorf("status != ok")
	}
	expected := []byte{128, 137, 255, 56, 1, 244}
	roomba.VerifyWritten( //server.Connections[conn_id].Roomba, expected, t)
		MakeDummyRoomba(), expected, t)
}
