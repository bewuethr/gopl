// Package ch12ex11 provides an encoder for URL query strings that takes struct
// tags for "http" into account.
package ch12ex11

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// Pack encodes the contents of struct s into the query string of URL u.
func Pack(s interface{}, u *url.URL) error {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		k := tag.Get("http")
		if k == "" {
			k = strings.ToLower(fieldInfo.Name)
		}
		if err := writeParam(u, k, v.Field(i)); err != nil {
			return err
		}
	}
	return nil
}

func writeParam(u *url.URL, k string, v reflect.Value) error {
	q := u.Query()
	switch v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if err := writeParam(u, k, v.Index(i)); err != nil {
				return err
			}
		}
		return nil

	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Bool:
		q.Add(k, fmt.Sprint(v))
		u.RawQuery = q.Encode()
		return nil
	}
	return fmt.Errorf("unsupported kind %s", v.Type())
}
