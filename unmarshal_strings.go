package qs

import (
	"encoding"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

type ptrUnmarshaler struct {
	Type            reflect.Type
	ElemType        reflect.Type
	ElemUnmarshaler Unmarshaler
}

func newPtrUnmarshaler(t reflect.Type, opts *UnmarshalOptions) (Unmarshaler, error) {
	if t.Kind() != reflect.Ptr {
		return nil, &wrongKindError{Expected: reflect.Ptr, Actual: t}
	}
	et := t.Elem()
	eu, err := opts.UnmarshalerFactory.Unmarshaler(et, opts)
	if err != nil {
		return nil, err
	}
	return &ptrUnmarshaler{
		Type:            t,
		ElemType:        et,
		ElemUnmarshaler: eu,
	}, nil
}

func (p *ptrUnmarshaler) Unmarshal(v reflect.Value, a []string, opts *UnmarshalOptions) error {
	t := v.Type()
	if t != p.Type {
		return &wrongTypeError{Actual: t, Expected: p.Type}
	}
	if a == nil {
		return nil
	}
	if v.IsNil() {
		v.Set(reflect.New(p.ElemType))
	}
	return p.ElemUnmarshaler.Unmarshal(v.Elem(), a, opts)
}

type arrayUnmarshaler struct {
	Type            reflect.Type
	ElemUnmarshaler Unmarshaler
	Len             int
}

func newArrayUnmarshaler(t reflect.Type, opts *UnmarshalOptions) (Unmarshaler, error) {
	if t.Kind() != reflect.Array {
		return nil, &wrongKindError{Expected: reflect.Array, Actual: t}
	}

	eu, err := opts.UnmarshalerFactory.Unmarshaler(t.Elem(), opts)
	if err != nil {
		return nil, err
	}
	return &arrayUnmarshaler{
		Type:            t,
		ElemUnmarshaler: eu,
		Len:             t.Len(),
	}, nil
}

func (p *arrayUnmarshaler) Unmarshal(v reflect.Value, a []string, opts *UnmarshalOptions) error {
	t := v.Type()
	if t != p.Type {
		return &wrongTypeError{Actual: t, Expected: p.Type}
	}

	if a == nil {
		return nil
	}
	if len(a) != p.Len {
		return fmt.Errorf("array length == %v, want %v", len(a), p.Len)
	}
	for i := range a {
		err := p.ElemUnmarshaler.Unmarshal(v.Index(i), a[i:i+1], opts)
		if err != nil {
			return fmt.Errorf("error unmarshaling array index %v :: %v", i, err)
		}
	}
	return nil
}

type sliceUnmarshaler struct {
	Type            reflect.Type
	ElemUnmarshaler Unmarshaler
}

func newSliceUnmarshaler(t reflect.Type, opts *UnmarshalOptions) (Unmarshaler, error) {
	if t.Kind() != reflect.Slice {
		return nil, &wrongKindError{Expected: reflect.Slice, Actual: t}
	}

	eu, err := opts.UnmarshalerFactory.Unmarshaler(t.Elem(), opts)
	if err != nil {
		return nil, err
	}
	return &sliceUnmarshaler{
		Type:            t,
		ElemUnmarshaler: eu,
	}, nil
}

func (p *sliceUnmarshaler) Unmarshal(v reflect.Value, a []string, opts *UnmarshalOptions) error {
	t := v.Type()
	if t != p.Type {
		return &wrongTypeError{Actual: t, Expected: p.Type}
	}

	if v.IsNil() {
		v.Set(reflect.MakeSlice(t, len(a), len(a)))
	}

	for i := range a {
		err := p.ElemUnmarshaler.Unmarshal(v.Index(i), a[i:i+1], opts)
		if err != nil {
			return fmt.Errorf("error unmarshaling slice index %v :: %v", i, err)
		}
	}

	return nil
}

// unmarshalString can unmarshal an ini file entry into a value with an
// underlying type (kind) of string.
func unmarshalString(v reflect.Value, s string, opts *UnmarshalOptions) error {
	if v.Kind() != reflect.String {
		return &wrongKindError{Expected: reflect.String, Actual: v.Type()}
	}
	v.SetString(s)
	return nil
}

// unmarshalBool can unmarshal an ini file entry into a value with an
// underlying type (kind) of bool.
func unmarshalBool(v reflect.Value, s string, opts *UnmarshalOptions) error {
	if v.Kind() != reflect.Bool {
		return &wrongKindError{Expected: reflect.Bool, Actual: v.Type()}
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	v.SetBool(b)
	return nil
}

// unmarshalInt can unmarshal an ini file entry into a signed integer value
// with an underlying type (kind) of int, int8, int16, int32 or int64.
func unmarshalInt(v reflect.Value, s string, opts *UnmarshalOptions) error {
	var bitSize int

	switch v.Kind() {
	case reflect.Int:
	case reflect.Int8:
		bitSize = 8
	case reflect.Int16:
		bitSize = 16
	case reflect.Int32:
		bitSize = 32
	case reflect.Int64:
		bitSize = 64
	default:
		return &wrongKindError{Expected: reflect.Int, Actual: v.Type()}
	}

	i, err := strconv.ParseInt(s, 0, bitSize)
	if err != nil {
		return err
	}

	v.SetInt(i)
	return nil
}

// unmarshalUint can unmarshal an ini file entry into an unsigned integer value
// with an underlying type (kind) of uint, uint8, uint16, uint32 or uint64.
func unmarshalUint(v reflect.Value, s string, opts *UnmarshalOptions) error {
	var bitSize int

	switch v.Kind() {
	case reflect.Uint:
	case reflect.Uint8:
		bitSize = 8
	case reflect.Uint16:
		bitSize = 16
	case reflect.Uint32:
		bitSize = 32
	case reflect.Uint64:
		bitSize = 64
	default:
		return &wrongKindError{Expected: reflect.Uint, Actual: v.Type()}
	}

	i, err := strconv.ParseUint(s, 0, bitSize)
	if err != nil {
		return err
	}

	v.SetUint(i)
	return nil
}

func unmarshalFloat(v reflect.Value, s string, opts *UnmarshalOptions) error {
	var bitSize int

	switch v.Kind() {
	case reflect.Float32:
		bitSize = 32
	case reflect.Float64:
		bitSize = 64
	default:
		return &wrongKindError{Expected: reflect.Float32, Actual: v.Type()}
	}

	f, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return err
	}

	v.SetFloat(f)
	return nil
}

func unmarshalTime(v reflect.Value, s string, opts *UnmarshalOptions) error {
	t := v.Type()
	if t != timeType {
		return &wrongTypeError{Actual: t, Expected: timeType}
	}

	tm, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(tm))
	return nil
}

