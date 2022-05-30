package msgpack

import (
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/vmihailenco/msgpack/v5"
)

// Name is the name registered for the msgpack codec.
const Name = "msgpack"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with json.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}