package querystring_test

import (
	"fmt"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"

	"github.com/wenerme/uo/querystring"
)

func TestValues(t *testing.T) {
	vs := "a"
	vi := 1
	var vnil *int = nil
	for _, test := range []struct {
		v interface{}
		o url.Values
	}{
		{v: nil, o: nil},
		{v: map[string]interface{}{
			"a": &vs,
			"b": &vi,
			"c": vnil,
		}, o: map[string][]string{
			"a": {vs},
			"b": {strconv.Itoa(vi)},
		}},
		{v: map[string]interface{}{
			"a": time.Time{},
		}, o: map[string][]string{
			"a": {"0001-01-01T00:00:00Z"},
		}},
		{v: map[string]interface{}{
			"a": 1,
			"b": []string{"a", "b"},
		}, o: map[string][]string{
			"a": {"1"},
			"b": {"a", "b"},
		}},
		{v: map[string]string{
			"a": "1",
			"b": "",
		}, o: map[string][]string{
			"a": {"1"},
			"b": {""},
		}},
		{v: map[string][]string{
			"a": {"1"},
			"b": {"", "a"},
		}, o: map[string][]string{
			"a": {"1"},
			"b": {"", "a"},
		}},
		{v: map[string]string{}, o: nil},
		{v: struct {
			A string   `url:"a"`
			B []string `url:"b"`
		}{
			A: "1",
			B: []string{"a", "b", "c"},
		}, o: map[string][]string{
			"a": {"1"},
			"b": {"a", "b", "c"},
		}},
	} {
		values, err := querystring.Values(test.v)
		assert.NoError(t, err)
		if !assert.Equal(t, fmt.Sprint(test.o), fmt.Sprint(values)) {
			spew.Dump(test.v, values)
		}
	}
}
