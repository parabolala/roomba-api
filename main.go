package main

import (
	"github.com/ant0ine/go-json-rest"
	"net/http"
	"roomba-api"
)

func hello(w *rest.ResponseWriter, req *rest.Request) {
	w.WriteJson(&roomba_api.Status{Status: "ok"})
}

func main() {
	handler := roomba_api.MakeHttpHandler()
	http.ListenAndServe(":8080", &handler)
}
