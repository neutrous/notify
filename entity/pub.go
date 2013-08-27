//
// Publication side's relevant entities.
//
// Author: neutrous
// Requires: github.com/alecthomas/gozmq
//
package entity

import (
	"errors"
	
	zmq "github.com/alecthomas/gozmq"
)

const (
	PubName = "Publisher"
)

// Abstraction for publisher
type Publisher struct {
	Endpoint
}

// Uses connecting role to intialze the publisher instance.
func (pub *Publisher) InitialConnecting(context *zmq.Context) error {
	pub.tp = PubName
	return pub.initial(context, pub.connect, zmq.PUB)
}

// Uses binding role to intialize the publisher instance.
func (pub *Publisher) InitialBinding(context *zmq.Context) error {
	pub.tp = PubName
	return pub.initial(context, pub.bind, zmq.PUB)
}

// Send the specified data, the data must be serializable
func (pub *Publisher) Write(data Serializable) error {
	if pub.sock == nil || pub.err != nil {
		return errors.New("Publisher hasn't been initialized.")
	}

	// Construct the data header/name
	parts := make([][]byte, 2)
	parts[0] = make([]byte, len(data.Name()))
	copy(parts[0], data.Name())
	
	// Construct the data content
	content, err := data.Serialize()
	if err != nil {
		return err
	}
	parts[1] = make([]byte, len(content))
	copy(parts[1], content)

	return pub.sock.SendMultipart(parts, 0)
}

