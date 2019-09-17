package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/comm"
	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/helpers"
	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/modules"
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

	resolver := modules.Resolver{}

	// init module
	urbModule := modules.UrbModule{ID: id, P: P, Resolver: &resolver}
	urbModule.Init()
	hbfdModule := modules.HbfdModule{ID: id, P: P, Resolver: &resolver}
	hbfdModule.Init()
	thetafdModule := modules.ThetafdModule{ID: id, P: P, Resolver: &resolver}
	thetafdModule.Init()

	// init resolver and attach modules
	resolver.Modules = make(map[modules.ModuleType]interface{})
	resolver.Modules[modules.URB] = urbModule
	resolver.Modules[modules.HBFD] = hbfdModule
	resolver.Modules[modules.THETAFD] = thetafdModule

	// init waitgroup to keep track of all goroutines
	var wg sync.WaitGroup

	// setup communication
	server := comm.Server{IP: []byte{127, 0, 0, 1}, Port: 4000 + id}
	wg.Add(1)
	go func(s *comm.Server) {
		defer wg.Done()
		s.Start()
	}(&server)

	// launch hbfd module
	wg.Add(1)
	go func(module modules.HbfdModule) {
		defer wg.Done()
		module.DoForever()
	}(hbfdModule)

	// launch thetafd module
	wg.Add(1)
	go func(module modules.ThetafdModule) {
		defer wg.Done()
		module.DoForever()
	}(thetafdModule)

	// launch urb module
	wg.Add(1)
	go func(module modules.UrbModule) {
		defer wg.Done()
		module.DoForever()
	}(urbModule)

	// wait forever and allow modules and communication to run concurrently
	wg.Wait()
}
