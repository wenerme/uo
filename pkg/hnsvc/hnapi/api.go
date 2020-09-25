package hnapi

import "github.com/wenerme/uo/pkg/srpc"

type HackerNewsServiceClient struct {
	GetItem      func(id int) (v *Item, err error)
	GetUser      func(name string) (v *User, err error)
	MaxItemID    func() (v int, err error)
	TopStoryIds  func() (v []int, err error)
	NewsStoryIds func() (v []int, err error)
	BestStoryIds func() (v []int, err error)
	AskStoryIds  func() (v []int, err error)
	ShowStoryIds func() (v []int, err error)
	JobStoryIds  func() (v []int, err error)
	Updates      func() (v *Updates, err error)
}

func (HackerNewsServiceClient) ServiceCoordinate() srpc.ServiceCoordinate {
	return srpc.ServiceCoordinate{
		ServiceName: "HackerNewsService",
		PackageName: "me.wener.hnsvc",
	}
}

type Item struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Parent      int    `json:"parent,omitempty"`
	Kids        []int  `json:"kids,omitempty"`
	Parts       []int  `json:"parts,omitempty"` // poll
	Score       int    `json:"score,omitempty"`
	Time        int    `json:"time"` // unix epoch
	Title       string `json:"title"`
	Text        string `json:"text,omitempty"`
	Type        string `json:"type"` // story, comment, job, poll, pollopt
	URL         string `json:"url,omitempty"`
	Poll        int    `json:"poll,omitempty"`
}

type User struct {
	ID string `json:"id"`
	// The user's optional self-description. HTML.
	About   string `json:"about"`
	Created int    `json:"created"`
	// Delay in minutes between a comment's creation and its visibility to other users.
	Delay     int   `json:"delay"`
	Karma     int   `json:"karma"`
	Submitted []int `json:"submitted"`
}

type Updates struct {
	Items    []int    `json:"items"`
	Profiles []string `json:"profiles"`
}
