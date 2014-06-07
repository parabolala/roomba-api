package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/xa4a/roomba-api"
)

func hello(w rest.ResponseWriter, req *rest.Request) {
	w.WriteJson(&roomba_api.Status{Status: "ok"})
}

func main() {
	fmt.Println("Serving..")
	handler := roomba_api.MakeHttpHandler()
	err := http.ListenAndServe(":8080", &handler)
	if err != nil {
		log.Fatalf("failed starting server: %s", err)
	}
}
