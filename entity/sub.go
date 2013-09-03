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

// Abstraction for subscriber
type Subscriber struct {
	Endpoint
	subs map[string]Deserializer
}

// Uses connecting role to initialze the subscriber instance
func (sub *Subscriber) InitialConnecting(context CommEnv) error {
	sub.tpstr = SubName
	sub.tp = zmq.SUB
	sub.subs = make(map[string]Deserializer)
	return sub.initial(context, sub.connect)
}

// Uses binding role to intialze the subscriber instance.
func (sub *Subscriber) InitialBinding(context CommEnv) error {
	sub.tpstr = SubName
	sub.tp = zmq.SUB
	sub.subs = make(map[string]Deserializer)
	return sub.initial(context, sub.bind)
}

// Subscribe the specified data type.
func (sub *Subscriber) Subscribe(inst Deserializer) error {
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
			log.Println("ReceivingEvent failure: ", err)
			continue
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
