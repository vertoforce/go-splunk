package splunk

import (
	"fmt"
	"reflect"
	"time"
)

// UnMarshal Fills in a struct with the splunk result.  The passed in interface must be a pointer to a struct
//
// It will fill each field in the struct looking for the "splunk" tag on the field based on the splunk field
func (e SearchResult) UnMarshal(v interface{}) error {
	typeOf := reflect.TypeOf(v).Elem()
	valueOf := reflect.ValueOf(v).Elem()

	for i := 0; i < typeOf.NumField(); i++ {
		thisField := typeOf.Field(i)
		label := thisField.Name
		if splunkLabel := typeOf.Field(i).Tag.Get("splunk"); splunkLabel != "" {
			label = splunkLabel
		}

		// Find this field in our SearchResult
		for key, value := range e {
			if key == label {
				// We found this field!  Set it in the struct
				switch thisField.Type.Kind() {
				case reflect.String:
					valueOf.Field(i).SetString(fmt.Sprintf("%v", value))
				case reflect.TypeOf(time.Time{}).Kind():
					// Parse the time and set it
					parsedTime, err := ParseTime(fmt.Sprintf("%v", value))
					if err != nil {
						return fmt.Errorf("Could not parse a time field from splunk that we want to fill as time in the struct.  splunk field: %s, struct field: %s, value: %v", key, thisField.Name, value)
					}
					valueOf.Field(i).Set(reflect.ValueOf(parsedTime))
				case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
					// Make sure the value is a float64
					floatValue, ok := value.(float64)
					if !ok {
						return fmt.Errorf("Could not parse a number field from splunk that we want to fill as a number in the struct.  splunk field: %s, struct field: %s, value: %v", key, thisField.Name, value)
					}
					valueOf.Field(i).SetInt(int64(floatValue))
				case reflect.Float32, reflect.Float64:
					// Make sure the value is a float64
					floatValue, ok := value.(float64)
					if !ok {
						return fmt.Errorf("Could not parse a number field from splunk that we want to fill as a number in the struct.  splunk field: %s, struct field: %s, value: %v", key, thisField.Name, value)
					}
					valueOf.Field(i).SetFloat(floatValue)
				case reflect.Interface:
					valueOf.Field(i).Set(reflect.ValueOf(value))
				}
			}
		}
	}
	return nil
}
