//go:build darwin || freebsd || linux
// +build darwin freebsd linux

package json

import (
	"encoding/json"
	jsonC "github.com/bytedance/sonic"
	"reflect"

	"github.com/go-kratos/kratos/v2/encoding"
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
		return jsonC.Marshal(m)
	case proto.Message:
		return MarshalOptions.Marshal(m)
	default:
		return jsonC.Marshal(m)
	}
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return jsonC.Unmarshal(data, m)
	case proto.Message:
		return UnmarshalOptions.Unmarshal(data, m)
	default:
		rv := reflect.ValueOf(v)
		for rv := rv; rv.Kind() == reflect.Ptr; {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		if m, ok := reflect.Indirect(rv).Interface().(proto.Message); ok {
			return UnmarshalOptions.Unmarshal(data, m)
		}
		return jsonC.Unmarshal(data, m)
	}
}

func (codec) Name() string {
	return Name
}
