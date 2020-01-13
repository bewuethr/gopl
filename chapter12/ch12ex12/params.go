// Package ch12ex12 provides a reflection-based parser for URL parameters with
// support for parameter validation.
package ch12ex12

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type fieldProps struct {
	f     reflect.Value
	valid string
}

// Unpack populates the fields of the struct pointed to by ptr from the HTTP
// request parameters in req. The http tag supports the validation options
// "email" and "pan" for strings, and "zipcode" for integers.
func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	// Build map of fields keyed by effective name.
	fields := make(map[string]fieldProps)
	v := reflect.ValueOf(ptr).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		httpTag := tag.Get("http")
		name, option := getOption(httpTag)
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = fieldProps{
			f:     v.Field(i),
			valid: option,
		}
	}

	// Update struct field for each parameter in the request.
	for name, values := range req.Form {
		f, valid := fields[name].f, fields[name].valid
		if !f.IsValid() {
			continue // ignore unrecognized HTTP parameters
		}
		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
				if err := validate(value, valid, elem); err != nil {
					return fmt.Errorf("validation failed: %v", err)
				}
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				if err := validate(value, valid, f); err != nil {
					return fmt.Errorf("validation failed: %v", err)
				}
			}
		}
	}
	return nil
}

// getOption extracts the validation option from the tag string if present.
func getOption(tag string) (string, string) {
	if i := strings.Index(tag, ","); i != -1 {
		fields := strings.Split(tag, ",")
		return fields[0], fields[1]
	}
	return tag, ""
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}

func validate(param, valid string, v reflect.Value) error {
	switch valid {
	case "":
		return nil
	case "email":
		if v.Kind() != reflect.String {
			return fmt.Errorf("email validation can only be applied to string, not %v", v.Type().Kind())
		}
		return isValidEmail(param)
	case "pan":
		if v.Kind() != reflect.String {
			return fmt.Errorf("PAN validation can only be applied to string, not %v", v.Type().Kind())
		}
		return isValidPAN(param)
	case "zipcode":
		if v.Kind() != reflect.Int {
			return fmt.Errorf("ZIP code validation can only be applied to integer, not %v", v.Type().Kind())
		}
		zipCode, _ := strconv.ParseInt(param, 10, 64) // error was already checked
		return isValidZIPCode(zipCode)
	default:
		return fmt.Errorf("invalid validation option %q", valid)
	}
}

var emailRE, panRE *regexp.Regexp

func init() {
	const (
		emailREStr = `^[[:alnum:]]+(?:\.[[:alnum:]]+)*(?:\+[[:alnum:]]+)?@(?:[[:alnum:]]+\.)+[[:alpha:]]{2,}$`
		panREStr   = `[^[:digit:]]`
	)
	emailRE = regexp.MustCompile(emailREStr)
	panRE = regexp.MustCompile(panREStr)
}

// isValidEmail checks if string s is a valid email address (using a simple
// regex).
func isValidEmail(s string) error {
	if !emailRE.MatchString(s) {
		return fmt.Errorf("not a valid email address: %v", s)
	}
	return nil
}

// isValidPAN checks if s contains any non-digits, has a length of exactly 16
// and is valid according to the Luhn algorithm.
func isValidPAN(s string) error {
	const panLen = 16

	if panRE.MatchString(s) {
		return fmt.Errorf("PAN %v contains non-digits", s)
	}

	if len(s) != panLen {
		return fmt.Errorf("PAN %v has %d digits, want %d", s, len(s), panLen)
	}

	// Luhn algorithm
	var sum int
	for i, c := range strings.Split(s, "") {
		d, err := strconv.Atoi(c)
		if err != nil { // shouldn't happen, we've checked for non-digits
			return err
		}
		if i%2 == 0 {
			d *= 2
			if d >= 10 {
				d -= 9
			}
		}
		sum += d
	}

	if sum%10 != 0 {
		return fmt.Errorf("PAN %v is invalid according to Luhn algorithm", s)
	}

	return nil
}

// isValidZIPCode just checks if i has 5 digits.
func isValidZIPCode(i int64) error {
	if i < 10000 || i > 99999 {
		return fmt.Errorf("ZIP code %d does not have 5 digits", i)
	}
	return nil
}
