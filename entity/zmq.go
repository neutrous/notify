//
// ZMQ implementation of communication environment
//
// Author: neutrous
// Required: github.com/alecthomas/gozmq
//
package entity

import (
	"errors"
	"os"
	"os/signal"

	zmq "github.com/alecthomas/gozmq"
)

// Uses zmq to implement the interface entityCommEnv
type ZMQContext struct {
	ctx       *zmq.Context
	entities  map[*Endpoint]bool
	handleSig bool
}

// Initialize the zmq context
func CreateZMQCommEnv(handleSig bool) (CommEnv, error) {
	retval := &ZMQContext{}
	var err error
	retval.ctx, err = zmq.NewContext()
	if err != nil {
		return nil, err
	}
	retval.entities = make(map[*Endpoint]bool)
	retval.handleSig = handleSig
	if handleSig {
		sig := make(chan os.Signal)
		signal.Notify(sig)

		go retval.handleSignal(sig)
	}
	return retval, nil
}

func (zmq *ZMQContext) addEntity(obj *Endpoint) error {
	if _, ok := zmq.entities[obj]; ok {
		return errors.New("Already contained this entity.")
	}

	// Attempt to initialize the endpoint
	var err error
	obj.sock, err = zmq.ctx.NewSocket(obj.tp)
	if err != nil {
		obj.err = err
		return err
	}

	obj.ctx = zmq
	zmq.entities[obj] = true
	return nil
}

func (zmq *ZMQContext) removeEntity(obj *Endpoint) {
	if _, ok := zmq.entities[obj]; ok {
		if obj.sock != nil {
			obj.sock.Close()
		}
		delete(zmq.entities, obj)
	}
}

func (zmq *ZMQContext) handleSignal(sig chan os.Signal) {
	for {
		<-sig
		// if val.String() == os.Interrupt.String() ||
		// 	val.String() == os.Kill.String() {
		zmq.clearAndDestroy()
		// }
		return
	}
}

func (zmq *ZMQContext) Close() {
	if zmq.handleSig {
		os.Kill.Signal()
	} else {
		zmq.clearAndDestroy()
	}
}

func (zmq *ZMQContext) clearAndDestroy() {
	if zmq.ctx != nil {
		for key, _ := range zmq.entities {
			key.sock.Close()
		}
		// remove all entities
		zmq.entities = nil
		zmq.ctx.Close()
	}
}
