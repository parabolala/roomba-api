package roomba_api

import (
	"testing"

	rt "github.com/xa4a/go-roomba/testing"
)

func TestDriveOk(t *testing.T) {
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

	drive_req := DriveRequest{port_name, -200, 500}
	err = client.Call("RoombaServer.Drive", drive_req, &DriveResponse{})

	if err != nil {
		t.Fatalf("rpc call failed unexpectedly: %s", err)
	}

	expected := []byte{128, 131, 137, 255, 56, 1, 244}
	rt.VerifyWritten(testRoombaServer.Connections[port_name].Roomba, expected, t)

	// Special cases radius.
	for _, radius := range []int16{-2000, 2000, 0, 10, -10, 32767, -32768} {
		drive_req = DriveRequest{port_name, -200, radius}
		err = client.Call("RoombaServer.Drive", drive_req, &DriveResponse{})

		if err != nil {
			t.Fatalf("rpc call failed unexpectedly for radius %d: %s", radius, err)
		}
	}
}

func TestDriveWrongConnId(t *testing.T) {
	StartTestServer()
	defer StopTestServer()

	client := NewTestClient(t)
	defer client.Close()

	drive_req := DriveRequest{"foo", -200, 500}
	err := client.Call("RoombaServer.Drive", drive_req, &DriveResponse{})

	if err == nil {
		t.Fatalf("rpc call succeeded unexpectedly")
	}
}

func TestDriveWrongVelocityRadius(t *testing.T) {
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

	for _, velocity := range []int16{-501, -1000, 501, 1001} {
		drive_req := DriveRequest{port_name, velocity, 500}
		err = client.Call("RoombaServer.Drive", drive_req, &DriveResponse{})

		if err == nil {
			t.Errorf("rpc call succeeded unexpectedly")
		}
	}

	for _, radius := range []int16{-2001, -10000, 2001, 10000} {
		drive_req := DriveRequest{port_name, 315, radius}
		err = client.Call("RoombaServer.Drive", drive_req, &DriveResponse{})

		if err == nil {
			t.Errorf("rpc call succeeded unexpectedly")
		}
	}
}

func TestDirectDriveOk(t *testing.T) {
	StartTestServer()
	defer StopTestServer()
	defer rt.ClearTestRoomba()

	client := NewTestClient(t)
	defer client.Close()

	conn_req := AcquireConnectionRequest{Port: DUMMY_PORT_NAME}
	conn_resp := AcquireConnectionResponse{}

	err := client.Call("RoombaServer.AcquireConnection", conn_req, &conn_resp)

	if err != nil {
		t.Errorf("failed acquiring dummy connection: %s", err)
	}

	port_name := conn_resp.Port

	drive_req := DirectDriveRequest{port_name, 127, 256}
	err = client.Call("RoombaServer.DirectDrive", drive_req, &DirectDriveResponse{})

	if err != nil {
		t.Fatalf("rpc call failed unexpectedly: %s", err)
	}

	expected := []byte{128, 131, 145, 0, 127, 1, 0}
	rt.VerifyWritten(testRoombaServer.Connections[port_name].Roomba, expected, t)
}

func TestDirectDriveWrongConnId(t *testing.T) {
	StartTestServer()
	defer StopTestServer()

	client := NewTestClient(t)
	defer client.Close()

	drive_req := DirectDriveRequest{"foo", 127, 256}
	err := client.Call("RoombaServer.DirectDrive", drive_req, &DirectDriveResponse{})

	if err == nil {
		t.Fatalf("rpc call succeeded unexpectedly")
	}
}

func TestDirectDriveWrongVelocity(t *testing.T) {
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
	for _, velocities := range [][2]int16{{-501, 500},
		{-1000, 501},
		{35, 1002}} {
		left := velocities[0]
		right := velocities[1]

		drive_req := DirectDriveRequest{port_name, left, right}
		err = client.Call("RoombaServer.DirectDrive", drive_req, &DirectDriveResponse{})

		if err == nil {
			t.Errorf("rpc call succeeded unexpectedly")
		}
	}
}
