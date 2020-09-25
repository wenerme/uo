package httpmore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wenerme/uo/pkg/httpmore"
)

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
}
