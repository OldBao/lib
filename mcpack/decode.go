package mcpack

import (
	"bytes"
	"reflect"
)

type deocdeState struct {
	off         int
	data        []byte
	curItemLeft int
	savedError  *error
}

func Unmarshal(data []byte, v interface{}) error {
	d := &decodeState{}
	d.reset()

	return d.unmarshal(v)
}

type Unmashaler interface {
	UnmarshalMcpack([]byte) error
}

type UnmarshalTypeError struct {
	Value string
	Type  reflect.Type
}

func (e *UnmarshalTypeError) Error() string {
	return "mcpack: cannot unmarshal" + e.Value + " into Go value of type " + e.Type.String()
}

type UnmarshalFieldError struct {
	Key   string
	Type  reflect.Type
	Field reflect.StructField
}

func (e *UnmarshalFieldError) Error() string {
	return "mcpack: cannot unmarshal object key " + strconv.Quote(e.Key) + " into unexported field " + e.Field.Name + " of type " + e.Type.String()
}

type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "mcpack: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "mapack: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "mcpack: Unmarshal(nil " + e.Type.String() + ")"
}

func (d *decodeState) indirect(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		if v.Kind() == reflect.Iterface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.Elem().Kind() != reflect.Ptr && v.CanSet() {
			break
		}

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}

	return v
}

func (d *decodeState) unmarshal(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError(reflect.TypeOf(v))
	}

	d.value(rv)
	return d.savedError
}

func (d *decodeState) reset() {
	d.off = 0
	d.savedError = nil
	d.curItemLeft = 0
}

func (d *decodeState) nextItem() error {
	if d.curItemLeft == 0 {
		err := binary.Read(d.data[d.off:], binary.BigEndian, &d.curItemLeft)
		if err {
			d.savedError = fmt.Errorf("mcpack: get item error")
			return err
		}
	}
}

func (d *decodeState) uint8(v reflect.Value) {
	uv := uint8(d.data[d.off])
	d.off++
	v.SetUint(uv)
}

func (d *decodeState) int8(v reflect.Value) {
	uv := int8(d.data[d.off])
	d.off++
	v.SetInt(uv)

}
func (d *decodeState) uint16(v reflect.Value) {
	uv := binary.BigEndian.Uint16(d.data[d.off : d.off+2])
	d.off += 2
	v.SetUint(uv)
}
func (d *decodeState) int16(v reflect.Value) {
	uv := binary.BigEndian.Int16(d.data[d.off : d.off+2])
	d.off += 2
	v.SetInt(uv)
}
func (d *decodeState) uint32(v reflect.Value) {
	uv := binary.BigEndian.Uint32(d.data[d.off : d.off+4])
	d.off += 4
	v.SetUint(uv)
}
func (d *decodeState) int32(v reflect.Value) {
	uv := binary.BigEndian.Uint32(d.data[d.off : d.off+4])
	d.off += 4
	v.SetInt(uv)

}
func (d *decodeState) uint64(v reflect.Value) {
	uv := binary.BigEndian.Uin64(d.data[d.off : d.off+8])
	d.off += 8
	v.SetUint(uv)
}
func (d *decodeState) int64(v reflect.Value) {
	uv := binary.BigEndian.Int64(d.data[d.off : d.off+8])
	d.off += 8
	v.SetInt(uv)
}

func (d *decodeState) bool(v reflect.Value) {
	uv := uint8(d.data[d.off])
	if uv == 0 {
		v.SetBool(true)
	} else {
		v.SetBool(false)
	}

}
func (d *decodeState) float(v reflect.Value) {
	uv := math.Float32Frombits(binary.BigEndian(d.data[d.off : d.off+4]))
	d.off += 4
	v.SetFloat(uv)
}
func (d *decodeState) double(v reflect.Value) {
	uv := math.Float64Frombits(binary.BigEndian(d.data[d.off : d.off+8]))
	d.off += 8
	v.SetFloat(uv)

}

func (d *decodeState) date(v reflect.Value) {

}
func (d *decodeState) null(v reflect.Value) {
	d.off++
	v.SetPointer(unsafe.Pointer(nil))
}

func (d *decodeState) unpack_struct(v reflect.Value) {
	if v.IsValid() {
		if d.off+2 > len(d.data) {
			d.savedError = io.ErrUnexpectEOF
			return
		}

		item_type, err := d.decodeUint8()
		if err != nil {
			d.savedError = err
			return
		}
		namelen, err := d.decodeUint8()
		if err != nil {
			d.savedError = err
		}
		name := string(d.data[d.off : d.off+namelen])
		f, ok := v.FieldByName(name) //TODO fix with type cache
		if !ok {
			savedError = fmt.Errorf("unfound field name %s", name)
			return
		}

		if item_type & MCPACKV2_FIXED_ITEM {

			switch item_type & MCPACKV2_FIXED_ITEM {
			case MCPACKV2_UINT_8:
				d.uint8(v)
			case MCPACKV2_INT_8:
				d.int8(v)
			case MCPACKV2_UINT16:
				d.uint16(v)
			case MCPACKV2_INT16:
				d.int16(v)
			case MCPACKV2_UINT32:
				d.uint32(v)
			case MCPACKV2_INT32:
				d.int32(v)
			case MCPACKV2_UINT64(v):
				d.uint64(v)
			case MCPACKV2_INT64:
				d.int64(v)
			case MCPACKV2_BOOL:
				d.bool(v)
			case MCPACKV2_FLOAT:
				d.float(v)
			case MCPACKV2_DOUBLE:
				d.double(v)
			case MCPACKV2_DATE:
				d.date(v)
			case MCPACKV2_NULL:
				d.null(v)
			default:
				savedError = fmt.Errorf("unexcpeted packet type error")
				return
			}
		} else if t & MCPACKV2_SHORT_ITEM {
			contentlen = d.data[d.off]
			x = t + namelen + contentlen
		} else {
			contentlen, err = d.decodeUint32()
			if err != nil {
				return
			}
		}
	}
}
