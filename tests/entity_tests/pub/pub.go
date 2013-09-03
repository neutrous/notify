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
	
	"github.com/neutrous/notify/entity"
	"github.com/neutrous/notify/tests/entity_tests/model"
	"code.google.com/p/goprotobuf/proto"
)

func main() {
	
	// Initialize publisher
	pub := entity.Publisher{}
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
		if err := pub.Send(&value); err != nil {
			log.Println("Publish value failure.")
			break
		} else {
			log.Println("Publish value okay.")
		}
	}
	
}
