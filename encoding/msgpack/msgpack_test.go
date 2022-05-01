package msgpack

import "testing"

type testModel struct {
	Field1 string
	Field2 int
	Field3 bool
}

func Test_Msgpack(t *testing.T) {
	m := &testModel{
		Field1: "test",
		Field2: 1,
		Field3: true,
	}
	t.Log(m)
	co := codec{}
	bytes, _ := co.Marshal(m)
	var m2 *testModel
	_ = co.Unmarshal(bytes, &m2)
	t.Log(m2)
}
