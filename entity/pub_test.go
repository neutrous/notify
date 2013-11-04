//
// Testing for publisher component.
//
// Author: neutrous
//
package entity_test

import (
	"testing"

	. "github.com/neutrous/notify/entity"
)

func bindingPub(addr Address) error {

	inst := NewPublisher()
	inst.AppendAddress(addr)
		
	context, _ := CreateZMQCommEnv(false)
	defer context.Close()

	return inst.InitialBinding(context)
}

// tcp address should bindable.
func TestBindingTcp(t *testing.T) {
	if err := bindingPub("tcp://*:6601"); err != nil {
		t.Errorf("Address should bind on %s, while return %s\n",
			string("tcp://*:6601"), err)
	}
}

// unknown address type shouldn't be bindable
func TestBindingUnknownType(t *testing.T) {
	if err := bindingPub("unk://*:231"); err != nil {
		return
	}	
	t.Errorf("Specified address [unk://*:231] is not available, but return value indicates accepted.\n")

}










