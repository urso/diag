package diag

import (
	"math"
	"time"
)

// Field to be stored in a context.
type Field struct {
	Key   string
	Value Value

	// Standardized indicates that the field its key and value are standardized
	// according to some external schema. Consumers of a context might decide to
	// handle Standardized and non-standardized fields differently.
	Standardized bool
}

// Value represents a reportable value to be stored in a Field.
// The Value struct provides a slot for primitive values that require only
// 64bits, a string, or an arbitrary interface. The interpretation of the slots is up to the Reporter.
type Value struct {
	Primitive uint64
	String    string
	Ifc       interface{}

	Reporter Reporter
}

// Reporter defines the type and supports unpacking, querying the decoded Value.
type Reporter interface {
	Type() Type

	// Ifc decodes the Value and reports the decoded value to as `interface{}`
	// to the provided callback.
	Ifc(*Value, func(interface{}))
}

// Type represents the possible types a Value can have.
type Type uint8

const (
	IfcType Type = iota
	IntType
	Int64Type
	Uint64Type
	Float64Type
	DurationType
	TimestampType
	StringType
)

// Interface decodes and returns the value stored in Value.
func (v *Value) Interface() (ifc interface{}) {
	v.Reporter.Ifc(v, func(tmp interface{}) {
		ifc = tmp
	})
	return ifc
}

func userField(k string, v Value) Field {
	return Field{Key: k, Value: v}
}

// Int creates a new user-field storing an int.
func Int(key string, i int) Field { return userField(key, ValInt(i)) }

// ValInt create a new Value representing an int.
func ValInt(i int) Value { return Value{Primitive: uint64(i), Reporter: _intReporter} }

type intReporter struct{}

var _intReporter Reporter = intReporter{}

func (intReporter) Type() Type                         { return IntType }
func (intReporter) Ifc(v *Value, fn func(interface{})) { fn(int(v.Primitive)) }

// Int64 creates a new user-field storing an int64 value.
func Int64(key string, i int64) Field { return userField(key, ValInt64(i)) }

// ValInt64 creates a new Value representing an int64.
func ValInt64(i int64) Value { return Value{Primitive: uint64(i), Reporter: _int64Reporter} }

type int64Reporter struct{}

var _int64Reporter Reporter = int64Reporter{}

func (int64Reporter) Type() Type                           { return Int64Type }
func (int64Reporter) Ifc(v *Value, fn func(v interface{})) { fn(int64(v.Primitive)) }

// Uint creates a new user-field storing an uint.
func Uint(key string, i uint) Field { return userField(key, ValUint(i)) }

// ValUint creates a new Value representing an uint.
func ValUint(i uint) Value { return ValUint64(uint64(i)) }

// Uint64 creates a new user-field storing an uint64.
func Uint64(key string, i uint64) Field { return userField(key, ValUint64(i)) }

// ValUint64 creates a new Value representing an uint64.
func ValUint64(u uint64) Value { return Value{Primitive: u, Reporter: _uint64Reporter} }

type uint64Reporter struct{}

var _uint64Reporter Reporter = uint64Reporter{}

func (uint64Reporter) Type() Type                           { return Int64Type }
func (uint64Reporter) Ifc(v *Value, fn func(v interface{})) { fn(uint64(v.Primitive)) }

// Float creates a new user-field storing a float.
func Float(key string, f float64) Field { return userField(key, ValFloat(f)) }

// ValFloat creates a new Value representing a float.
func ValFloat(f float64) Value {
	return Value{Primitive: math.Float64bits(f), Reporter: _float64Reporter}
}

type float64Reporter struct{}

var _float64Reporter Reporter = float64Reporter{}

func (float64Reporter) Type() Type                           { return Float64Type }
func (float64Reporter) Ifc(v *Value, fn func(v interface{})) { fn(math.Float64frombits(v.Primitive)) }

// String creates a new user-field storing a string.
func String(key, str string) Field { return userField(key, ValString(str)) }

// ValString creates a new Value representing a string.
func ValString(str string) Value { return Value{String: str, Reporter: _strReporter} }

type strReporter struct{}

var _strReporter Reporter = strReporter{}

func (strReporter) Type() Type                           { return StringType }
func (strReporter) Ifc(v *Value, fn func(v interface{})) { fn(v.String) }

// Duration creates a new user-field storing a duration.
func Duration(key string, dur time.Duration) Field { return userField(key, ValDuration(dur)) }

// ValDuration creates a new Value representing a duration.
func ValDuration(dur time.Duration) Value {
	return Value{Primitive: uint64(dur), Reporter: _durReporter}
}

type durReporter struct{}

var _durReporter Reporter = durReporter{}

func (durReporter) Type() Type                           { return DurationType }
func (durReporter) Ifc(v *Value, fn func(v interface{})) { fn(time.Duration(v.Primitive)) }

// Timestamp creates a new user-field storing a time value.
func Timestamp(key string, ts time.Time) Field { return userField(key, ValTime(ts)) }

// ValTime creates a new Value representing a timestamp.
func ValTime(ts time.Time) Value {
	return Value{Ifc: ts, Reporter: _timeReporter}
}

type timeReporter struct{}

var _timeReporter Reporter = timeReporter{}

func (timeReporter) Type() Type { return TimestampType }
func (timeReporter) Ifc(v *Value, fn func(v interface{})) {
	fn(v.Ifc)
}

// Any creates a new user-field storing any value as interface.
func Any(key string, ifc interface{}) Field {
	// TODO: use type switch + reflection to select concrete Field
	return userField(key, ValAny(ifc))
}

// ValAny creates a new Value representing any value as interface.
func ValAny(ifc interface{}) Value { return Value{Ifc: ifc, Reporter: _anyReporter} }
func reportAny(v *Value, fn func(v interface{})) {
	fn(v.Ifc)
}

type anyReporter struct{}

var _anyReporter Reporter = anyReporter{}

func (anyReporter) Type() Type                           { return IfcType }
func (anyReporter) Ifc(v *Value, fn func(v interface{})) { fn(v.Ifc) }
