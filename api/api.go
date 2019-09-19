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

type broadcastPayload struct {
	Text string `json:"text"`
}

func index(w http.ResponseWriter, r *http.Request) {
	res := response{Endpoint: "/", StatusCode: 200, Data: map[string]interface{}{"foo": "bar", "trusted": resolver.Trusted()}}
	json.NewEncoder(w).Encode(res)
}

func broadcast(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var payload broadcastPayload
	err := decoder.Decode(&payload)
	if err != nil {
		panic(err)
	}

	msg := ssurb.UrbMessage{Text: payload.Text}
	go resolver.UrbBroadcast(&msg)
}

// SetUp launches the API and registers all handlers
func SetUp(id int, r *ssurb.Resolver) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/broadcast", broadcast).Methods("POST")
	port := 5000 + id
	resolver = r

	log.Printf("Launching API on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
