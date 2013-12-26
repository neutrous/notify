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

func bindingPub(addr Address, context CommEnv) (*Publisher, error) {
	inst := NewPublisher()
	inst.AppendAddress(addr)
	return inst, inst.InitialBinding(context)
}

// tcp address should bindable.
func TestBindingTcp(t *testing.T) {
	context, _ := CreateZMQCommEnv(false)
	defer context.Close()

	if _, err := bindingPub("tcp://*:6601", context); err != nil {
		t.Errorf("Address should bind on %s, while return %s\n",
			string("tcp://*:6601"), err)
	}
}

// unknown address type shouldn't be bindable
func TestBindingUnknownType(t *testing.T) {
	context, _ := CreateZMQCommEnv(false)
	defer context.Close()

	if _, err := bindingPub("unk://*:231", context); err != nil {
		return
	}
	t.Errorf("Specified address [unk://*:231] is " +
		"not available, but return value indicates accepted.\n")

}

func TestPubClosable(t *testing.T) {
	context, _ := CreateZMQCommEnv(false)
	defer context.Close()

	if inst, err := bindingPub("tcp://*:6621", context); err == nil {
		inst.Close()
		return
	}

	t.Errorf("Initialized publisher instance should be closable.\n")
}

func TestSendNilInterface(t *testing.T) {
	context, _ := CreateZMQCommEnv(false)
	defer context.Close()

	if inst, err := bindingPub("tcp://*:6621", context); err == nil {
		if inst.Send(nil) == nil {
			t.Errorf("When sending nil data, error should be returned.\n")
		}
		return
	}

	t.Errorf("Initialized failure.\n")
}

func TestSendSerializedData(t *testing.T) {
	context, _ := CreateZMQCommEnv(false)
	defer context.Close()

	if inst, err := bindingPub("tcp://*:6621", context); err == nil {
		gen := func() Serializer {
			return &testData{}
		}
		if inst.Send(1, gen()) != nil {
			t.Errorf("Sending failure.\n")
		}
		return
	}

	t.Errorf("Initialized failure.\n")
}
