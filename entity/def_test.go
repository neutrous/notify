//
//  Base configuration for testing pub or sub
//
//  Author: neutrous
//

package entity_test

type testData struct {
}

func (this *testData) Name() string {
	return "testData"
}

func (this *testData) Serialize() ([]byte, error) {
	return []byte{}, nil
}

func (this *testData) DataReceived(inst interface{}, ok bool) {
}

func (this *testData) Deserialize(data []byte) (interface{}, bool) {
	return nil, true
}
