//
// Common definitions of entities module.
//
// Author: neutrous
//
package entity

import (
	"errors"
	"log"

	zmq "github.com/alecthomas/gozmq"
)

type Address string
type Addresses []Address

type Serializable interface {
	// The typename used to identify the data
	Name() string
	// Serialize the data into byte array
	Serialize() ([]byte, error)
}

// Indicate whether specified addr is valid.
func (addr *Address) valid() bool {
	// TODO: uses reg to valid the specified address.
	return true
}

// Abstraction for I/O on the net's entities.
type Endpoint struct {
	Addr Addresses
	sock *zmq.Socket
	err  error
	tp   string					// Indicates the type name of entity
}

// Determine the different action in context
type action func(addr Address) error

// Tell the instance to initialize the information.
func (ep *Endpoint) initial(context *zmq.Context,
	act action, tp zmq.SocketType) error {

	if ep.sock != nil {
		return errors.New("Instance has already been initialized.")
	}

	// Check the pointer's validation.
	if context == nil {
		return errors.New("Specified context is null.")
	}

	var err error
	ep.sock, err = context.NewSocket(tp)
	if err != nil {
		return err
	}

	// Check the addrs' validation.
	for _, addr := range ep.Addr {
		if !addr.valid() {
			return ep.handleError(
				errors.New("Specified addr not valid."))
		}
		if err = act(addr); err != nil {
			return ep.handleError(err)
		}
	}

	return nil
}

func (ep *Endpoint) handleError(err error) error {
	ep.sock.Close()
	return err
}

// Destroy the initialized publisher instance.
func (ep *Endpoint) Destroy() {
	if ep.sock != nil {
		ep.sock.Close()
	}
}

func (ep *Endpoint) connect(addr Address) (err error) {
	err = ep.sock.Connect(string(addr))
	if err != nil {
		log.Printf("%s bind on the address %s failure.\n",
			ep.tp, addr)
	}
	return
}

func (ep *Endpoint) bind(addr Address) (err error) {
	err = ep.sock.Bind(string(addr))
	if err != nil {
		log.Printf("%s bind on the address %s failure.\n",
			ep.tp, addr)
	}
	return
}

// Append relevant addresses into publisher
func (obj *Endpoint) AppendAddress(addr Address) {
	obj.Addr = append(obj.Addr, addr)
}