func unmarshalURL(v reflect.Value, s string, opts *UnmarshalOptions) error {
	t := v.Type()
	if t != urlType {
		return &wrongTypeError{Actual: t, Expected: urlType}
	}

	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(*u))
	return nil
}

func unmarshalWithUnmarshalQS(v reflect.Value, a []string, opts *UnmarshalOptions) error {
	if !v.CanAddr() {
		return fmt.Errorf("expected and addressable value, got %v", v)
	}
	unmarshalQS, ok := v.Addr().Interface().(UnmarshalQS)
	if !ok {
		return fmt.Errorf("expected a type that implements UnmarshalQS, got %v", v.Type())
	}
	return unmarshalQS.UnmarshalQS(a, opts)
}

func unmarshalWithTextUnmarshaler(v reflect.Value, a []string, opts *UnmarshalOptions) error {
	if !v.CanAddr() {
		return fmt.Errorf("expected and addressable value, got %v", v)
	}
	unmarshaler, ok := v.Addr().Interface().(encoding.TextUnmarshaler)
	if !ok {
		return fmt.Errorf("expected a type that implements encoding.TextUnmarshaler, got %v", v.Type())
	}
	text, err := opts.SliceToString(a)
	if err != nil {
		return err
	}
	return unmarshaler.UnmarshalText([]byte(text))
}
