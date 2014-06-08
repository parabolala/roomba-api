package roomba_api

import (
	"testing"
)

func TestPortsGet(t *testing.T) {
	StartTestServer()
	defer StopTestServer()

	client := NewTestClient(t)
	defer client.Close()

	resp := GetPortsResponse{}
	err := client.Call("RoombaServer.GetPorts", GetPortsRequest{}, &resp)

	if err != nil {
		t.Fatalf("rpc call failed: %s", err)
	}

	if len(resp.Ports) < 1 {
		t.Errorf("no ports found (not even DummyPort")
	}

	dummyPresent := false
	for _, Port := range resp.Ports {
		if Port.Name == DUMMY_PORT_NAME {
			dummyPresent = true
			break
		}
	}

	if !dummyPresent {
		t.Errorf("DummyPort is not present in ports")
	}
}

func TestAcquireReleaseDummy(t *testing.T) {
	StartTestServer()
	defer StopTestServer()

	client := NewTestClient(t)
	defer client.Close()

	req := &AcquireConnectionRequest{
		Name: DUMMY_PORT_NAME,
	}
	resp := &AcquireConnectionResponse{}
	err := client.Call("RoombaServer.AcquireConnection", req, &resp)

	if err != nil {
		t.Fatalf("rpc call failed: %s", err)
	}

	if resp.Name != DUMMY_PORT_NAME {
		t.Errorf("returned name (%s) != requested (%s)",
			resp.Name, DUMMY_PORT_NAME)
	}

	err = client.Call(
		"RoombaServer.ReleaseConnection",
		&ReleaseConnectionRequest{ConnectionId: resp.ConnectionId},
		&ReleaseConnectionResponse{},
	)

	if err != nil {
		t.Fatalf("rpc call returned an error: %s", err)
	}
}

func TestAcquireWrongPort(t *testing.T) {
	StartTestServer()
	defer StopTestServer()

	client := NewTestClient(t)
	defer client.Close()

	req := &AcquireConnectionRequest{Name: "/ports/wrong"}
	resp := &AcquireConnectionResponse{}
	err := client.Call("RoombaServer.AcquireConnection", req, &resp)

	if err == nil {
		t.Errorf("acquiring wrong port name has not errored out")
	}
}

func TestDeleteWrongConn(t *testing.T) {
	StartTestServer()
	defer StopTestServer()

	client := NewTestClient(t)
	defer client.Close()

	req := &ReleaseConnectionRequest{
		ConnectionId: 42,
	}

	err := client.Call("RoombaServer.ReleaseConnection", req, &ReleaseConnectionResponse{})

	if err == nil {
		t.Errorf("releasing bad connection id has not errored out")
	}
}
