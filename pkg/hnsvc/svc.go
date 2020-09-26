package hnsvc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wenerme/uo/pkg/hnsvc/hnapi"
	"github.com/wenerme/uo/pkg/srpc"

	"github.com/wenerme/uo/pkg/httpmore"
)

// https://github.com/HackerNews/API
type HackerNewsService struct {
	BaseRequest *httpmore.RequestInit
	Client      *http.Client
}

func (HackerNewsService) ServiceCoordinate() srpc.ServiceCoordinate {
	return srpc.ServiceCoordinate{
		ServiceName: "HackerNewsService",
		PackageName: "me.wener.hnsvc",
	}
}

func (s *HackerNewsService) GetItem(id int) (v *hnapi.Item, err error) {
	return v, s.request(httpmore.RequestInit{
		URL: fmt.Sprintf("/item/%v.json", id),
	}, &v)
}

func (s *HackerNewsService) GetUser(name string) (v *hnapi.User, err error) {
	return v, s.request(httpmore.RequestInit{
		URL: fmt.Sprintf("/user/%v.json", name),
	}, &v)
}
func (s *HackerNewsService) MaxItemID() (v int, err error) {
	return v, s.request(httpmore.RequestInit{
		URL: "/maxitem.json",
	}, &v)
}
func (s *HackerNewsService) TopStoryIds() (v []int, err error) {
	return v, s.request(httpmore.RequestInit{
		URL: "/topstories.json",
	}, &v)
}
func (s *HackerNewsService) NewsStoryIds() (v []int, err error) {
	return v, s.request(httpmore.RequestInit{
		URL: "/newstories.json",
	}, &v)
}
func (s *HackerNewsService) BestStoryIds() (v []int, err error) {
	return v, s.request(httpmore.RequestInit{
		URL: "/beststories.json",
	}, &v)
}
func (s *HackerNewsService) AstStoryIds() (v []int, err error) {
	return v, s.request(httpmore.RequestInit{
		URL: "/askstories.json",
	}, &v)
}
func (s *HackerNewsService) ShowStoryIds() (v []int, err error) {
	return v, s.request(httpmore.RequestInit{
		URL: "/showstories.json",
	}, &v)
}
func (s *HackerNewsService) JobStoryIds() (v []int, err error) {
	return v, s.request(httpmore.RequestInit{
		URL: "/jobstories.json",
	}, &v)
}
func (s *HackerNewsService) Updates() (v *hnapi.Updates, err error) {
	return v, s.request(httpmore.RequestInit{
		URL: "/updates.json",
	}, &v)
}

func (s *HackerNewsService) request(r httpmore.RequestInit, out interface{}) error {
	if s.BaseRequest == nil {
		s.BaseRequest = DefaultRequestInit()
	}
	if s.Client == nil {
		s.Client = http.DefaultClient
	}

	base := s.BaseRequest
	client := s.Client
	req, err := base.WithOverride(r).NewRequest()
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(out)
}
func DefaultRequestInit() *httpmore.RequestInit {
	return &httpmore.RequestInit{
		BaseURL: "https://hacker-news.firebaseio.com/v0",
	}
}
