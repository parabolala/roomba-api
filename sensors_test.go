package roomba_api

import (
	"testing"

	"github.com/xa4a/go-roomba/constants"
	"github.com/xa4a/go-roomba/sim"
	rt "github.com/xa4a/go-roomba/testing"
)

func TestSensorOk(t *testing.T) {
	StartTestServer()
	defer StopTestServer()
	defer rt.ClearTestRoomba()

	client := NewTestClient(t)
	defer client.Close()

	conn_req := AcquireConnectionRequest{Port: DUMMY_PORT_NAME}
	conn_resp := AcquireConnectionResponse{}

	err := client.Call("RoombaServer.AcquireConnection", conn_req, &conn_resp)

	if err != nil {
		t.Fatalf("failed acquiring dummy connection: %s", err)
	}

	port_name := conn_resp.Port

	output := sim.MockSensorValues[constants.SENSOR_CLIFF_RIGHT]

	s_req := SensorRequest{
		Port:     port_name,
		PacketId: constants.SENSOR_CLIFF_RIGHT,
	}
	resp := &SensorResponse{}
	err = client.Call("RoombaServer.Sensor", s_req, resp)

	if err != nil {
		t.Fatalf("rpc call failed: %s", err)
	}

	if len(resp.Value) != len(output) || (resp.Value[0] != output[0]) {
		t.Errorf("returned value (%v)!= expected output (%v)", resp.Value, output)
	}
}

func TestBadSensorPacketId(t *testing.T) {
	StartTestServer()
	defer StopTestServer()
	defer rt.ClearTestRoomba()

	client := NewTestClient(t)
	defer client.Close()

	conn_req := AcquireConnectionRequest{Port: DUMMY_PORT_NAME}
	conn_resp := AcquireConnectionResponse{}

	err := client.Call("RoombaServer.AcquireConnection", conn_req, &conn_resp)

	if err != nil {
		t.Fatalf("failed acquiring dummy connection: %s", err)
	}

	port_name := conn_resp.Port

	s_req := SensorRequest{Port: port_name, PacketId: 250}
	s_resp := &SensorResponse{}
	err = client.Call("RoombaServer.Sensor", s_req, s_resp)

	if err == nil {
		t.Errorf("rpc call unexpectedly succeeded")
	}
}

func TestBadSensorConnId(t *testing.T) {
	StartTestServer()
	defer StopTestServer()

	client := NewTestClient(t)
	defer client.Close()

	s_req := SensorRequest{Port: "foo", PacketId: 250}
	s_resp := &SensorResponse{}
	err := client.Call("RoombaServer.Sensor", s_req, s_resp)

	if err == nil {
		t.Errorf("rpc call unexpectedly succeeded")
	}
}

func TestSensorList(t *testing.T) {
	StartTestServer()
	defer StopTestServer()
	defer rt.ClearTestRoomba()

	client := NewTestClient(t)
	defer client.Close()

	conn_req := AcquireConnectionRequest{Port: DUMMY_PORT_NAME}
	conn_resp := AcquireConnectionResponse{}

	err := client.Call("RoombaServer.AcquireConnection", conn_req, &conn_resp)

	if err != nil {
		t.Fatalf("failed acquiring dummy connection: %s", err)
	}

	port_name := conn_resp.Port

	requested_sensors := []byte{
		constants.SENSOR_DISTANCE,
		constants.SENSOR_WALL,
	}
	expected_values := [][]byte{
		sim.MockSensorValues[requested_sensors[0]],
		sim.MockSensorValues[requested_sensors[1]],
	}

	s_req := SensorListRequest{
		Port:      port_name,
		PacketIds: requested_sensors,
	}
	resp := &SensorListResponse{}

	err = client.Call("RoombaServer.SensorList", s_req, resp)

	if err != nil {
		t.Fatalf("rpc call failed: %s", err)
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

func TestSensorListBadConnId(t *testing.T) {
	StartTestServer()
	defer StopTestServer()

	client := NewTestClient(t)
	defer client.Close()

	requested_sensors := []byte{
		constants.SENSOR_DISTANCE,
	}

	s_req := SensorListRequest{
		Port:      "foo",
		PacketIds: requested_sensors,
	}
	resp := &SensorListResponse{}

	err := client.Call("RoombaServer.SensorList", s_req, resp)

	if err == nil {
		t.Fatalf("rpc call succeeded unexpectedly")
	}
}
