package hnsvc_test

import (
	stdlog "log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wenerme/uo/pkg/hnsvc"
)

func TestBasic(t *testing.T) {
	// os.Setenv("HTTP_PROXY", "socks5://127.0.0.1:8888")
	s := &hnsvc.HackerNewsService{}
	{
		it, err := s.GetItem(8863)
		assert.NoError(t, err)
		assert.Equal(t, 8863, it.ID)
		assert.Equal(t, "My YC app: Dropbox - Throw away your USB drive", it.Title)
	}
	{
		v, err := s.GetUser("wener")
		assert.NoError(t, err)
		assert.Equal(t, "wener", v.ID)
	}
	{
		v, err := s.MaxItemID()
		assert.NoError(t, err)
		assert.Greater(t, v, 100000)
		stdlog.Printf("Max Item ID %v", v)
	}
	{
		v, err := s.TopStoryIds()
		assert.NoError(t, err)
		assert.NotEmpty(t, v)
	}
	{
		v, err := s.NewsStoryIds()
		assert.NoError(t, err)
		assert.NotEmpty(t, v)
	}
	{
		v, err := s.BestStoryIds()
		assert.NoError(t, err)
		assert.NotEmpty(t, v)
	}
	{
		v, err := s.AstStoryIds()
		assert.NoError(t, err)
		assert.NotEmpty(t, v)
	}
	{
		v, err := s.ShowStoryIds()
		assert.NoError(t, err)
		assert.NotEmpty(t, v)
	}
	{
		v, err := s.JobStoryIds()
		assert.NoError(t, err)
		assert.NotEmpty(t, v)
	}
	{
		v, err := s.Updates()
		assert.NoError(t, err)
		assert.NotEmpty(t, v.Items)
		assert.NotEmpty(t, v.Profiles)
		stdlog.Printf("Updates %#v", v)
	}
}
