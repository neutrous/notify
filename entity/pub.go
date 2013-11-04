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

// Serializer defines the required method to marshal the application
// relevant datas into byte array.
type Serializer interface {
	// The typename used to identify the data
	Name() string
	// Serialize the data into byte array
	Serialize() ([]byte, error)
}

// Abstraction for publisher
type Publisher struct {
	endpoint
}

// NewPublisher creates a initialized publisher instance.
// WARNING!!! User should all this method to create a instance
// of Publisher, not by himself.
func NewPublisher() *Publisher {
	return &Publisher{endpoint{tp: zmq.PUB, tpstr: PubName}}
}

// Send the specified data, the data must be serializable
func (pub *Publisher) Send(data Serializer) error {
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
