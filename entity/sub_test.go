//
// Testing for subscriber component
//
// Author: neutrous
//
package entity_test

import (
	"testing"

	. "github.com/neutrous/notify/entity"
)

func connectSub(addr Address, context CommEnv) (*Subscriber, error) {
	inst := NewSubscriber(int32(1))
	inst.AppendAddress(addr)

	return inst, inst.InitialConnecting(context)
}

func TestConnectingUnknown(t *testing.T) {
	context, _ := CreateZMQCommEnv(false)
	defer context.Close()

	if _, err := connectSub("unk://*:6632", context); err == nil {
		t.Errorf("Specified address [unk://*:6632] shouldn't conncted.\n")
	}
}

func TestSubClosable(t *testing.T) {
	context, _ := CreateZMQCommEnv(false)
	defer context.Close()

	if inst, err := connectSub("tcp://127.0.0.1:6621", context); err == nil {
		inst.Close()
	} else {
		t.Errorf("Initialized failure.\n")
	}
}

func TestSubscribeData(t *testing.T) {
	context, _ := CreateZMQCommEnv(false)
	defer context.Close()

	if inst, err := connectSub("tcp://127.0.0.1:6621", context); err == nil {
		if inst.Subscribe(nil) == nil {
			t.Errorf("Shouldn't be subscribing nil data.\n")
		}
		if inst.Subscribe(&testData{}) != nil {
			t.Errorf("Subscribe failure.\n")
		}
		if inst.UnSubscribe(&testData{}) != nil {
			t.Errorf("Unsubscribe by instance failure.\n")
		}
		// test another unsubscribe
		data := &testData{}
		if inst.Subscribe(data) != nil {
			t.Errorf("Subscribe failure.\n")
		}
		if inst.UnSubscribeByName(data.Name()) != nil {
			t.Errorf("UnSubscribe by name failure.\n")
		}
	} else {
		t.Errorf("Initialized failure.\n")
	}
}
