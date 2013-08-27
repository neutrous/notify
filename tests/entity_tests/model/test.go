// 
// This is the example to implements the Serializable or
// Deserializable interface.
// 
package example

import (
	"code.google.com/p/goprotobuf/proto"
	"log"
)

// Implementation of Serializable interface.
func (obj *Test) Name() string {
	return "example.Test"
}

func (obj *Test) Serialize() (data []byte, err error) {
	data, err = proto.Marshal(obj)
	return
}

type TestHandler struct {
	count int32
}

func (obj *TestHandler) Deserialize(data []byte) (inst interface{}, ok bool) {
	var err error
	testObj := &Test{}
	inst, ok = testObj, true
	err = proto.Unmarshal(data, testObj)
	if err != nil {
		ok = false
	}
	return
}

func (obj *TestHandler) DataReceived(inst interface{}, ok bool) {

	if !ok {
		log.Printf("DataReceived: specified data Unmarshaled failure.\n")
		return
	}
	
	if value, is := inst.(*Test); is {
		obj.count++
		log.Printf("[%v], %s %v\n", obj.count, value.GetLabel(),
			value.GetType())
	} else {
		log.Printf("Specified inst object isn't the type of example.Test.\n")
	}
	
}

func (obj *TestHandler) Name() string {
	return "example.Test"
}
