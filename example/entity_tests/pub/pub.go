//
// The program of publisher
//
// Author: neutrous
// Requires: github.com/alecthomas/gozmq
//
package main

import (
	"log"
	"time"

	"code.google.com/p/goprotobuf/proto"
	"github.com/neutrous/notify/entity"
	"github.com/neutrous/notify/example/entity_tests/model"
)

const (
	filter = 28
)

func main() {

	// Initialize publisher
	pub := entity.NewPublisher()
	pub.AppendAddress("tcp://*:6602")

	context, _ := entity.CreateZMQCommEnv(true)
	defer context.Close()

	if err := pub.InitialBinding(context); err != nil {
		log.Fatalln("Initialize publisher failure.", err)
	}

	value := example.Test{}
	value.Label = proto.String("value1")
	value.Type = proto.Int32(0)

	for {
		time.Sleep(time.Second)
		*value.Type = value.GetType() + 2
		if err := pub.Send(filter, &value); err != nil {
			log.Println("Publish value failure.", err)
			break
		} else {
			log.Println("Publish value okay.")
		}
	}

}
