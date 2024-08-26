package util

import (
	"encoding/json"
	"net/url"
	"reflect"

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

			// since snowflake is special, need to take care of snowflake
			// TODO: add other custom types here when needed
			var val string
			if snowflake, ok := field.Interface().(structs.Snowflake); ok {
				val = snowflake.ToString()
			} else {
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
