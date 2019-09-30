package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/helpers"
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

type launchClientPayload struct {
	ReqCount int `json:"reqCount"`
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

func launchClient(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var payload launchClientPayload
	err := decoder.Decode(&payload)
	if err != nil {
		panic(err)
	}

	go func(reqCount int) {
		mod := resolver.GetUrbModule()
		log.Println("Launching client")

		// TODO look into optimising this one
		if reqCount != -1 {
			for i := 0; i < reqCount; i++ {
				mod.BlockUntilAvailableSpace()
				mod.UrbBroadcast(&ssurb.UrbMessage{Text: fmt.Sprintf("Message %d_%d", mod.ID, i)})
			}
		}
	}(payload.ReqCount)

	w.WriteHeader(200)
}

// SetUp launches the API and registers all handlers
func SetUp(id int, r *ssurb.Resolver) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/client/launch", launchClient).Methods("POST")
	router.HandleFunc("/broadcast", broadcast).Methods("POST")

	port := 4000 + id
	resolver = r
	ipString := helpers.GetIP()
	log.Printf("Launching API on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", ipString, port), router))
}
