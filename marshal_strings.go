package qs

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

type ptrMarshaler struct {
	Type          reflect.Type
	ElemMarshaler Marshaler
}

func newPtrMarshaler(t reflect.Type, opts *MarshalOptions) (Marshaler, error) {
	if t.Kind() != reflect.Ptr {
		return nil, &wrongKindError{Expected: reflect.Ptr, Actual: t}
	}
	et := t.Elem()
	em, err := opts.MarshalerFactory.Marshaler(et, opts)
	if err != nil {
		return nil, err
	}
	return &ptrMarshaler{
		Type:          t,
		ElemMarshaler: em,
	}, nil
}

func (p *ptrMarshaler) Marshal(v reflect.Value, opts *MarshalOptions) ([]string, error) {
	t := v.Type()
	if t != p.Type {
		return nil, &wrongTypeError{Actual: t, Expected: p.Type}
	}
	if v.IsNil() {
		return nil, nil
	}
	return p.ElemMarshaler.Marshal(v.Elem(), opts)
}

type arrayAndSliceMarshaler struct {
	Type          reflect.Type
	ElemMarshaler Marshaler
}

func newArrayAndSliceMarshaler(t reflect.Type, opts *MarshalOptions) (Marshaler, error) {
	k := t.Kind()
	if k != reflect.Array && k != reflect.Slice {
		return nil, &wrongKindError{Expected: reflect.Array, Actual: t}
	}

	em, err := opts.MarshalerFactory.Marshaler(t.Elem(), opts)
	if err != nil {
		return nil, err
	}
	return &arrayAndSliceMarshaler{
		Type:          t,
		ElemMarshaler: em,
	}, nil
}

func (p *arrayAndSliceMarshaler) Marshal(v reflect.Value, opts *MarshalOptions) ([]string, error) {
	t := v.Type()
	if t != p.Type {
		return nil, &wrongTypeError{Actual: t, Expected: p.Type}
	}

	vlen := v.Len()
	if vlen == 0 {
		return nil, nil
	}

	a := make([]string, vlen)
	for i := 0; i < vlen; i++ {
		a2, err := p.ElemMarshaler.Marshal(v.Index(i), opts)
		if err != nil {
			return nil, fmt.Errorf("error marshaling array/slice index %v :: %v", i, err)
		}
		if len(a2) != 1 {
			return nil, fmt.Errorf("marshaler returned a slice of length %v for array/slice index %v", len(a2), i)
		}
		a[i] = a2[0]
	}
	return a, nil
}

func marshalString(v reflect.Value, opts *MarshalOptions) (string, error) {
	if v.Kind() != reflect.String {
		return "", &wrongKindError{Expected: reflect.String, Actual: v.Type()}
	}
	return v.String(), nil
}

func marshalBool(v reflect.Value, opts *MarshalOptions) (string, error) {
	if v.Kind() != reflect.Bool {
		return "", &wrongKindError{Expected: reflect.Bool, Actual: v.Type()}
	}
	return strconv.FormatBool(v.Bool()), nil
}

func marshalInt(v reflect.Value, opts *MarshalOptions) (string, error) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	default:
		return "", &wrongKindError{Expected: reflect.Int, Actual: v.Type()}
	}
}

func marshalUint(v reflect.Value, opts *MarshalOptions) (string, error) {
	switch v.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), nil
	default:
		return "", &wrongKindError{Expected: reflect.Uint, Actual: v.Type()}
	}
}

func marshalFloat(v reflect.Value, opts *MarshalOptions) (string, error) {
	var bitSize int

	switch v.Kind() {
	case reflect.Float32:
		bitSize = 32
	case reflect.Float64:
		bitSize = 64
	default:
		return "", &wrongKindError{Expected: reflect.Float32, Actual: v.Type()}
	}

	return strconv.FormatFloat(v.Float(), 'f', -1, bitSize), nil
}

func marshalTime(v reflect.Value, opts *MarshalOptions) (string, error) {
	t := v.Type()
	if t != timeType {
		return "", &wrongTypeError{Actual: t, Expected: timeType}
	}
	return v.Interface().(time.Time).Format(time.RFC3339), nil
}

func marshalURL(v reflect.Value, opts *MarshalOptions) (string, error) {
	t := v.Type()
	if t != urlType {
		return "", &wrongTypeError{Actual: t, Expected: urlType}
	}
	u := v.Interface().(url.URL)
	return u.String(), nil
}

func marshalWithMarshalQS(v reflect.Value, opts *MarshalOptions) ([]string, error) {
	marshalQS, ok := v.Interface().(MarshalQS)
	if !ok {
		return nil, fmt.Errorf("expected a type that implements MarshalQS, got %v", v.Type())
	}
	return marshalQS.MarshalQS(opts)
}
