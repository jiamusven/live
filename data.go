package live

import (
	"encoding/binary"
	"math"
	"reflect"

	"github.com/edwingeng/live/internal"
)

var (
	Nil Data
)

type Data struct {
	v interface{}
}

func (d Data) ToBool() bool {
	return d.v.(*internal.Data).N == 1
}

func (d Data) ToInt() int {
	return int(d.v.(*internal.Data).N)
}

func (d Data) ToInt8() int8 {
	return int8(d.v.(*internal.Data).N)
}

func (d Data) ToInt16() int16 {
	return int16(d.v.(*internal.Data).N)
}

func (d Data) ToInt32() int32 {
	return int32(d.v.(*internal.Data).N)
}

func (d Data) ToInt64() int64 {
	return d.v.(*internal.Data).N
}

func (d Data) ToUint() uint {
	v, _ := binary.Uvarint(d.v.(*internal.Data).X)
	return uint(v)
}

func (d Data) ToUint8() uint8 {
	return uint8(d.v.(*internal.Data).N)
}

func (d Data) ToUint16() uint16 {
	return uint16(d.v.(*internal.Data).N)
}

func (d Data) ToUint32() uint32 {
	return uint32(d.v.(*internal.Data).N)
}

func (d Data) ToUint64() uint64 {
	v, _ := binary.Uvarint(d.v.(*internal.Data).X)
	return v
}

func (d Data) ToFloat32() float32 {
	return math.Float32frombits(uint32(d.v.(*internal.Data).N))
}

func (d Data) ToFloat64() float64 {
	v, _ := binary.Uvarint(d.v.(*internal.Data).X)
	return math.Float64frombits(v)
}

func (d Data) ToString() string {
	return string(d.v.(*internal.Data).X)
}

func (d Data) ToBytes() []byte {
	if d.v != nil {
		return d.v.(*internal.Data).X
	} else {
		return nil
	}
}

func (d Data) V() interface{} {
	return d.v
}

func (d Data) ToProtobufObj(obj interface {
	Unmarshal([]byte) error
}) {
	if len(d.v.(*internal.Data).X) == 0 {
		return
	}
	err := obj.Unmarshal(d.v.(*internal.Data).X)
	if err != nil {
		panic(err)
	}
}

func (d Data) ToJSONObj(obj interface {
	UnmarshalJSON([]byte) error
}) {
	if len(d.v.(*internal.Data).X) == 0 {
		return
	}
	err := obj.UnmarshalJSON(d.v.(*internal.Data).X)
	if err != nil {
		panic(err)
	}
}

func (d Data) Marshal() (dAtA []byte, err error) {
	if d.v == nil {
		return nil, nil
	}
	switch v := d.v.(type) {
	case *internal.Data:
		return v.Marshal()
	default:
		panic("Marshal does not support type " + reflect.TypeOf(d.v).Name())
	}
}

func (d Data) MarshalTo(dAtA []byte) (int, error) {
	if d.v == nil {
		return 0, nil
	}
	switch v := d.v.(type) {
	case *internal.Data:
		return v.MarshalTo(dAtA)
	default:
		panic("MarshalTo does not support type " + reflect.TypeOf(d.v).Name())
	}
}

func (d Data) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	if d.v == nil {
		return 0, nil
	}
	switch v := d.v.(type) {
	case *internal.Data:
		return v.MarshalToSizedBuffer(dAtA)
	default:
		panic("MarshalToSizedBuffer does not support type " + reflect.TypeOf(d.v).Name())
	}
}

func (d Data) Size() (n int) {
	if d.v == nil {
		return 0
	}
	switch v := d.v.(type) {
	case *internal.Data:
		return v.Size()
	default:
		panic("Size does not support type " + reflect.TypeOf(d.v).Name())
	}
}

func (d Data) MarshalJSON() ([]byte, error) {
	if d.v == nil {
		return nil, nil
	}
	switch v := d.v.(type) {
	case *internal.Data:
		return v.MarshalJSON()
	default:
		panic("MarshalJSON does not support type " + reflect.TypeOf(d.v).Name())
	}
}
