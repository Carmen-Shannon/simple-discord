package util

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"

	"github.com/Carmen-Shannon/simple-discord/structs"
)

func BuildQueryString(obj interface{}) string {
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	if v.Kind() != reflect.Struct {
		return ""
	}

	// need to init a bool to check if we need to add a ? or &
	var query string
	first := true

	// scan all the properties of obj
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// ignore non-pointer fields and nil pointers
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			// fetch the actual field name by the json tag
			tag := fieldType.Tag.Get("json")
			if tag == "" {
				tag = fieldType.Name
			} else {
				tag = tag[:len(tag)-len(",omitmepty")]
			}

			// first things
			if first {
				query += "?"
				first = false
			} else {
				query += "&"
			}

			// we can use custom structs for their ToString methods if we implement them
			var val string
			switch field.Interface().(type) {
			case structs.Snowflake:
				snowflakeVal := field.Interface().(structs.Snowflake)
				val = snowflakeVal.ToString()
			case structs.Bitfield[any]:
				bitfieldVal := field.Interface().(structs.Bitfield[any])
				val = bitfieldVal.ToString()
			case *bool:
				boolVal := field.Interface().(*bool)
				val = strconv.FormatBool(*boolVal)
			default:
				val = reflect.Indirect(field).String()
			}

			query += url.QueryEscape(tag) + "=" + url.QueryEscape(val)
		}
	}

	return query
}

func EncodeStructToURL(str interface{}) string {
	raw, err := json.Marshal(str)
	if err != nil {
		return ""
	}

	strStruct := string(raw)
	return url.QueryEscape(strStruct)
}
