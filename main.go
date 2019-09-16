package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

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
	zeroedSlice := []int{}
	for _, p := range hosts {
		P = append(P, p.ID)
		zeroedSlice = append(zeroedSlice, 0)
	}

	resolver := modules.Resolver{}

	// init module
	urbModule := modules.UrbModule{ID: id, P: P, Resolver: &resolver, Seq: 0, Buffer: modules.Buffer{}, RxObsS: []int{}, TxObsS: []int{}}
	hbfdModule := modules.HbfdModule{ID: id, P: P, Resolver: &resolver, Hb: zeroedSlice}
	thetafdModule := modules.ThetafdModule{ID: id, P: P, Resolver: &resolver, Vector: zeroedSlice}

	// init resolver and attach modules
	resolver.Modules = make(map[modules.ModuleType]interface{})
	resolver.Modules[modules.URB] = urbModule
	resolver.Modules[modules.HBFD] = hbfdModule
	resolver.Modules[modules.THETAFD] = thetafdModule

	// init waitgroup to keep track of all goroutines
	var wg sync.WaitGroup

	// setup communication

	// launch hbfd module
	wg.Add(1)
	go func(module modules.HbfdModule) {
		defer wg.Done()
		log.Printf("Starting HBFD module")
		module.DoForever()
	}(hbfdModule)

	// launch thetafd module
	wg.Add(1)
	go func(module modules.ThetafdModule) {
		defer wg.Done()
		log.Printf("Starting THETAFD module")
		module.DoForever()
	}(thetafdModule)

	// launch urb module
	wg.Add(1)
	go func(module modules.UrbModule) {
		defer wg.Done()
		log.Printf("Starting URB module")
		module.DoForever()
	}(urbModule)

	// wait forever and allow modules and communication to run concurrently
	wg.Wait()
}
