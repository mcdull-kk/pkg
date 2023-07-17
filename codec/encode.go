package codec

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-playground/form/v4"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const (
	xmlName   = "xml"
	protoName = "proto"
	jsonName  = "json"
	yamlName  = "yaml"
	formName  = "x-www-form-urlencoded"
)

func init() {
	formDecoder.SetTagName("json")
	formEncoder.SetTagName("json")
	registerCodec(xmlCodec{})
	registerCodec(jsonCodec{})
	registerCodec(protoCodec{})
	registerCodec(yamlCodec{})
	registerCodec(formCodec{encoder: formEncoder, decoder: formDecoder})
}

var (
	formEncoder = form.NewEncoder()
	formDecoder = form.NewDecoder()
	// MarshalOptions is a configurable JSON format marshaller.
	MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true,
	}
	// UnmarshalOptions is a configurable JSON format parser.
	UnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

// Codec defines the interface Transport uses to encode and decode messages.  Note
// that implementations of this interface must be thread safe; a Codec's
// methods can be called from concurrent goroutines.
type Codec interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
	// Name returns the name of the Codec implementation. The returned string
	// will be used as part of content type in transmission.  The result must be
	// static; the result cannot change between calls.
	Name() string
}

var registeredCodecs = make(map[string]Codec)

// registerCodec registers the provided Codec for use with all Transport clients and
// servers.
func registerCodec(codec Codec) {
	if codec == nil {
		panic("cannot register a nil Codec")
	}
	if codec.Name() == "" {
		panic("cannot register Codec with empty string result for Name()")
	}
	contentSubtype := strings.ToLower(codec.Name())
	registeredCodecs[contentSubtype] = codec
}

// GetCodec gets a registered Codec by content-subtype, or nil if no Codec is
// registered for the content-subtype.
//
// The content-subtype is expected to be lowercase.
func GetCodec(contentSubtype string) Codec {
	return registeredCodecs[contentSubtype]
}

type codec struct{}

type (
	xmlCodec   codec
	jsonCodec  codec
	protoCodec codec
	yamlCodec  codec
	formCodec  struct {
		encoder *form.Encoder
		decoder *form.Decoder
	}
)

// xmlCodec
func (xmlCodec) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (xmlCodec) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func (xmlCodec) Name() string {
	return xmlName
}

// jsonCodec
func (jsonCodec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return MarshalOptions.Marshal(m)
	default:
		return json.Marshal(m)
	}
}

func (jsonCodec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
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
		return json.Unmarshal(data, m)
	}
}

func (jsonCodec) Name() string {
	return jsonName
}

//protoCodec
func (protoCodec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (protoCodec) Unmarshal(data []byte, v interface{}) error {
	pm, err := getProtoMessage(v)
	if err != nil {
		return err
	}
	return proto.Unmarshal(data, pm)
}

func (protoCodec) Name() string {
	return protoName
}

func getProtoMessage(v interface{}) (proto.Message, error) {
	if msg, ok := v.(proto.Message); ok {
		return msg, nil
	}
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return nil, errors.New("not proto message")
	}

	val = val.Elem()
	return getProtoMessage(val.Interface())
}

// yamlCodec
func (yamlCodec) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (yamlCodec) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func (yamlCodec) Name() string {
	return yamlName
}

//formCodec
func (c formCodec) Marshal(v interface{}) ([]byte, error) {
	var vs url.Values
	var err error

	vs, err = c.encoder.Encode(v)
	if err != nil {
		return nil, err
	}

	for k, v := range vs {
		if len(v) == 0 {
			delete(vs, k)
		}
	}
	return []byte(vs.Encode()), nil
}

func (c formCodec) Unmarshal(data []byte, v interface{}) error {
	vs, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}

	return c.decoder.Decode(v, vs)
}

func (formCodec) Name() string {
	return formName
}
