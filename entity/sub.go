//
// Subscription side's relevant entities.
//
// Author: neutrous
// Requires: github.com/alecthomas/gozmq
//
package entity

import (
	"errors"
	"log"

	zmq "github.com/alecthomas/gozmq"
)

const (
	SubName = "Subscriber"
)

// Deserializer defines the required methods to demarshal the
// specified raw bytes array to application relevant datas.
type Deserializer interface {
	// The typename of which the interface support to decode.
	Name() string
	// Callback, the subscriber would invoke this method to give
	// the control right.
	DataReceived(inst interface{}, ok bool)
	// Deserialize the data
	Deserialize(data []byte) (interface{}, bool)
}

// Subscriber is a abstraction for subscriber
type Subscriber struct {
	endpoint
	subs map[string]Deserializer
}

// NewSubscriber creates a initialized instance of Subscriber.
// WARNING!!! User should call this method to create a new instance
// of Subscriber, not do it by himself.
func NewSubscriber() *Subscriber {
	return &Subscriber{endpoint: endpoint{tpstr: SubName, tp: zmq.SUB},
		subs: make(map[string]Deserializer)}
}

// Subscribe the specified data type.
func (sub *Subscriber) Subscribe(inst Deserializer) error {
	if inst == nil {
		return errors.New("Couldn't subscribe nil data.")
	}

	if _, ok := sub.subs[inst.Name()]; ok {
		return errors.New("Specified data has already been subscribed.")
	}

	if err := sub.initiaFailure(); err != nil {
		return err
	}

	// Change the subscribe behavior
	if err := sub.sock.SetSockOptString(zmq.SUBSCRIBE,
		inst.Name()); err != nil {
		return err
	}
	sub.subs[inst.Name()] = inst

	return nil
}

// Unsubscribe the specified data type
func (sub *Subscriber) UnSubscribeByName(tpname string) error {
	if _, ok := sub.subs[tpname]; !ok {
		return errors.New("Specified data hasn't been subscribed.")
	}

	if err := sub.initiaFailure(); err != nil {
		return err
	}

	// Change the subscription behavior
	if err := sub.sock.SetSockOptString(zmq.UNSUBSCRIBE,
		tpname); err != nil {
		return err
	}
	delete(sub.subs, tpname)

	return nil
}

// Overload of unsubscription function.
func (sub *Subscriber) UnSubscribe(inst Deserializer) error {
	return sub.UnSubscribeByName(inst.Name())
}

// Data receiving event, which will be blocked.
func (sub *Subscriber) ReceivingEvent() error {
	if err := sub.initiaFailure(); err != nil {
		return err
	}

	// Loop for subscribing packages.
	for {
		msgbytes, err := sub.sock.RecvMultipart(0)
		if err != nil {
			if err == zmq.ENOTSOCK || err == zmq.ETERM {
				// socket has been closed or env terminated.
				log.Println("ReceivingEvent gracefully terminated.")
				return nil
			}
			log.Println("ReceivingEvent failure: ", err)
			return err
		}
		// Currently, we known every packages have only two parts.
		// The first part indicates the type name;
		// The reset part contains the whole package contents.
		if len(msgbytes) != 2 {
			log.Println("ReceivingEvent failure: Unknown data package part.")
			continue
		}
		if inst, ok := sub.subs[string(msgbytes[0])]; ok {
			inst.DataReceived(inst.Deserialize(msgbytes[1]))
		} else {
			log.Println("ReceivingEvent failure: Doesn't support data type.")
		}
	}
}

func (sub *Subscriber) initiaFailure() error {
	if sub.sock == nil || sub.err != nil {
		return errors.New("Subscriber instance hasn't been initialized.")
	}
	return nil
}
