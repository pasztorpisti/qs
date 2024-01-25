package qs

import (
	"errors"
	"reflect"
	"testing"
)

type fakeValuesMarshalerFactory struct {
	m     ValuesMarshaler
	err   error
	calls []reflect.Type
}

func (o *fakeValuesMarshalerFactory) ValuesMarshaler(t reflect.Type, opts *MarshalOptions) (ValuesMarshaler, error) {
	o.calls = append(o.calls, t)
	return o.m, o.err
}

func TestValuesMarshalerCacheSuccess(t *testing.T) {
	expected := &structMarshaler{}
	wrapped := &fakeValuesMarshalerFactory{m: expected}
	cache := newValuesMarshalerCache(wrapped)
	tp := reflect.TypeOf((*fakeValuesMarshalerFactory)(nil)).Elem()

	// cache miss
	m, err := cache.ValuesMarshaler(tp, nil)
	if err != nil {
		t.Fatal(err)
	}
	if m != expected {
		t.Fatalf("got %v, want %v", m, expected)
	}
	if len(wrapped.calls) != 1 || wrapped.calls[0] != tp {
		t.Fatalf("got %v, want %v", wrapped.calls, []reflect.Type{tp})
	}

	// cache hit
	m, err = cache.ValuesMarshaler(tp, nil)
	if err != nil {
		t.Fatal(err)
	}
	if m != expected {
		t.Fatalf("got %v, want %v", m, expected)
	}
	if len(wrapped.calls) != 1 || wrapped.calls[0] != tp {
		t.Fatalf("got %v, want %v", wrapped.calls, []reflect.Type{tp})
	}
}

func TestValuesMarshalerCacheError(t *testing.T) {
	e := errors.New("test error")
	wrapped := &fakeValuesMarshalerFactory{err: e}
	cache := newValuesMarshalerCache(wrapped)
	tp := reflect.TypeOf((*fakeValuesMarshalerFactory)(nil)).Elem()

	// cache miss
	_, err := cache.ValuesMarshaler(tp, nil)
	if err != e {
		t.Fatalf("got %q, want %q", err, e)
	}
	if len(wrapped.calls) != 1 || wrapped.calls[0] != tp {
		t.Fatalf("got %v, want %v", wrapped.calls, []reflect.Type{tp})
	}

	// cache hit
	_, err = cache.ValuesMarshaler(tp, nil)
	if err != e {
		t.Fatalf("got %q, want %q", err, e)
	}
	if len(wrapped.calls) != 1 || wrapped.calls[0] != tp {
		t.Fatalf("got %v, want %v", wrapped.calls, []reflect.Type{tp})
	}
}

type fakeMarshalerFactory struct {
	m     Marshaler
	err   error
	calls []reflect.Type
}

func (o *fakeMarshalerFactory) Marshaler(t reflect.Type, opts *MarshalOptions) (Marshaler, error) {
	o.calls = append(o.calls, t)
	return o.m, o.err
}

type fakeMarshaler struct{}

func (o *fakeMarshaler) Marshal(v reflect.Value, opts *MarshalOptions) ([]string, error) {
	return nil, nil
}

func TestMarshalerCacheSuccess(t *testing.T) {
	// we need a comparable fakeMarshaler object to be able to assert
	expected := &fakeMarshaler{}
	wrapped := &fakeMarshalerFactory{m: expected}
	cache := newMarshalerCache(wrapped)
	tp := reflect.TypeOf((*fakeMarshalerFactory)(nil)).Elem()

	// cache miss
	m, err := cache.Marshaler(tp, nil)
	if err != nil {
		t.Fatal(err)
	}
	if m != expected {
		t.Fatalf("got %v, want %v", m, expected)
	}
	if len(wrapped.calls) != 1 || wrapped.calls[0] != tp {
		t.Fatalf("got %v, want %v", wrapped.calls, []reflect.Type{tp})
	}

	// cache hit
	m, err = cache.Marshaler(tp, nil)
	if err != nil {
		t.Fatal(err)
	}
	if m != expected {
		t.Fatalf("got %v, want %v", m, expected)
	}
	if len(wrapped.calls) != 1 || wrapped.calls[0] != tp {
		t.Fatalf("got %v, want %v", wrapped.calls, []reflect.Type{tp})
	}
}

func TestMarshalerCacheError(t *testing.T) {
	e := errors.New("test error")
	wrapped := &fakeMarshalerFactory{err: e}
	cache := newMarshalerCache(wrapped)
	tp := reflect.TypeOf((*fakeMarshalerFactory)(nil)).Elem()

	// cache miss
	_, err := cache.Marshaler(tp, nil)
	if err != e {
		t.Fatalf("got %q, want %q", err, e)
	}
	if len(wrapped.calls) != 1 || wrapped.calls[0] != tp {
		t.Fatalf("got %v, want %v", wrapped.calls, []reflect.Type{tp})
	}

	// cache hit
	_, err = cache.Marshaler(tp, nil)
	if err != e {
		t.Fatalf("got %q, want %q", err, e)
	}
	if len(wrapped.calls) != 1 || wrapped.calls[0] != tp {
		t.Fatalf("got %v, want %v", wrapped.calls, []reflect.Type{tp})
	}
}
