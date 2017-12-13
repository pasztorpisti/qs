package qs

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
)

// UnmarshalPresence is an enum that controls the unmarshaling of fields.
// This option is used by the unmarshaler only if the given field isn't present
// in the query string or url.Values that is being unmarshaled.
type UnmarshalPresence int

const (
	// UPUnspecified is the zero value of UnmarshalPresence. In most cases
	// you will use this implicitly by simply leaving the
	// UnmarshalOptions.DefaultUnmarshalPresence field uninitialised which results
	// in using the default UnmarshalPresence which is Opt.
	UPUnspecified UnmarshalPresence = iota

	// Opt tells the unmarshaler to leave struct fields as they are when they
	// aren't present in the query string. However, nil pointers and arrays are
	// created and initialised with new objects.
	Opt

	// Nil is the same as Opt except that it doesn't initialise nil pointers
	// and slices during unmarshal when they are missing from the query string.
	Nil

	// Req tells the unmarshaler to fail with an error that can be detected
	// using qs.IsRequiredFieldError if the given field is
	// missing from the query string. While this is rather validation than
	// unmarshaling it is practical to have this in case of simple programs.
	// If you don't want to mix unmarshaling and validation then you can use the
	// Nil option instead with nil pointers and nil arrays to be able to detect
	// missing fields after unmarshaling.
	Req
)

func (v UnmarshalPresence) String() string {
	switch v {
	case UPUnspecified:
		return "UPUnspecified"
	case Opt:
		return "Opt"
	case Nil:
		return "Nil"
	case Req:
		return "Req"
	default:
		return fmt.Sprintf("UnmarshalPresence(%v)", int(v))
	}
}

// UnmarshalOptions is used as a parameter by the NewUnmarshaler function.
type UnmarshalOptions struct {
	// NameTransformer is used to transform struct field names into a query
	// string names when they aren't set explicitly in the struct field tag.
	// If this field is nil then NewUnmarshaler uses a default function that
	// converts the CamelCase field names to snake_case which is popular
	// with query strings.
	NameTransformer NameTransformFunc

	// SliceToString is used by Unmarshaler.Unmarshal when it unmarshals into a
	// primitive non-array struct field. In such cases unmarshaling a []string
	// (which is the value type of the url.Values map) requires transforming
	// the []string into a single string before unmarshaling.
	//
	// E.g.: If you have a struct field "Count int" but you receive a query
	// string "count=5&count=6&count=8" then the incoming []string{"5", "6", "8"}
	// has to be converted into a single string before setting the "Count int"
	// field.
	//
	// If you don't initialise this field then a default function is used that
	// fails if the input array doesn't contain exactly one item.
	//
	// In some cases you might want to provide your own function that is more
	// forgiving. E.g.: you can provide a function that picks the first or last
	// item, or concatenates/joins the whole list into a single string.
	SliceToString func([]string) (string, error)

	// ValuesUnmarshalerFactory is used by QSUnmarshaler to create ValuesUnmarshaler
	// objects for specific types. If this field is nil then NewUnmarshaler uses
	// a default builtin factory.
	ValuesUnmarshalerFactory ValuesUnmarshalerFactory

	// UnmarshalerFactory is used by QSUnmarshaler to create Unmarshaler
	// objects for specific types. If this field is nil then NewUnmarshaler uses
	// a default builtin factory.
	UnmarshalerFactory UnmarshalerFactory

	// DefaultUnmarshalPresence is used for the unmarshaling of struct fields
	// that don't have an explicit UnmarshalPresence option set in their tags.
	DefaultUnmarshalPresence UnmarshalPresence
}

// DefaultUnmarshaler is the unmarshaler used by the Unmarshal, UnmarshalValues,
// CanUnmarshal and CanUnmarshalType functions.
var DefaultUnmarshaler = NewUnmarshaler(&UnmarshalOptions{})

// Unmarshal unmarshals a query string and stores the result to the object
// pointed to by the given pointer.
//
// Unmarshal uses the inverse of the encodings that Marshal uses.
//
// A struct field tag can optionally contain one of the opt, nil and req options
// for unmarshaling. If it contains none of these then opt is the default but
// the default can also be changed by using a custom marshaler. The
// UnmarshalPresence of a field is used only when the query string doesn't
// contain a value for it:
//  - nil succeeds and keeps the original field value
//  - opt succeeds and keeps the original field value but in case of
//    pointer-like types (pointers, slices) with nil field value it initialises
//    the field with a newly created object.
//  - req causes the unmarshal operation to fail with an error that can be
//    detected using qs.IsRequiredFieldError.
//
// When unmarshaling a nil pointer field that is present in the query string
// the pointer is automatically initialised even if it has the nil option in
// its tag.
func Unmarshal(into interface{}, queryString string) error {
	return DefaultUnmarshaler.Unmarshal(into, queryString)
}

// UnmarshalValues is the same as Unmarshal but it unmarshals from a url.Values
// instead of a query string.
func UnmarshalValues(into interface{}, values url.Values) error {
	return DefaultUnmarshaler.UnmarshalValues(into, values)
}

// CheckUnmarshal returns an error if the type of the given object can't be
// unmarshaled from a url.Vales or query string. By default only maps and structs
// can be unmarshaled from query strings given that all of their fields or values
// can be unmarshaled from []string (which is the value type of the url.Values map).
//
// It performs the check on the type of the object without traversing or
// unmarshaling the object.
func CheckUnmarshal(into interface{}) error {
	return DefaultUnmarshaler.CheckUnmarshal(into)
}

