//go:build !((linux || darwin) && amd64)

package json

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/encoding"
	gojson "github.com/goccy/go-json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Name is the name registered for the json codec.
const Name = "json"

var (
	// MarshalOptions is a configurable JSON format marshaller.
	MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true,
	}
	// UnmarshalOptions is a configurable JSON format parser.
	UnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with json.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return gojson.Marshal(m)
	case proto.Message:
		return MarshalOptions.Marshal(m)
	default:
		return gojson.Marshal(m)
	}
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return gojson.Unmarshal(data, m)
	case proto.Message:
		return UnmarshalOptions.Unmarshal(data, m)
	default:
		return gojson.Unmarshal(data, m)
	}
}

func (codec) Name() string {
	return Name
}
