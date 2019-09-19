package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/api"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/helpers"
	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/ssurb"
)

func getNodeIDs() []int {
	return []int{0, 1, 2, 3, 4, 5}
}

func getID() int {
	idStr, exists := os.LookupEnv("ID")
	if !exists {
		log.Fatal("Environment variable ID missing, aborting")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal("Badly formatted ID env var")
	}
	return id
}

func main() {
	id := getID()

	// setup logging
	log.SetPrefix(fmt.Sprintf("[Node %d]: ", id))
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.Printf("Instance %d starting\n", id)

	// parse hosts and build P which is a essentially a slice of all node ids
	hosts, _ := helpers.ParseHostsFile()
	P := []int{}
	for _, p := range hosts {
		P = append(P, p.ID)
	}

	resolver := ssurb.Resolver{}

	// init module
	urbModule := ssurb.UrbModule{ID: id, P: P, Resolver: &resolver}
	urbModule.Init()
	hbfdModule := ssurb.HbfdModule{ID: id, P: P, Resolver: &resolver}
	hbfdModule.Init()
	thetafdModule := ssurb.ThetafdModule{ID: id, P: P, Resolver: &resolver}
	thetafdModule.Init()

	// init resolver and attach modules
	resolver.Modules = make(map[ssurb.ModuleType]interface{})
	resolver.Modules[ssurb.URB] = urbModule
	resolver.Modules[ssurb.HBFD] = hbfdModule
	resolver.Modules[ssurb.THETAFD] = thetafdModule

	// init waitgroup to keep track of all goroutines
	var wg sync.WaitGroup

	// setup communication
	server := ssurb.Server{IP: []byte{127, 0, 0, 1}, Port: 4000 + id, Resolver: &resolver}
	wg.Add(1)
	go func(s *ssurb.Server) {
		defer wg.Done()
		err := s.Start()
		if err != nil {
			log.Fatal(err)
		}
		s.Listen()
	}(&server)

	// launch hbfd module
	wg.Add(1)
	go func(module ssurb.HbfdModule) {
		defer wg.Done()
		module.DoForever()
	}(hbfdModule)

	// launch thetafd module
	wg.Add(1)
	go func(module ssurb.ThetafdModule) {
		defer wg.Done()
		module.DoForever()
	}(thetafdModule)

	// launch urb module
	wg.Add(1)
	go func(module ssurb.UrbModule) {
		defer wg.Done()
		module.DoForever()
	}(urbModule)

	// launch API
	api.SetUp(id, &resolver)

	// wait forever and allow modules and communication to run concurrently
	wg.Wait()
}
