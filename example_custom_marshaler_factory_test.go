package qs_test

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"time"

	"github.com/pamburus/qs"
)

// This example shows how to create QSMarshaler and QSUnmarshaler objects
// that have custom marshaler and unmarshaler factories to provide custom
// marshaling and unmarshaling for some types.
//
// In this example we change the default marshaling and unmarshaling of the
// []byte type and we compare our custom marshaler with the default one. You can
// not only change the behavior of already supported types (like []byte) but you
// can also add types that aren't supported by default - in this example we
// add time.Duration as one such type.
//
// Builtin unnamed golang types (like []byte) can't implement the MarshalQS and
// UnmarshalQS interfaces to provide their own marshaling, this is why we have
// to create custom QSMarshaler and QSUnmarshaler with custom factories for them.
func Example_customMarshalerFactory() {
	customMarshaler := qs.NewMarshaler(&qs.MarshalOptions{
		MarshalerFactory: &marshalerFactory{qs.NewDefaultMarshalOptions().MarshalerFactory},
	})
	customUnmarshaler := qs.NewUnmarshaler(&qs.UnmarshalOptions{
		UnmarshalerFactory: &unmarshalerFactory{qs.NewDefaultUnmarshalOptions().UnmarshalerFactory},
	})

	performSliceTest("Default", qs.DefaultMarshaler, qs.DefaultUnmarshaler)
	performSliceTest("Custom", customMarshaler, customUnmarshaler)
	performDurationTest(customMarshaler, customUnmarshaler)

	// Output:
	// Default-Marshal-Result: a=0&a=1&a=2&b=3&b=4&b=5 <nil>
	// Default-Unmarshal-Result: len=2 a=[0 1 2] b=[3 4 5] <nil>
	// Custom-Marshal-Result: a=000102&b=030405 <nil>
	// Custom-Unmarshal-Result: len=2 a=[0 1 2] b=[3 4 5] <nil>
	// Duration-Marshal-Result: duration=1m1.2s <nil>
	// Duration-Unmarshal-Result: len=1 duration=1m1.2s <nil>
}

func performSliceTest(name string, m *qs.QSMarshaler, um *qs.QSUnmarshaler) {
	queryStr, err := m.Marshal(map[string][]byte{
		"a": {0, 1, 2},
		"b": {3, 4, 5},
	})
	fmt.Printf("%v-Marshal-Result: %v %v\n", name, queryStr, err)

	var query map[string][]byte
	err = um.Unmarshal(&query, queryStr)
	fmt.Printf("%v-Unmarshal-Result: len=%v a=%v b=%v %v\n",
		name, len(query), query["a"], query["b"], err)
}

func performDurationTest(m *qs.QSMarshaler, um *qs.QSUnmarshaler) {
	queryStr, err := m.Marshal(map[string]time.Duration{
		"duration": time.Millisecond * (61*1000 + 200),
	})
	fmt.Printf("Duration-Marshal-Result: %v %v\n", queryStr, err)

	var query map[string]time.Duration
	err = um.Unmarshal(&query, queryStr)
	fmt.Printf("Duration-Unmarshal-Result: len=%v duration=%v %v\n",
		len(query), query["duration"].String(), err)
}

var byteSliceType = reflect.TypeOf([]byte(nil))
var durationType = reflect.TypeOf(time.Duration(0))

// marshalerFactory implements the MarshalerFactory interface and provides
// custom Marshaler for the []byte type.
type marshalerFactory struct {
	orig qs.MarshalerFactory
}

func (f *marshalerFactory) Marshaler(t reflect.Type, opts *qs.MarshalOptions) (qs.Marshaler, error) {
	switch t {
	case byteSliceType:
		return byteSliceMarshaler{}, nil
	case durationType:
		return durationMarshalerInstance, nil
	default:
		return f.orig.Marshaler(t, opts)
	}
}

// unmarshalerFactory implements the UnmarshalerFactory interface and provides
// custom Unmarshaler for the []byte type.
type unmarshalerFactory struct {
	orig qs.UnmarshalerFactory
}

func (f *unmarshalerFactory) Unmarshaler(t reflect.Type, opts *qs.UnmarshalOptions) (qs.Unmarshaler, error) {
	switch t {
	case byteSliceType:
		return byteSliceMarshaler{}, nil
	case durationType:
		return durationMarshalerInstance, nil
	default:
		return f.orig.Unmarshaler(t, opts)
	}
}

// byteSliceMarshaler implements the Marshaler and Unmarshaler interfaces to
// provide custom marshaling and unmarshaling for the []byte type.
type byteSliceMarshaler struct{}

func (byteSliceMarshaler) Marshal(v reflect.Value, opts *qs.MarshalOptions) ([]string, error) {
	return []string{hex.EncodeToString(v.Interface().([]byte))}, nil
}

func (byteSliceMarshaler) Unmarshal(v reflect.Value, a []string, opts *qs.UnmarshalOptions) error {
	s, err := opts.SliceToString(a)
	if err != nil {
		return err
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(b))
	return nil
}

var durationMarshalerInstance = &durationMarshaler{}

// durationMarshaler implements the Marshaler and Unmarshaler interfaces to
// provide custom marshaling and unmarshaling for the time.Duration type.
type durationMarshaler struct{}

func (o *durationMarshaler) Marshal(v reflect.Value, opts *qs.MarshalOptions) ([]string, error) {
	return []string{v.Interface().(time.Duration).String()}, nil
}

func (o *durationMarshaler) Unmarshal(v reflect.Value, a []string, opts *qs.UnmarshalOptions) error {
	s, err := opts.SliceToString(a)
	if err != nil {
		return err
	}
	t, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("unsupported time format: %v", s)
	}
	v.Set(reflect.ValueOf(t))
	return nil
}
