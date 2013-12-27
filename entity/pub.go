//
// Publication side's relevant entities.
//
// Author: neutrous
// Requires: github.com/alecthomas/gozmq
//
package entity

import (
	"encoding/binary"
	"errors"
	"time"

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
// WARNING!!! User should use this method to create a instance
// of Publisher, not by itself.
func NewPublisher() *Publisher {
	return &Publisher{endpoint{tp: zmq.PUB, tpstr: PubName}}
}

// Broadcast make the data to all of the attached subscribers.
// func (pub *Publisher) Broadcast(data ...Serializer) error {
// 	return pub.send_with_filter([]byte{}, data...)
// }

// Send make the data to the specified subscriber who aims to
// receive the data with flag of filter.
// WARNING!!! Currently, the filter should be integer value.
func (pub *Publisher) Send(filter interface{}, data ...Serializer) error {
	// create filter value every time.
	if filter_bytes := create_filter_bytes(filter); filter_bytes != nil {
		return pub.send_with_filter(filter_bytes, data...)
	}
	return errors.New("unknown type")
}

func (pub *Publisher) send_with_filter(filter []byte, data ...Serializer) error {
	if pub.sock == nil || pub.err != nil {
		return errors.New("Publisher hasn't been initialized.")
	}

	if data == nil {
		return errors.New("Couldn't send nil data.")
	}

	// Construct the data.
	// The protocol is described in `push_service', please refer to
	// it for details.
	parts := make([][]byte, 3+2*len(data))
	filter_len := len(filter)
	parts[0] = make([]byte, len(string(filter)))
	if filter_len != 0 {
		copy(parts[0], string(filter))
	}
	parts[1] = make([]byte, len(ProtocolName))
	parts[2] = make([]byte, 8)
	copy(parts[1], ProtocolName)
	binary.LittleEndian.PutUint64(parts[2], uint64(time.Now().Unix()))

	// Construct the data content
	for i, idx := 0, 0; i < len(data); i, idx = i+1, idx+2 {
		content, err := data[i].Serialize()
		if err != nil {
			return err
		}
		parts[3+idx] = make([]byte, len(data[i].Name()))
		copy(parts[3+idx], data[i].Name())
		parts[4+idx] = make([]byte, len(content))
		copy(parts[4+idx], content)
	}

	return pub.sock.SendMultipart(parts, 0)
}
