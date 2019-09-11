package main

import (
	"fmt"
	"log"
	"os"
	"self-stabilizing-uniform-reliable-broadcast/helpers"
	"self-stabilizing-uniform-reliable-broadcast/modules/urb"
	"strconv"
	"sync"
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

	// parse hosts and build P which is a slice of all node ids
	hosts, _ := helpers.ParseHostsFile()
	P := []int{}
	for _, p := range hosts {
		P = append(P, p.ID)
	}

	// init module
	urbModule := urb.UrbModule{ID: id, P: P}

	// init waitgroup to keep track of all goroutines
	var wg sync.WaitGroup

	// setup communication

	// launch urb module in a goroutine
	wg.Add(1)
	go func(module urb.UrbModule) {
		defer wg.Done()
		log.Printf("Starting URB module")
		module.DoForever()
	}(urbModule)

	// wait forever and allow modules and communication to run concurrently
	wg.Wait()
}
