//
// The program of subscriber
//
// Author: neutrous
// Requires: github.com/alecthomas/gozmq
//
package main

import (
	"log"

	"github.com/neutrous/notify/entity"
	"github.com/neutrous/notify/example/entity_tests/model"
)

const (
	filter = 28
)

func main() {

	// Initialize subscriber
	sub := entity.NewSubscriber(filter)
	sub.AppendAddress("tcp://localhost:6602")

	context, _ := entity.CreateZMQCommEnv(true)
	defer context.Close()

	if err := sub.InitialConnecting(context); err != nil {
		log.Fatalln("Initialize subscriber failure.", err)
	}

	value := &example.TestHandler{}

	if err := sub.Subscribe(value); err != nil {
		log.Fatalln("Subscribe example.Test data type failure.")
	}

	log.Println("Start receiving datas.")

	sub.ReceivingEvent()
}
