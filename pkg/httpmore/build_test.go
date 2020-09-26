package httpmore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wenerme/uo/pkg/httpmore"
)

func TestOverride(t *testing.T) {
	{
		assert.Equal(t,
			httpmore.RequestInit{}.WithOverride(httpmore.RequestInit{Options: httpmore.Values{}.Add("1", "2").Add("1", "3")}),
			httpmore.RequestInit{Options: httpmore.Values{}.Add("1", "2")}.WithOverride(httpmore.RequestInit{Options: httpmore.Values{}.Add("1", "3")}),
		)
	}
}

func TestUrlBuild(t *testing.T) {
	{
		r, err := httpmore.RequestInit{
			BaseURL: "https://wener.me",
			URL:     "/token",
		}.NewRequest()
		assert.NoError(t, err)
		assert.Equal(t, "https://wener.me/token", r.URL.String())
	}
	{
		r, err := httpmore.RequestInit{
			BaseURL: "https://wener.me",
			URL:     "/token",
			Query: map[string][]string{
				"name": {"wener"},
			},
		}.NewRequest()
		assert.NoError(t, err)
		assert.Equal(t, "https://wener.me/token?name=wener", r.URL.String())
	}
	{
		r, err := httpmore.RequestInit{
			BaseURL: "https://wener.me",
			URL:     "/token",
			Query: map[string]string{
				"name": "wener",
			},
		}.NewRequest()
		assert.NoError(t, err)
		assert.Equal(t, "https://wener.me/token?name=wener", r.URL.String())
	}
	{
		r, err := httpmore.RequestInit{
			BaseURL: "https://wener.me",
			URL:     "/token",
			Query: map[string]interface{}{
				"name": "wener",
				"age":  18,
			},
		}.NewRequest()
		assert.NoError(t, err)
		assert.Equal(t, "https://wener.me/token?age=18&name=wener", r.URL.String())
	}
}
