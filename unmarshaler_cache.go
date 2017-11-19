package qs

import "reflect"

func newValuesUnmarshalerCache(wrapped ValuesUnmarshalerFactory) ValuesUnmarshalerFactory {
	return &valuesUnmarshalerCache{
		wrapped: wrapped,
		cache:   newSyncMap(),
	}
}

type valuesUnmarshalerCache struct {
	wrapped ValuesUnmarshalerFactory
	cache   syncMap
}

func (o *valuesUnmarshalerCache) ValuesUnmarshaler(t reflect.Type, opts *UnmarshalOptions) (ValuesUnmarshaler, error) {
	if item, ok := o.cache.Load(t); ok {
		if m, ok := item.(ValuesUnmarshaler); ok {
			return m, nil
		}
		return nil, item.(error)
	}

	u, err := o.wrapped.ValuesUnmarshaler(t, opts)
	if err != nil {
		o.cache.Store(t, err)
	} else {
		o.cache.Store(t, u)
	}
	return u, err
}

func newUnmarshalerCache(wrapped UnmarshalerFactory) UnmarshalerFactory {
	return &unmarshalerCache{
		wrapped: wrapped,
		cache:   newSyncMap(),
	}
}

type unmarshalerCache struct {
	wrapped UnmarshalerFactory
	cache   syncMap
}

func (o *unmarshalerCache) Unmarshaler(t reflect.Type, opts *UnmarshalOptions) (Unmarshaler, error) {
	if item, ok := o.cache.Load(t); ok {
		if m, ok := item.(Unmarshaler); ok {
			return m, nil
		}
		return nil, item.(error)
	}

	u, err := o.wrapped.Unmarshaler(t, opts)
	if err != nil {
		o.cache.Store(t, err)
	} else {
		o.cache.Store(t, u)
	}
	return u, err
}
