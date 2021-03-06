package live

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"reflect"
	"strings"

	"github.com/edwingeng/live/internal"
)

var (
	UnsafeMode bool
)

type Helper struct {
	Blacklist blacklist
}

var (
	liveDataType = reflect.TypeOf(Nil)
)

type blacklist []string

func (bl blacklist) covers(pkgPath string) bool {
	for _, x := range bl {
		if strings.HasPrefix(pkgPath, x) {
			if n := len(x); pkgPath[n] == '/' {
				return true
			}
		}
	}
	return false
}

func NewHelper(blacklist []string) Helper {
	var h Helper
	for _, v := range blacklist {
		if v := strings.TrimRight(v, `/`); v != "" {
			h.Blacklist = append(h.Blacklist, v)
		}
	}
	return h
}

func (h Helper) WrapBool(v bool) Data {
	if v {
		return Data{v: &internal.Data{
			N: 1,
		}}
	} else {
		return Data{v: &internal.Data{
			N: 0,
		}}
	}
}

func (h Helper) WrapInt(v int) Data {
	return Data{v: &internal.Data{
		N: int64(v),
	}}
}

func (h Helper) WrapInt8(v int8) Data {
	return Data{v: &internal.Data{
		N: int64(v),
	}}
}

func (h Helper) WrapInt16(v int16) Data {
	return Data{v: &internal.Data{
		N: int64(v),
	}}
}

func (h Helper) WrapInt32(v int32) Data {
	return Data{v: &internal.Data{
		N: int64(v),
	}}
}

func (h Helper) WrapInt64(v int64) Data {
	return Data{v: &internal.Data{
		N: v,
	}}
}

func (h Helper) WrapUint(v uint) Data {
	switch v {
	case 0:
		return Data{v: &internal.Data{}}
	default:
		var buf [10]byte
		n := binary.PutUvarint(buf[:], uint64(v))
		return Data{v: &internal.Data{
			X: buf[:n],
		}}
	}
}

func (h Helper) WrapUint8(v uint8) Data {
	return Data{v: &internal.Data{
		N: int64(v),
	}}
}

func (h Helper) WrapUint16(v uint16) Data {
	return Data{v: &internal.Data{
		N: int64(v),
	}}
}

func (h Helper) WrapUint32(v uint32) Data {
	return Data{v: &internal.Data{
		N: int64(v),
	}}
}

func (h Helper) WrapUint64(v uint64) Data {
	switch v {
	case 0:
		return Data{v: &internal.Data{}}
	default:
		var buf [10]byte
		n := binary.PutUvarint(buf[:], v)
		return Data{v: &internal.Data{
			X: buf[:n],
		}}
	}
}

func (h Helper) WrapFloat32(v float32) Data {
	b := math.Float32bits(v)
	return Data{v: &internal.Data{
		N: int64(b),
	}}
}

func (h Helper) WrapFloat64(v float64) Data {
	b := math.Float64bits(v)
	var buf [10]byte
	n := binary.PutUvarint(buf[:], b)
	return Data{v: &internal.Data{
		X: buf[:n],
	}}
}

func (h Helper) WrapString(v string) Data {
	return Data{v: &internal.Data{
		X: []byte(v),
	}}
}

func (h Helper) WrapBytes(v []byte) Data {
	return Data{v: &internal.Data{
		X: v,
	}}
}

func (h Helper) WrapProtobufObj(v interface{}) Data {
	x, ok := v.(interface {
		Marshal() ([]byte, error)
		Unmarshal([]byte) error
	})
	if !ok {
		panic("v is not protobuf compatible")
	}
	if x != nil {
		bts, err := x.Marshal()
		if err != nil {
			panic(err)
		}
		return Data{v: &internal.Data{
			X: bts,
		}}
	} else {
		return Data{v: &internal.Data{}}
	}
}

func (h Helper) WrapJSONObj(v interface{}) Data {
	x, ok := v.(interface {
		MarshalJSON() ([]byte, error)
		UnmarshalJSON([]byte) error
	})
	if ok {
		if x != nil {
			bts, err := x.MarshalJSON()
			if err != nil {
				panic(err)
			}
			return Data{v: &internal.Data{
				X: bts,
			}}
		} else {
			return Data{v: &internal.Data{}}
		}
	} else {
		bts, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		return Data{v: &internal.Data{
			X: bts,
		}}
	}
}

func (h Helper) FromProtobufBytes(v []byte) Data {
	var d internal.Data
	if err := d.Unmarshal(v); err != nil {
		panic(err)
	}
	return Data{v: &d}
}

func (h Helper) FromJSONBytes(v []byte) Data {
	var d internal.Data
	if err := d.UnmarshalJSON(v); err != nil {
		panic(err)
	}
	return Data{v: &d}
}

func (h Helper) WrapValue(v interface{}) Data {
	if v == nil {
		return Data{}
	}
	if UnsafeMode {
		return Data{v: v}
	}
	h.checkType(reflect.TypeOf(v))
	return Data{v: v}
}

func (h Helper) checkType(t reflect.Type) {
	if t == liveDataType {
		return
	}

	switch t.Kind() {
	case reflect.Bool:
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Uintptr:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.Array:
		h.checkType(t.Elem())
	case reflect.Chan:
		h.checkType(t.Elem())
	case reflect.Func:
		panic("live data does not support func")
	case reflect.Interface:
		panic("live data does not support interface")
	case reflect.Map:
		h.checkMapKeyType(t.Key())
		h.checkType(t.Elem())
	case reflect.Ptr:
		h.checkType(t.Elem())
	case reflect.Slice:
		h.checkType(t.Elem())
	case reflect.String:
	case reflect.Struct:
		pkgPath := t.PkgPath()
		if len(h.Blacklist) > 0 {
			if h.Blacklist.covers(pkgPath) {
				panic(pkgPath + " is in the blacklist of live.Helper")
			}
		}

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if live, ok := f.Tag.Lookup("live"); ok {
				if live == "true" || live == "1" {
					continue
				}
			}
			h.checkType(f.Type)
		}
	case reflect.UnsafePointer:
		panic("live data does not support unsafe pointer")
	default:
		panic("impossible")
	}
}

func (h Helper) checkMapKeyType(t reflect.Type) {
	switch t.Kind() {
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.String:
	case reflect.Uintptr:
	default:
		panic("unsupported map key type: " + t.Name())
	}
}
