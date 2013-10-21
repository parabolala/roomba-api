package roomba_api

import (
    "testing"
    "io/ioutil"
    "io"
    "fmt"
    "net/http"
    "net/http/httptest"
    "log"
	"github.com/ant0ine/go-json-rest"
    "encoding/json"
    )

func GetResponseData(h rest.ResourceHandler, method, path string, body io.Reader) (int, []byte) {
    req, err := http.NewRequest(method,path, body)
    if err != nil {
        log.Fatal(err)
    }

    w := httptest.NewRecorder()
    h.ServeHTTP(w, req)

    bytes, err := ioutil.ReadAll(w.Body)
    if err != nil {
        log.Fatal(err)
    }

    return w.Code, bytes
}

func GetResponse(h rest.ResourceHandler, method, path string) (int, []byte) {
    return GetResponseData(h, method, path, nil)
}

func TestGetPorts(t *testing.T) {
    server := MakeServer()

    handler := MakeHttpHandlerForServer(server)
    code, body := GetResponse(handler, "GET", "/ports")
    
    if code != 200 {
        t.Errorf("response code is not 200")
    }

    resp := PortGetResponse{}
    json.Unmarshal(body, &resp)

    if resp.Status.Status != "ok" {
        t.Errorf("status != ok")
    }

    if len(resp.Ports) < 1 {
        t.Errorf("no ports found (not even DummyPort")
    }

    dummyPresent := false
    for _, Port := range(resp.Ports) {
        if Port.Name == DUMMY_PORT_NAME {
            dummyPresent = true
            break
        }
    }

    if !dummyPresent {
        t.Fatalf("DummyPort is not present in ports")
    }
}

func TestAcquireReleaseDummy(t *testing.T) {
    server := MakeServer()

    handler := MakeHttpHandlerForServer(server)
    code, body := GetResponse(handler, "POST", "/ports/" + DUMMY_PORT_NAME)

    if code != 200 {
        t.Errorf("response code is not 200")
    }

    resp := PortPostResponse{}
    json.Unmarshal(body, &resp)

    if resp.Status.Status != "ok" {
        t.Errorf("status != ok")
    }

    if resp.Name != DUMMY_PORT_NAME {
        t.Errorf("returned name (%s) != requested (%s)",
                 resp.Name, DUMMY_PORT_NAME)
    }


    code, body = GetResponse(handler, "DELETE", 
                              fmt.Sprintf("/connection/%d", resp.ConnectionId))

    if code != 200 {
        t.Errorf("response code is not 200")
    }

    del_resp := Status{}
    json.Unmarshal(body, &del_resp)

    if del_resp.Status != "ok" {
        t.Errorf("status != ok")
    }
}

func TestAcquireWrongPort(t *testing.T) {
    server := MakeServer()

    handler := MakeHttpHandlerForServer(server)
    code, body := GetResponse(handler, "POST", "/ports/wrong")

    if code != 404 {
        t.Errorf("response code is not 404")
    }

    resp := ErrorStatus{}
    json.Unmarshal(body, &resp)

    if resp.Status.Status == "ok" {
        t.Errorf("status == ok")
    }
}

func TestDeleteWrongConn(t *testing.T) {
    server := MakeServer()

    handler := MakeHttpHandlerForServer(server)
    code, body := GetResponse(handler, "DELETE", "/connection/42")

    if code != 404 {
        t.Errorf("response code is not 404")
    }

    resp := ErrorStatus{}
    json.Unmarshal(body, &resp)

    if resp.Status.Status == "ok" {
        t.Errorf("status == ok")
    }
}
