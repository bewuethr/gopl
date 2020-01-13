package ch12ex12

import (
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type destType struct {
	Email   string `http:"e,email"`
	PAN     string `http:"p,pan"`
	ZIPCode int    `http:"z,zipcode"`
}

func TestUnpack(t *testing.T) {
	var tests = []struct {
		desc      string
		qs        string
		got, want destType
		wantErr   string
	}{
		{
			desc: "all valid fields",
			qs:   "e=" + url.QueryEscape("foobar@example.com") + "&p=4492086053472425&z=12345",
			want: destType{
				Email:   "foobar@example.com",
				PAN:     "4492086053472425",
				ZIPCode: 12345,
			},
		},
		{
			desc: "valid email address",
			qs:   "e=" + url.QueryEscape("foobar@example.com"),
			want: destType{Email: "foobar@example.com"},
		},
		{
			desc: "valid PAN",
			qs:   "p=4492086053472425",
			want: destType{PAN: "4492086053472425"},
		},
		{
			desc: "valid ZIP code",
			qs:   "z=12345",
			want: destType{ZIPCode: 12345},
		},
		{
			desc:    "invalid email address",
			qs:      "e=" + url.QueryEscape("foobar@example.c"),
			wantErr: "validation failed: not a valid email address",
		},
		{
			desc:    "invalid PAN",
			qs:      "p=4492086053472426",
			wantErr: "validation failed: PAN 4492086053472426 is invalid",
		},
		{
			desc:    "PAN with non-digits",
			qs:      "p=111111111111111a",
			wantErr: "validation failed: PAN 111111111111111a contains non-digits",
		},
		{
			desc:    "PAN with incorrect length",
			qs:      "p=4492",
			wantErr: "validation failed: PAN 4492 has 4 digits, want 16",
		},
		{
			desc:    "invalid ZIP code",
			qs:      "z=1234",
			wantErr: "validation failed: ZIP code 1234 does not have 5 digits",
		},
	}

	for _, test := range tests {
		req := &http.Request{}
		req.URL = &url.URL{RawQuery: test.qs}
		err := Unpack(req, &test.got)
		if test.wantErr == "" {
			if err != nil {
				t.Errorf("%v: got error %v, expected nil", test.desc, err)
				continue
			}

			if !reflect.DeepEqual(test.got, test.want) {
				t.Errorf("%v: Unpack(%v) == %+v, want %+v", test.desc, test.qs, test.got, test.want)
			}
		} else {
			if err == nil {
				t.Errorf("%v: got nil error, want %v", test.desc, test.wantErr)
				continue
			}
			if !strings.HasPrefix(err.Error(), test.wantErr) {
				t.Errorf("%v: got %v, want %v", test.desc, err, test.wantErr)
			}
		}
	}
}

type invDestType struct {
	Einv int    `http:",email"`   // invalid email struct
	Pinv int    `http:",pan"`     // invalid PAN struct
	Zinv string `http:",zipcode"` // invalid ZIP code struct
	Inv  string `http:",invalid"` // invalid validation
}

func TestUnpackInvalidTags(t *testing.T) {
	var tests = []struct {
		desc string
		qs   string
		dest invDestType
		want string
	}{
		{
			desc: "email validation on non-string",
			qs:   "einv=100",
			want: "validation failed: email validation can only be applied to string",
		},
		{
			desc: "PAN validation on non-string",
			qs:   "pinv=100",
			want: "validation failed: PAN validation can only be applied to string",
		},
		{
			desc: "ZIP code validation on non-integer",
			qs:   "zinv=abc",
			want: "validation failed: ZIP code validation can only be applied to integer",
		},
		{
			desc: "invalid validation option",
			qs:   "inv=abc",
			want: `validation failed: invalid validation option "invalid"`,
		},
	}

	for _, test := range tests {
		req := &http.Request{}
		req.URL = &url.URL{RawQuery: test.qs}
		err := Unpack(req, &test.dest)
		if err == nil {
			t.Errorf("%v: got nil error, want %v", test.desc, test.want)
			continue
		}

		if !strings.HasPrefix(err.Error(), test.want) {
			t.Errorf("%v: got %v, want %v", test.desc, err, test.want)
		}
	}
}
