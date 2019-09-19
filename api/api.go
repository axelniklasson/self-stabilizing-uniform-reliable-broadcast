package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/ssurb"

	"github.com/gorilla/mux"
)

var resolver *ssurb.Resolver

type response struct {
	Endpoint   string
	StatusCode int
	Data       interface{}
	Error      error
}

func index(w http.ResponseWriter, r *http.Request) {
	res := response{Endpoint: "/", StatusCode: 200, Data: map[string]interface{}{"foo": "bar", "trusted": resolver.Trusted()}}
	json.NewEncoder(w).Encode(res)
}

// SetUp launches the API and registers all handlers
func SetUp(id int, r *ssurb.Resolver) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	port := 5000 + id
	resolver = r

	log.Printf("Launching API on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
