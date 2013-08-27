//
// The node implements proxy/broker for internal messages
// transmitting.
//
// Every external/internal subscribing clients should connect
// to the address using tcp protocols.
//
// Usage: proxy pub_address[:port] [sub_address[:port]]
//
// If port not specified the default value(6000 for pub) would be
// used.
//
// Author: neutrous
//
package main

import (
	"fmt"
	"log"
	"os"

	zmq "github.com/alecthomas/gozmq"
	"os/signal"
)

// The proxy is used to decouple the concrete service node from any
// connected clients.
func main() {

	initOk := make(chan bool)

	inter_chan := make(chan os.Signal)
	signal.Notify(inter_chan)
	
	// Launch the proxy logic
	go proxyRouting(initOk)

	select {
	case value := <-initOk:
		if value != true {
			return
		}
	case <-inter_chan:
		log.Println("Proxy exit gracefully...")
	}
}

func proxyRouting(status chan bool) {
		// intialize the zmq context.
	context, err := zmq.NewContext()
	if err != nil {
		status <-false
		log.Fatal("Intialize the zeromq context failure.\n")
	}
	defer context.Close()

	var subscriber, publisher *zmq.Socket
	subscriber, err = context.NewSocket(zmq.XSUB)
	if err != nil {
		status <-false
		log.Fatal("Intialize the subscriber failure.\n")
	}
	defer subscriber.Close()

	var (
		sub_address, pub_address = "*", "*"
		subPort, pubPort         = 6001, 6000
	)

	// Bind the subscriber
	address := fmt.Sprintf("tcp://%s:%v", sub_address, subPort)
	err = subscriber.Bind(address)
	if err != nil {
		status <-false
		log.Fatalf("Subscriber bind on the address %s failure\n", address)
	}
	log.Printf("Subscriber bind on the address %s.\n", address)

	publisher, err = context.NewSocket(zmq.XPUB)
	if err != nil {
		status <-false
		log.Fatal("Intialize the publisher failure.\n")
	}
	defer publisher.Close()

	// Bind the publisher
	address = fmt.Sprintf("tcp://%s:%v", pub_address, pubPort)
	err = publisher.Bind(address)
	if err != nil {
		status <-false
		log.Fatalf("Publisher bind on the address %s failure.\n", address)
	}
	log.Printf("Publisher bind on the address %s.\n", address)

	log.Println("Proxy successfully launched...")
	// Poll the events on relevant sockets.
	zmq.Proxy(subscriber, publisher, nil)
}
