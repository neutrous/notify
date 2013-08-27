//
// Testing for publisher component.
//
// Author: neutrous
//
package entity_test

import (
	"testing"

	zmq "github.com/alecthomas/gozmq"
	. "github.com/neutrous/notify/entity"
)

func bindingPub(addr Address, t *testing.T) {

	inst := Publisher{}
	inst.Addr = make([]Address, 1)
	inst.Addr[0] = addr

	err := make(chan error)
	context, _ := zmq.NewContext()
	defer context.Close()

	// select {
	// case err <- inst.InitialBinding(context):
	// 	if err != nil {
	// 		t.Errorf("Address should bind on %s, while return %s",
	// 			string(inst.Addr[0]), err)
	// 	}
	// case <-time.After(time.Second):
	// 	return
	// }
	if err <-inst.InitialBinding(context); err != nil {
		t.Errorf("Address should bind on %s, while return %s\n",
			string(inst.Addr[0]), <-err)
	}

}

// tcp address should bindable.
func TestBindingTcp(t *testing.T) {
	bindingPub("tcp://localhost:6601", t)
}

// unknown address type shouldn't be bindable
func TestBindingUnknownType(t *testing.T) {
	bindingPub("unk://localhost:231", t)
}
