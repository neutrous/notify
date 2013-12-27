//
// The file contains test cases to detect the functions of
// notify component.
//

package entity_test

import (
	"log"
	"math/rand"
	"testing"
	"time"

	. "github.com/neutrous/notify/entity"
)

const (
	filter1 = 0x12345678
	filter2 = 0x78563412
)

type testStruct struct {
	ReceivedCount int
}

func (this *testStruct) Name() string {
	return "testStruct"
}

func (this *testStruct) Serialize() ([]byte, error) {
	return []byte{}, nil
}

func (this *testStruct) DataReceived(inst interface{}, ok bool) {
	this.ReceivedCount++
}

func (this *testStruct) Deserialize(data []byte) (interface{}, bool) {
	return nil, true
}

func TestFilterFunction(t *testing.T) {
	context, _ := CreateZMQCommEnv(true)

	// Create a publisher
	pub := NewPublisher()
	pub.AppendAddress("tcp://*:6622")

	if pub.InitialBinding(context) != nil {
		t.Errorf("Publisher binding failure.\n")
	}

	ret1, ret2 := &testStruct{}, &testStruct{}
	go subscribe(context, filter1, ret1, "tcp://localhost:6622")
	go subscribe(context, filter2, ret2, "tcp://localhost:6622")

	// for connect between each other
	time.Sleep(time.Second)

	for idx := 0; idx < 5; idx++ {
		if pub.Send(filter1, &testStruct{}) != nil {
			t.Errorf("pub send filter1 failure.\n")
		}
	}

	for idx := 0; idx < 10; idx++ {
		if pub.Send(filter2, &testStruct{}) != nil {
			t.Errorf("pub send filter2 failure.\n")
		}
	}

	time.Sleep(time.Second)
	// Close the context to force pub and sub sides to quit.
	context.Close()

	if ret1.ReceivedCount != 5 {
		t.Errorf("filter1 received count failure, expected 5 while got %d\n",
			ret1.ReceivedCount)
	}

	if ret2.ReceivedCount != 10 {
		t.Errorf("filter2 received count failure, expected 10 while got %d\n",
			ret2.ReceivedCount)
	}

}

func subscribe(context CommEnv, filter int, ret *testStruct, addr Address) {
	sub := NewSubscriber(filter)
	sub.AppendAddress(addr)
	if sub.InitialConnecting(context) != nil {
		log.Printf("Subscriber at %d connecting failure.\n", filter)
		return
	}

	if sub.Subscribe(ret) != nil {
		log.Printf("Subscribe at %d failure.\n", filter)
		return
	}

	sub.ReceivingEvent()
}

func TestCommEnv(t *testing.T) {
	context1, _ := CreateZMQCommEnv(false)
	context2, _ := CreateZMQCommEnv(false)

	var pub_addr, sub_addr Address = "tcp://*:6623", "tcp://localhost:6623"

	pub := NewPublisher()
	pub.AppendAddress(Address(pub_addr))

	if pub.InitialBinding(context1) != nil {
		t.Errorf("Pubisher initial failure.\n")
	}

	ret1, ret2 := &testStruct{}, &testStruct{}

	go subscribe(context1, filter1, ret1, sub_addr)
	go subscribe(context2, filter1, ret2, sub_addr)

	// wait for connecting
	time.Sleep(time.Second)

	// context2 directly close to test whether it could receive datas.
	context2.Close()
	for idx := 0; idx < 10; idx++ {
		if pub.Send(filter1, &testStruct{}) != nil {
			t.Errorf("Send failure.\n")
		}
	}

	// wait for receiving
	time.Sleep(time.Millisecond * 10)
	context1.Close()

	if ret1.ReceivedCount != 10 {
		t.Errorf("Should received 10 packages while got %d.\n", ret1.ReceivedCount)
	}

	if ret2.ReceivedCount != 0 {
		t.Errorf("Should receive none package while got %d.\n", ret2.ReceivedCount)
	}
}

func TestSendMultiPackages(t *testing.T) {
	context, _ := CreateZMQCommEnv(false)

	pub := NewPublisher()
	pub.AppendAddress("tcp://*:6624")

	data := &testStruct{}

	if pub.InitialBinding(context) != nil {
		t.Errorf("Publication initialized failuer.\n")
	}

	go subscribe(context, filter1, data, "tcp://localhost:6624")

	// Wait for connection.
	time.Sleep(time.Second)

	totalcount := 0
	for idx := 0; idx < 10; idx++ {
		subpackages_count := rand.Intn(5)
		totalcount += subpackages_count
		datas := make([]Serializer, subpackages_count)
		for i := 0; i < subpackages_count; i++ {
			datas[i] = &testStruct{}
		}
		if pub.Send(filter1, datas...) != nil {
			t.Errorf("Send multipackages failure.\n")
		}
	}

	// Wait for receiving.
	time.Sleep(time.Second)

	context.Close()
	if totalcount != data.ReceivedCount {
		t.Errorf("Expect receive %d packages, while got %d.\n",
			totalcount, data.ReceivedCount)
	}

}