// CheckUnmarshalType returns an error if the given type can't be unmarshaled
// from a url.Vales or query string. By default only maps and structs
// can be unmarshaled from query strings given that all of their fields or values
// can be unmarshaled from []string (which is the value type of the url.Values map).
func CheckUnmarshalType(t reflect.Type) error {
	return DefaultUnmarshaler.CheckUnmarshalType(t)
}

// QSUnmarshaler objects can be created by calling NewUnmarshaler and they can be
// used to unmarshal query strings or url.Values into structs or maps.
type QSUnmarshaler struct {
	opts *UnmarshalOptions
}

// NewUnmarshaler returns a new QSUnmarshaler object.
func NewUnmarshaler(opts *UnmarshalOptions) *QSUnmarshaler {
	return &QSUnmarshaler{
		opts: prepareUnmarshalOptions(*opts),
	}
}

// Unmarshal unmarshals an object from a query string.
// See the documentation of the global Unmarshal func.
func (p *QSUnmarshaler) Unmarshal(into interface{}, queryString string) error {
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return fmt.Errorf("error parsing query string %q :: %v", queryString, err)
	}
	return p.UnmarshalValues(into, values)
}

// UnmarshalValues unmarshals an object from a url.Values.
// See the documentation of the global UnmarshalValues func.
func (p *QSUnmarshaler) UnmarshalValues(into interface{}, values url.Values) error {
	pv := reflect.ValueOf(into)
	if !pv.IsValid() {
		return errors.New("received an empty interface")
	}
	if pv.Kind() != reflect.Ptr {
		return fmt.Errorf("expected a pointer, got %T", into)
	}
	if pv.IsNil() {
		return fmt.Errorf("nil pointer of type %T", into)
	}
	v := pv.Elem()

	vum, err := p.opts.ValuesUnmarshalerFactory.ValuesUnmarshaler(v.Type(), p.opts)
	if err != nil {
		return err
	}
	return vum.UnmarshalValues(v, values, p.opts)
}

// CheckUnmarshal check whether the type of the given object supports
// unmarshaling from query strings.
// See the documentation of the global CheckUnmarshal func.
func (p *QSUnmarshaler) CheckUnmarshal(into interface{}) error {
	return p.CheckUnmarshalType(reflect.TypeOf(into))
}

// CheckUnmarshalType check whether the given type supports unmarshaling from
// query strings. See the documentation of the global CheckUnmarshalType func.
func (p *QSUnmarshaler) CheckUnmarshalType(t reflect.Type) error {
	if t == nil {
		return errors.New("nil type")
	}
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("expected a pointer, got %v", t)
	}
	_, err := p.opts.ValuesUnmarshalerFactory.ValuesUnmarshaler(t.Elem(), p.opts)
	return err
}

// NewDefaultUnmarshalOptions creates a new UnmarshalOptions in which every field
// is set to its default value.
func NewDefaultUnmarshalOptions() *UnmarshalOptions {
	return prepareUnmarshalOptions(UnmarshalOptions{})
}

// defaultSliceToString is used by the NewUnmarshaler function when
// its UnmarshalOptions.SliceToString parameter is nil.
var defaultSliceToString = func(a []string) (string, error) {
	if len(a) != 1 {
		return "", fmt.Errorf("SliceToString expects array length == 1. array=%v", a)
	}
	return a[0], nil
}

// defaultValuesUnmarshalerFactory is used by the NewUnmarshaler function when
// its UnmarshalOptions.ValuesUnmarshalerFactory parameter is nil.
var defaultValuesUnmarshalerFactory = newValuesUnmarshalerFactory()

// defaultUnmarshalerFactory is used by the NewUnmarshaler function when its
// UnmarshalOptions.UnmarshalerFactory parameter is nil. This variable is set
// to a factory object that handles most builtin types (arrays, pointers,
// bool, int, etc...). If a type implements the UnmarshalQS interface then this
// factory returns an unmarshaler object that allows instances of the given type
// to unmarshal themselves.
var defaultUnmarshalerFactory = newUnmarshalerFactory()

// defaultUnmarshalPresence is used by the NewUnmarshaler function when its
// UnmarshalOptions.DefaultUnmarshalPresence parameter is UPUnspecified.
const defaultUnmarshalPresence = Opt

func prepareUnmarshalOptions(opts UnmarshalOptions) *UnmarshalOptions {
	if opts.NameTransformer == nil {
		opts.NameTransformer = snakeCase
	}
	if opts.SliceToString == nil {
		opts.SliceToString = defaultSliceToString
	}

	if opts.ValuesUnmarshalerFactory == nil {
		opts.ValuesUnmarshalerFactory = defaultValuesUnmarshalerFactory
	}
	opts.ValuesUnmarshalerFactory = newValuesUnmarshalerCache(opts.ValuesUnmarshalerFactory)

	if opts.UnmarshalerFactory == nil {
		opts.UnmarshalerFactory = defaultUnmarshalerFactory
	}
	opts.UnmarshalerFactory = newUnmarshalerCache(opts.UnmarshalerFactory)

	if opts.DefaultUnmarshalPresence == UPUnspecified {
		opts.DefaultUnmarshalPresence = defaultUnmarshalPresence
	}
	return &opts
}
