package gopartial

import (
	"errors"
	"fmt"
	"reflect"
)

var errDestinationMustBeStructType = errors.New("Destination must be a struct type")
var errDestinationMustBePointerType = errors.New("Destination must be pointer to struct")

// PartialUpdate updates destination object (Must be a pointer to a struct)
// from a map[string]interface{} where struct tag name is equals to the map key.
// This function can extended through updaters. A list of function that accepts
// destination Value and the to be assigned Value and return true if updates is successful
// Returns list of struct field names that was successfully updated.
// If tagName is not provided, the default lookup value would be the field's name.
func PartialUpdate(dest interface{}, partial map[string]interface{}, tagName string, skipConditions []func(reflect.StructField) bool, updaters []func(reflect.Value, reflect.Value) bool) ([]string, error) {
	valueOfDest := reflect.ValueOf(dest)
	// Must be a pointer to a struct so that it can be updated
	if valueOfDest.Kind() != reflect.Ptr {
		return nil, errDestinationMustBePointerType
	}
	valueOfDest = valueOfDest.Elem()

	typeOfDest := valueOfDest.Type()
	// Must be a pointer to a struct so that it can be updated
	if typeOfDest.Kind() != reflect.Struct {
		return nil, errDestinationMustBeStructType
	}

	// fieldsUpdated is to keep track all the field names that were successfuly updated
	fieldsUpdated := make([]string, 0)

	for i := 0; i < typeOfDest.NumField(); i++ {
		field := typeOfDest.Field(i)

		// skip this field if it cant be set
		if !valueOfDest.Field(i).CanSet() {
			continue
		}

		skip := false
		// go through all extended skip conditions
		for _, skipCondition := range skipConditions {
			skip = skipCondition(field)
			if skip {
				// break on the first skip condition found
				break
			}
		}
		if skip {
			continue
		}

		fieldName := field.Name
		if tagName != "" {
			fieldName = field.Tag.Get(tagName)
		}

		// get the partial value based on the tagName
		if val, ok := partial[fieldName]; ok {
			v := reflect.ValueOf(val)
			updateSuccess := false

			if valueOfDest.Field(i).Kind() == reflect.Slice {
				updateSuccess = SliceUpdater(valueOfDest.Field(i), v)
				if updateSuccess {
					updateSuccess = true
				}
			} else if valueOfDest.Field(i).Kind() == v.Kind() {
				valueOfDest.Field(i).Set(v)
				updateSuccess = true
			} else {
				// go through all extended process types
				for _, updater := range updaters {
					updateSuccess = updater(valueOfDest.Field(i), v)
					if updateSuccess {
						// the first updateSuccess found, break the loop
						break
					}
				}
			}

			if updateSuccess {
				fieldsUpdated = append(fieldsUpdated, field.Name)
			} else {
				if !v.IsValid() {
					errMsg := fmt.Sprintf("%v.%v cannot be assigned with value null", typeOfDest.Name(), field.Name)
					return nil, errors.New(errMsg)
				} else {
					errMsg := fmt.Sprintf("%v.%v cannot be assigned with value %v", typeOfDest.Name(), field.Name, v.Interface())
					return nil, errors.New(errMsg)
				}
			}
		}

	}

	return fieldsUpdated, nil
}
