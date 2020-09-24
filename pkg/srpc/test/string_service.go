package test

import "strings"

type StringService struct {
}
type StringPart struct {
	A   string
	B   string
	Sep string `json:"sep,omitempty"`
}

func (*StringService) Uppercase(s string) (string, error) {
	return strings.ToUpper(s), nil
}
func (*StringService) UppercasePtr(s *string) (string, error) {
	return strings.ToUpper(*s), nil
}
func (*StringService) Join(r StringPart) (string, error) {
	return r.A + r.Sep + r.B, nil
}
func (s *StringService) JoinPtr(r *StringPart) (string, error) {
	if r == nil {
		r = &StringPart{}
	}
	return s.Join(*r)
}
func (*StringService) Sep(r string) (StringPart, error) {
	sep := strings.SplitN(r, ".", 2)
	return StringPart{
		A:   sep[0],
		B:   sep[1],
		Sep: ".",
	}, nil
}
func (s *StringService) SepPtr(r string) (*StringPart, error) {
	re, _ := s.Sep(r)
	return &re, nil
}

type StringServiceClient struct {
	Uppercase    func(string) (string, error)
	UppercasePtr func(string) (string, error)
	Join         func(r StringPart) (string, error)
	JoinPtr      func(r *StringPart) (string, error)
	Sep          func(r string) (StringPart, error)
	SepPtr       func(r string) (*StringPart, error)
}
