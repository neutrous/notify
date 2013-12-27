//
// Subscription side's relevant entities.
//
// Author: neutrous
// Requires: github.com/alecthomas/gozmq
//
package entity

import (
	"encoding/binary"
	"errors"
	"log"
	"time"

	zmq "github.com/alecthomas/gozmq"
)

const (
	SubName     = "Subscriber"
	LostTooLong = "The subscribed package has been lost too long."
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
	subs         map[string]Deserializer
	filter       []byte
	last         time.Time
	max_duration time.Duration
}

// NewSubscriber creates a initialized instance of Subscriber.
// WARNING!!! User should call this method to create a new instance
// of Subscriber, not do it by himself.
func NewSubscriber(filter interface{}) *Subscriber {
	ret := &Subscriber{endpoint: endpoint{tpstr: SubName, tp: zmq.SUB},
		subs: make(map[string]Deserializer),
		// TODO: i don't know how to specify the maximum value of
		// time.Duration, here uses 24 hour as the maximu value.
		max_duration: time.Hour * 24}
	if ret.filter = create_filter_bytes(filter); ret.filter == nil {
		return nil
	}
	return ret
}

func (sub *Subscriber) InitialConnecting(context CommEnv) error {
	if err := sub.endpoint.InitialConnecting(context); err != nil {
		return err
	}
	// Change the subscribe behavior
	if err := sub.sock.SetSockOptUInt64(zmq.UInt64SocketOption(zmq.SUBSCRIBE),
		binary.LittleEndian.Uint64(sub.filter)); err != nil {
		return err
	}
	return nil
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

	delete(sub.subs, tpname)

	return nil
}

// Overload of unsubscription function.
func (sub *Subscriber) UnSubscribe(inst Deserializer) error {
	return sub.UnSubscribeByName(inst.Name())
}

// ReceivingEvent starts a data receiving event loop, which would be blocked.
// User could use relevant CommEnv instant or the subscriber instant to
// make this function return, by calling its Close Method.
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

		// The protocol should be refered to `push_service' description.
		if len(msgbytes) < 3 {
			log.Println("ReceivingEvent failure: Unknown data package part.")
			continue
		}
		// check the protocol version.
		if string(msgbytes[1]) != ProtocolName {
			log.Println("Not support protocol version.")
			continue
		}
		// check the validation of timestamp
		if !sub.last.Equal(time.Time{}) {
			current := time.Unix(int64(binary.LittleEndian.Uint64(msgbytes[2])), 0)
			if current.After(sub.last) {
				// if the received package is too long since last
				// timestamp, we have to close the sub side, because
				// it's a too slow one who may cause pub side slow.
				if current.Sub(sub.last) > sub.max_duration {
					log.Println(LostTooLong)
					sub.sock.Close()
					sub.err = errors.New(LostTooLong)
					return sub.err
				}
			} else {
				log.Println("Received unordered package.")
				continue
			}
		}
		// does the real decoding things.  Because the data may
		// aggregate as one package, we have to fectch it out one by one.
		for idx := 3; idx < len(msgbytes); idx += 2 {
			if inst, ok := sub.subs[string(msgbytes[idx])]; ok {
				inst.DataReceived(inst.Deserialize(msgbytes[idx+1]))
			} else {
				log.Println("ReceivingEvent failure: Doesn't support data type.",
					string(msgbytes[idx]))
			}
		}
	}
}

func (sub *Subscriber) initiaFailure() error {
	if sub.sock == nil || sub.err != nil {
		return errors.New("Subscriber instance hasn't been initialized.")
	}
	return nil
}

func (sub *Subscriber) SetMaximumDuration(val time.Duration) {
	sub.max_duration = val
}

func (sub *Subscriber) GetMaximumDuration() time.Duration {
	return sub.max_duration
}
