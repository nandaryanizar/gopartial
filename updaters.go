package gopartial

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/guregu/null"
)

// NullStringUpdater update null.String
func NullStringUpdater(fieldValue reflect.Value, v reflect.Value) bool {
	switch fieldValue.Interface().(type) {
	case null.String:
		// if its null value
		if !v.IsValid() {
			newValue := reflect.ValueOf(null.String{NullString: sql.NullString{Valid: false}})
			fieldValue.Set(newValue)
			return true
		}
		// only set if underlying type is string
		if v.Kind() == reflect.String {
			newValue := reflect.ValueOf(null.String{NullString: sql.NullString{Valid: true, String: v.String()}})
			fieldValue.Set(newValue)
			return true
		}
	}

	return false
}

// NullFloatUpdater update null.Float64
func NullFloatUpdater(fieldValue reflect.Value, v reflect.Value) bool {
	switch fieldValue.Interface().(type) {
	case null.Float:
		// if its null value
		if !v.IsValid() {
			newValue := reflect.ValueOf(null.Float{NullFloat64: sql.NullFloat64{Valid: false}})
			fieldValue.Set(newValue)
			return true
		}
		// only set if underlying type is any int/float
		if v.Kind() == reflect.Int ||
			v.Kind() == reflect.Int8 ||
			v.Kind() == reflect.Int16 ||
			v.Kind() == reflect.Int32 ||
			v.Kind() == reflect.Int64 {
			newValue := reflect.ValueOf(null.Float{NullFloat64: sql.NullFloat64{Valid: true, Float64: float64(v.Int())}})
			fieldValue.Set(newValue)
			return true
		} else if v.Kind() == reflect.Float32 ||
			v.Kind() == reflect.Float64 {
			newValue := reflect.ValueOf(null.Float{NullFloat64: sql.NullFloat64{Valid: true, Float64: v.Float()}})
			fieldValue.Set(newValue)
			return true
		}
	}

	return false
}

// NullIntUpdater update null.Int
func NullIntUpdater(fieldValue reflect.Value, v reflect.Value) bool {
	switch fieldValue.Interface().(type) {
	case null.Int:
		// if its null value
		if !v.IsValid() {
			newValue := reflect.ValueOf(null.Int{NullInt64: sql.NullInt64{Valid: false}})
			fieldValue.Set(newValue)
			return true
		}
		// only set if underlying type is any int/float
		if v.Kind() == reflect.Int ||
			v.Kind() == reflect.Int8 ||
			v.Kind() == reflect.Int16 ||
			v.Kind() == reflect.Int32 ||
			v.Kind() == reflect.Int64 {
			newValue := reflect.ValueOf(null.Int{NullInt64: sql.NullInt64{Valid: true, Int64: v.Int()}})
			fieldValue.Set(newValue)
			return true
		} else if v.Kind() == reflect.Float32 ||
			v.Kind() == reflect.Float64 {
			newValue := reflect.ValueOf(null.Int{NullInt64: sql.NullInt64{Valid: true, Int64: int64(v.Float())}})
			fieldValue.Set(newValue)
			return true
		}
	}

	return false
}

// NullBoolUpdater update null.Bool
func NullBoolUpdater(fieldValue reflect.Value, v reflect.Value) bool {
	switch fieldValue.Interface().(type) {
	case null.Bool:
		// if its null value
		if !v.IsValid() {
			newValue := reflect.ValueOf(null.Bool{NullBool: sql.NullBool{Valid: false}})
			fieldValue.Set(newValue)
			return true
		}
		// only set if underlying type is bool
		if v.Kind() == reflect.Bool {
			newValue := reflect.ValueOf(null.Bool{NullBool: sql.NullBool{Valid: true, Bool: v.Bool()}})
			fieldValue.Set(newValue)
			return true
		}
	}

	return false
}

// NullTimeUpdater update null.Time
func NullTimeUpdater(fieldValue reflect.Value, v reflect.Value) bool {
	switch fieldValue.Interface().(type) {
	case null.Time:
		// if its null value
		if !v.IsValid() {
			newValue := reflect.ValueOf(null.Time{Valid: false})
			fieldValue.Set(newValue)
			return true
		}
		// only set if underlying type is string
		if v.Kind() == reflect.String {
			nullTime := null.Time{}
			if err := nullTime.UnmarshalJSON([]byte(`"` + v.String() + `"`)); err == nil {
				newValue := reflect.ValueOf(nullTime)
				fieldValue.Set(newValue)
				return true
			}
		}
	}

	return false
}

