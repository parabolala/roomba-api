package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest"
	"net/http"
	"roomba-api"
)

func hello(w *rest.ResponseWriter, req *rest.Request) {
	w.WriteJson(&roomba_api.Status{Status: "ok"})
}

func main() {
	fmt.Println("Serving..")
	handler := roomba_api.MakeHttpHandler()
	http.ListenAndServe(":8080", &handler)
}
