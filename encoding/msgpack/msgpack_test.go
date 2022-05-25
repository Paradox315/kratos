package msgpack

import (
	"github.com/go-kratos/kratos/v2/encoding"
	"testing"
)

const contentType = "msgpack"

type testModel struct {
	Field1 string
	Field2 int
	Field3 bool
}

func Test_Msgpack(t *testing.T) {
	m := map[string]float32{
		"pi": 0,
		"e":  0,
	}
	t.Log(m)
	co := encoding.GetCodec(contentType)
	bytes, _ := co.Marshal(m)
	var m2 map[string]float32
	_ = co.Unmarshal(bytes, &m2)
	t.Log(m2)
}
