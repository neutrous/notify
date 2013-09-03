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

// Indicates the level of logging function.
var logLevel = 1
const (
	Critical = iota
	Warning
	Notify
	Debug
)

// The specified lvl indicates which levels to be logged. In other
// words, if lvl is Warning, only Critical and Warning level loggs
// would be record.
func SetLoggingLevel(lvl int) { 
	if lvl < Critical { 
		lvl = Critical 
	} else if lvl > Debug { 
		lvl = Debug 
	} 
	logLevel = lvl 
}

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
	tpstr string					// Indicates the type name of entity
	tp   zmq.SocketType
	ctx  CommEnv
}

// Determine the different action in context
type action func(addr Address) error

// Tell the instance to initialize the information.
func (ep *Endpoint) initial(context CommEnv,
	act action) error {

	if ep.sock != nil {
		return errors.New("Instance has already been initialized.")
	}

	// Check the pointer's validation.
	if context == nil {
		return errors.New("Specified context is null.")
	}

	err := context.addEntity(ep)
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
	ep.Destroy()
	return err
}

// Destroy the initialized publisher instance.
func (ep *Endpoint) Destroy() {
	if ep.ctx != nil {
		ep.ctx.removeEntity(ep)
	}
}

func (ep *Endpoint) connect(addr Address) (err error) {
	err = ep.sock.Connect(string(addr))
	if err != nil && logLevel >= Debug {
		log.Printf("%s bind on the address %s failure.\n",
			ep.tpstr, addr)
	}
	return
}

func (ep *Endpoint) bind(addr Address) (err error) {
	err = ep.sock.Bind(string(addr))
	if err != nil && logLevel >= Debug {
		log.Printf("%s bind on the address %s failure.\n",
			ep.tpstr, addr)
	}
	return
}

// Append relevant addresses into publisher
func (obj *Endpoint) AppendAddress(addr Address) {
	obj.Addr = append(obj.Addr, addr)
}