// MapStringInterfaceUpdater update map[string]interface{}
func MapStringInterfaceUpdater(fieldValue reflect.Value, v reflect.Value) bool {
	// if fieldValue.Kind() == reflect.Struct {
	// 	// @TODO: right now, do nothing
	// 	return false
	// }
	return false
}

// BoolUpdater update bool (pointer or value)
func BoolUpdater(fieldValue reflect.Value, v reflect.Value) bool {
	if fieldValue.Kind() == reflect.Bool {
		if v.Kind() == reflect.Bool {
			fieldValue.Set(v)
			return true
		}
	} else if fieldValue.Kind() == reflect.Ptr {
		// only process if field is pointer to any bool
		if fieldValue.Type().String() == "*bool" {
			if !v.IsValid() {
				var newBoolValue *bool
				newValue := reflect.ValueOf(newBoolValue)
				fieldValue.Set(newValue)
				return true
			} else if v.Kind() == reflect.Bool {
				newBoolValue := v.Bool()
				newValue := reflect.ValueOf(&newBoolValue)
				fieldValue.Set(newValue)
				return true
			}
		}
	}

	return false
}

// IntUpdater update int (any int type Int8, Int16, Int32, Int64 and whether its a pointer or a value)
func IntUpdater(fieldValue reflect.Value, v reflect.Value) bool {
	switch fieldValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if fieldValue.OverflowInt(v.Int()) {
				return false
			}
			fieldValue.SetInt(v.Int())
			return true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if fieldValue.OverflowInt(int64(v.Uint())) {
				return false
			}
			fieldValue.SetInt(int64(v.Uint()))
			return true
		case reflect.Float32, reflect.Float64:
			if fieldValue.OverflowInt(int64(v.Float())) {
				return false
			}
			fieldValue.SetInt(int64(v.Float()))
			return true
		}

	case reflect.Ptr:

		if !v.IsValid() {
			if fieldValue.Type().String() == "*int" {
				var newNullInt *int
				newValue := reflect.ValueOf(newNullInt)
				fieldValue.Set(newValue)
				return true
			} else if fieldValue.Type().String() == "*int8" {
				var newNullInt8 *int8
				newValue := reflect.ValueOf(newNullInt8)
				fieldValue.Set(newValue)
				return true
			} else if fieldValue.Type().String() == "*int16" {
				var newNullInt16 *int16
				newValue := reflect.ValueOf(newNullInt16)
				fieldValue.Set(newValue)
				return true
			} else if fieldValue.Type().String() == "*int32" {
				var newNullInt32 *int32
				newValue := reflect.ValueOf(newNullInt32)
				fieldValue.Set(newValue)
				return true
			} else if fieldValue.Type().String() == "*int64" {
				var newNullInt64 *int64
				newValue := reflect.ValueOf(newNullInt64)
				fieldValue.Set(newValue)
				return true
			}
		} else {

			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				vv := v.Int()

				if fieldValue.Type().String() == "*int" {
					var newIntValue int
					if isOverflowInt(newIntValue, vv) {
						return false
					}
					newIntValue = int(vv)
					newValue := reflect.ValueOf(&newIntValue)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*int8" {
					var newInt8Value int8
					if isOverflowInt(newInt8Value, vv) {
						return false
					}
					newInt8Value = int8(vv)
					newValue := reflect.ValueOf(&newInt8Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*int16" {
					var newInt16Value int16
					if isOverflowInt(newInt16Value, vv) {
						return false
					}
					newInt16Value = int16(vv)
					newValue := reflect.ValueOf(&newInt16Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*int32" {
					var newInt32Value int32
					if isOverflowInt(newInt32Value, vv) {
						return false
					}
					newInt32Value = int32(vv)
					newValue := reflect.ValueOf(&newInt32Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*int64" {
					var newInt64Value int64
					if isOverflowInt(newInt64Value, vv) {
						return false
					}
					newInt64Value = vv
					newValue := reflect.ValueOf(&newInt64Value)
					fieldValue.Set(newValue)
					return true
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				vv := v.Uint()

				if fieldValue.Type().String() == "*int" {
					var newIntValue int
					if isOverflowInt(newIntValue, int64(vv)) {
						return false
					}
					newIntValue = int(vv)
					newValue := reflect.ValueOf(&newIntValue)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*int8" {
					var newInt8Value int8
					if isOverflowInt(newInt8Value, int64(vv)) {
						return false
					}
					newInt8Value = int8(vv)
					newValue := reflect.ValueOf(&newInt8Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*int16" {
					var newInt16Value int16
					if isOverflowInt(newInt16Value, int64(vv)) {
						return false
					}
					newInt16Value = int16(vv)
					newValue := reflect.ValueOf(&newInt16Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*int32" {
					var newInt32Value int32
					if isOverflowInt(newInt32Value, int64(vv)) {
						return false
					}
					newInt32Value = int32(vv)
					newValue := reflect.ValueOf(&newInt32Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*int64" {
					var newInt64Value int64
					if isOverflowInt(newInt64Value, int64(vv)) {
						return false
					}
					newInt64Value = int64(vv)
					newValue := reflect.ValueOf(&newInt64Value)
					fieldValue.Set(newValue)
					return true
				}
			case reflect.Float32, reflect.Float64:
				vv := v.Float()

				if fieldValue.Type().String() == "*int" {
					var newIntValue int
					if isOverflowInt(newIntValue, int64(vv)) {
						return false
					}
					newIntValue = int(vv)
					newValue := reflect.ValueOf(&newIntValue)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*int8" {
					var newInt8Value int8
					if isOverflowInt(newInt8Value, int64(vv)) {
						return false
					}
					newInt8Value = int8(vv)
					newValue := reflect.ValueOf(&newInt8Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*int16" {
					var newInt16Value int16
					if isOverflowInt(newInt16Value, int64(vv)) {
						return false
					}
					newInt16Value = int16(vv)
					newValue := reflect.ValueOf(&newInt16Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*int32" {
					var newInt32Value int32
					if isOverflowInt(newInt32Value, int64(vv)) {
						return false
					}
					newInt32Value = int32(vv)
					newValue := reflect.ValueOf(&newInt32Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*int64" {
					var newInt64Value int64
					if isOverflowInt(newInt64Value, int64(vv)) {
						return false
					}
					newInt64Value = int64(vv)
					newValue := reflect.ValueOf(&newInt64Value)
					fieldValue.Set(newValue)
					return true
				}
			}

		}

	}

	return false
}

// UintUpdater update int (any int type Uint8, Uint16, Uint32, Uint64 and whether its a pointer or a value)
func UintUpdater(fieldValue reflect.Value, v reflect.Value) bool {
	// check for underflow first
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() < 0 {
			return false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() < 0 {
			return false
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() < 0 {
			return false
		}
	}

	switch fieldValue.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if fieldValue.OverflowUint(uint64(v.Int())) {
				return false
			}
			fieldValue.SetUint(uint64(v.Int()))
			return true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if fieldValue.OverflowUint(v.Uint()) {
				return false
			}
			fieldValue.SetUint(v.Uint())
			return true
		case reflect.Float32, reflect.Float64:
			if fieldValue.OverflowUint(uint64(v.Float())) {
				return false
			}
			fieldValue.SetUint(uint64(v.Float()))
			return true
		}

	case reflect.Ptr:

		if !v.IsValid() {
			if fieldValue.Type().String() == "*uint" {
				var newNullUint *uint
				newValue := reflect.ValueOf(newNullUint)
				fieldValue.Set(newValue)
				return true
			} else if fieldValue.Type().String() == "*uint8" {
				var newNullUint8 *uint8
				newValue := reflect.ValueOf(newNullUint8)
				fieldValue.Set(newValue)
				return true
			} else if fieldValue.Type().String() == "*uint16" {
				var newNullUint16 *uint16
				newValue := reflect.ValueOf(newNullUint16)
				fieldValue.Set(newValue)
				return true
			} else if fieldValue.Type().String() == "*uint32" {
				var newNullUint32 *uint32
				newValue := reflect.ValueOf(newNullUint32)
				fieldValue.Set(newValue)
				return true
			} else if fieldValue.Type().String() == "*uint64" {
				var newNullUint64 *uint64
				newValue := reflect.ValueOf(newNullUint64)
				fieldValue.Set(newValue)
				return true
			}
		} else {

			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				vv := v.Int()

				if fieldValue.Type().String() == "*uint" {
					var newUintValue uint
					if isOverflowUint(newUintValue, uint64(vv)) {
						return false
					}
					newUintValue = uint(vv)
					newValue := reflect.ValueOf(&newUintValue)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*uint8" {
					var newUint8Value uint8
					if isOverflowUint(newUint8Value, uint64(vv)) {
						return false
					}
					newUint8Value = uint8(vv)
					newValue := reflect.ValueOf(&newUint8Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*uint16" {
					var newUint16Value uint16
					if isOverflowUint(newUint16Value, uint64(vv)) {
						return false
					}
					newUint16Value = uint16(vv)
					newValue := reflect.ValueOf(&newUint16Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*uint32" {
					var newUint32Value uint32
					if isOverflowUint(newUint32Value, uint64(vv)) {
						return false
					}
					newUint32Value = uint32(vv)
					newValue := reflect.ValueOf(&newUint32Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*uint64" {
					var newUint64Value uint64
					if isOverflowUint(newUint64Value, uint64(vv)) {
						return false
					}
					newUint64Value = uint64(vv)
					newValue := reflect.ValueOf(&newUint64Value)
					fieldValue.Set(newValue)
					return true
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				vv := v.Uint()

				if fieldValue.Type().String() == "*uint" {
					var newUintValue uint
					if isOverflowUint(newUintValue, vv) {
						return false
					}
					newUintValue = uint(vv)
					newValue := reflect.ValueOf(&newUintValue)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*uint8" {
					var newUint8Value uint8
					if isOverflowUint(newUint8Value, vv) {
						return false
					}
					newUint8Value = uint8(vv)
					newValue := reflect.ValueOf(&newUint8Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*uint16" {
					var newUint16Value uint16
					if isOverflowUint(newUint16Value, vv) {
						return false
					}
					newUint16Value = uint16(vv)
					newValue := reflect.ValueOf(&newUint16Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*uint32" {
					var newUint32Value uint32
					if isOverflowUint(newUint32Value, vv) {
						return false
					}
					newUint32Value = uint32(vv)
					newValue := reflect.ValueOf(&newUint32Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*uint64" {
					var newUint64Value uint64
					if isOverflowUint(newUint64Value, vv) {
						return false
					}
					newUint64Value = vv
					newValue := reflect.ValueOf(&newUint64Value)
					fieldValue.Set(newValue)
					return true
				}
			case reflect.Float32, reflect.Float64:
				vv := v.Float()

				if fieldValue.Type().String() == "*uint" {
					var newUintValue uint
					if isOverflowUint(newUintValue, uint64(vv)) {
						return false
					}
					newUintValue = uint(vv)
					newValue := reflect.ValueOf(&newUintValue)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*uint8" {
					var newUint8Value uint8
					if isOverflowUint(newUint8Value, uint64(vv)) {
						return false
					}
					newUint8Value = uint8(vv)
					newValue := reflect.ValueOf(&newUint8Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*uint16" {
					var newUint16Value uint16
					if isOverflowUint(newUint16Value, uint64(vv)) {
						return false
					}
					newUint16Value = uint16(vv)
					newValue := reflect.ValueOf(&newUint16Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*uint32" {
					var newUint32Value uint32
					if isOverflowUint(newUint32Value, uint64(vv)) {
						return false
					}
					newUint32Value = uint32(vv)
					newValue := reflect.ValueOf(&newUint32Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*uint64" {
					var newUint64Value uint64
					if isOverflowUint(newUint64Value, uint64(vv)) {
						return false
					}
					newUint64Value = uint64(vv)
					newValue := reflect.ValueOf(&newUint64Value)
					fieldValue.Set(newValue)
					return true
				}
			}

		}

	}

	return false
}

// FloatUpdater update int (any float type Float8, Float16, Float32, Float64 and whether its a pointer or a value)
func FloatUpdater(fieldValue reflect.Value, v reflect.Value) bool {
	switch fieldValue.Kind() {
	case reflect.Float32, reflect.Float64:

		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if fieldValue.OverflowFloat(float64(v.Int())) {
				return false
			}
			fieldValue.SetFloat(float64(v.Int()))
			return true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if fieldValue.OverflowFloat(float64(v.Uint())) {
				return false
			}
			fieldValue.SetFloat(float64(v.Uint()))
			return true
		case reflect.Float32, reflect.Float64:
			if fieldValue.OverflowFloat(v.Float()) {
				return false
			}
			fieldValue.SetFloat(v.Float())
			return true
		}

	case reflect.Ptr:

		// only process if field is pointer to any float
		if fieldValue.Type().String() != "*float32" && fieldValue.Type().String() != "*float64" {
			return false
		}

		if !v.IsValid() {
			if fieldValue.Type().String() == "*float32" {
				var newFloat32Value *float32
				newValue := reflect.ValueOf(newFloat32Value)
				fieldValue.Set(newValue)
				return true
			} else if fieldValue.Type().String() == "*float64" {
				var newFloat64Value *float64
				newValue := reflect.ValueOf(newFloat64Value)
				fieldValue.Set(newValue)
				return true
			}
		} else {

			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				vv := v.Int()

				if fieldValue.Type().String() == "*float32" {
					var newFloat32Value float32
					if isOverflowFloat(newFloat32Value, float64(vv)) {
						return false
					}
					newFloat32Value = float32(vv)
					newValue := reflect.ValueOf(&newFloat32Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*float64" {
					var newFloat64Value float64
					if isOverflowFloat(newFloat64Value, float64(vv)) {
						return false
					}
					newFloat64Value = float64(vv)
					newValue := reflect.ValueOf(&newFloat64Value)
					fieldValue.Set(newValue)
					return true
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				vv := v.Uint()

				if fieldValue.Type().String() == "*float32" {
					var newFloat32Value float32
					if isOverflowFloat(newFloat32Value, float64(vv)) {
						return false
					}
					newFloat32Value = float32(vv)
					newValue := reflect.ValueOf(&newFloat32Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*float64" {
					var newFloat64Value float64
					if isOverflowFloat(newFloat64Value, float64(vv)) {
						return false
					}
					newFloat64Value = float64(vv)
					newValue := reflect.ValueOf(&newFloat64Value)
					fieldValue.Set(newValue)
					return true
				}
			case reflect.Float32, reflect.Float64:
				vv := v.Float()

				if fieldValue.Type().String() == "*float32" {
					var newFloat32Value float32
					if isOverflowFloat(newFloat32Value, vv) {
						return false
					}
					newFloat32Value = float32(vv)
					newValue := reflect.ValueOf(&newFloat32Value)
					fieldValue.Set(newValue)
					return true
				} else if fieldValue.Type().String() == "*float64" {
					var newFloat64Value float64
					if isOverflowFloat(newFloat64Value, vv) {
						return false
					}
					newFloat64Value = vv
					newValue := reflect.ValueOf(&newFloat64Value)
					fieldValue.Set(newValue)
					return true
				}
			}

		}

	}

	return false
}

// TimeUpdater update time (pointer or value)
func TimeUpdater(fieldValue reflect.Value, v reflect.Value) bool {
	switch fieldValue.Interface().(type) {
	case time.Time:
		if !v.IsValid() {
			return false
		}

		switch v.Kind() {
		case reflect.Int64:
			t := time.Unix(v.Int(), 0)
			fieldValue.Set(reflect.ValueOf(t))
			return true
		case reflect.Float64:
			t := time.Unix(int64(v.Float()), 0)
			fieldValue.Set(reflect.ValueOf(t))
			return true
		case reflect.String:
			t := time.Now()
			// make sure date format is correct
			if err := t.UnmarshalJSON([]byte(`"` + v.String() + `"`)); err == nil {
				newValue := reflect.ValueOf(t)
				fieldValue.Set(newValue)
				return true
			}
		}
	case *time.Time:
		if !v.IsValid() {
			var t *time.Time
			newValue := reflect.ValueOf(t)
			fieldValue.Set(newValue)
			return true
		}

		switch v.Kind() {
		case reflect.Int64:
			t := time.Unix(v.Int(), 0)
			fieldValue.Set(reflect.ValueOf(&t))
			return true
		case reflect.Float64:
			t := time.Unix(int64(v.Float()), 0)
			fieldValue.Set(reflect.ValueOf(&t))
			return true
		case reflect.String:
			t := time.Now()
			// make sure date format is correct
			if err := t.UnmarshalJSON([]byte(`"` + v.String() + `"`)); err == nil {
				newValue := reflect.ValueOf(&t)
				fieldValue.Set(newValue)
				return true
			}
		}
	}

	return false
}

func isOverflowUint(field interface{}, value uint64) bool {
	return reflect.ValueOf(field).OverflowUint(value)
}

func isOverflowInt(field interface{}, value int64) bool {
	return reflect.ValueOf(field).OverflowInt(value)
}

func isOverflowFloat(field interface{}, value float64) bool {
	return reflect.ValueOf(field).OverflowFloat(value)
}

// AllUpdaters is a collection of all type updaters
var AllUpdaters = []func(reflect.Value, reflect.Value) bool{
	NullStringUpdater,
	NullFloatUpdater,
	NullIntUpdater,
	NullBoolUpdater,
	NullTimeUpdater,
	MapStringInterfaceUpdater,
	IntUpdater,
	UintUpdater,
	FloatUpdater,
	TimeUpdater,
	BoolUpdater,
}

// Updaters is a collection of standard type updaters
var Updaters = []func(reflect.Value, reflect.Value) bool{
	IntUpdater,
	UintUpdater,
	FloatUpdater,
	TimeUpdater,
	BoolUpdater,
}
