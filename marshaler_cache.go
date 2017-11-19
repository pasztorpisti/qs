package qs

import "reflect"

func newValuesMarshalerCache(wrapped ValuesMarshalerFactory) ValuesMarshalerFactory {
	return &valuesMarshalerCache{
		wrapped: wrapped,
		cache:   newSyncMap(),
	}
}

type valuesMarshalerCache struct {
	wrapped ValuesMarshalerFactory
	cache   syncMap
}

func (o *valuesMarshalerCache) ValuesMarshaler(t reflect.Type, opts *MarshalOptions) (ValuesMarshaler, error) {
	if item, ok := o.cache.Load(t); ok {
		if m, ok := item.(ValuesMarshaler); ok {
			return m, nil
		}
		return nil, item.(error)
	}

	m, err := o.wrapped.ValuesMarshaler(t, opts)
	if err != nil {
		o.cache.Store(t, err)
	} else {
		o.cache.Store(t, m)
	}
	return m, err
}

func newMarshalerCache(wrapped MarshalerFactory) MarshalerFactory {
	return &marshalerCache{
		wrapped: wrapped,
		cache:   newSyncMap(),
	}
}

type marshalerCache struct {
	wrapped MarshalerFactory
	cache   syncMap
}

func (o *marshalerCache) Marshaler(t reflect.Type, opts *MarshalOptions) (Marshaler, error) {
	if item, ok := o.cache.Load(t); ok {
		if m, ok := item.(Marshaler); ok {
			return m, nil
		}
		return nil, item.(error)
	}

	m, err := o.wrapped.Marshaler(t, opts)
	if err != nil {
		o.cache.Store(t, err)
	} else {
		o.cache.Store(t, m)
	}
	return m, err
}
