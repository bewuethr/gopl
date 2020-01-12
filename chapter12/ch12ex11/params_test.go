package ch12ex11

import (
	"net/url"
	"testing"
)

type data struct {
	Labels     []string `http:"l"`
	MaxResults uint     `http:"max"`
	Exact      bool     `http:"x"`
}

type dataExt struct {
	Labels     []string `http:"l"`
	MaxResults uint     `http:"max"`
	Exact      bool     `http:"x"`
	NoTag      string
}

func TestPack(t *testing.T) {
	var tests = []struct {
		desc  string
		u     *url.URL
		input interface{}
		want  string
	}{
		{
			desc: "empty query string",
			u:    new(url.URL),
			input: data{
				Labels:     []string{"l1", "l2"},
				MaxResults: 5,
				Exact:      true,
			},
			want: "l=l1&l=l2&max=5&x=true",
		},
		{
			desc: "non-empty query string",
			u:    &url.URL{RawQuery: "l=l0"},
			input: data{
				Labels:     []string{"l1", "l2"},
				MaxResults: 5,
				Exact:      true,
			},
			want: "l=l0&l=l1&l=l2&max=5&x=true",
		},
		{
			desc: "field without tag",
			u:    new(url.URL),
			input: dataExt{
				Labels:     []string{"l1", "l2"},
				MaxResults: 5,
				Exact:      true,
				NoTag:      "abc",
			},
			want: "l=l1&l=l2&max=5&notag=abc&x=true",
		},
	}

	for _, test := range tests {
		if err := Pack(test.input, test.u); err != nil {
			t.Errorf("%v: got error %v, expected nil", test.desc, err)
			continue
		}
		if got := test.u.Query().Encode(); got != test.want {
			t.Errorf("%v: Pack(%+v) == %v, want %v", test.desc, test.input, got, test.want)
		}
	}
}
