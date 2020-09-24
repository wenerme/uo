package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func StringServiceClientSpec(t *testing.T, svc *StringServiceClient) {
	{
		reply, err := svc.Uppercase("a")
		assert.NoError(t, err)
		assert.Equal(t, "A", reply)
	}
	{
		reply, err := svc.Join(StringPart{
			A:   "a",
			B:   "B",
			Sep: "|",
		})
		assert.NoError(t, err)
		assert.Equal(t, "a|B", reply)
	}
	{
		reply, err := svc.Sep("a.B")
		assert.NoError(t, err)
		assert.Equal(t, StringPart{A: "a", B: "B", Sep: "."}, reply)
	}
	// Ptr
	{
		reply, err := svc.UppercasePtr("aBc")
		assert.NoError(t, err)
		assert.Equal(t, "ABC", reply)
	}
	{
		reply, err := svc.JoinPtr(&StringPart{
			A:   "a",
			B:   "B",
			Sep: "|",
		})
		assert.NoError(t, err)
		assert.Equal(t, "a|B", reply)
	}
	{
		reply, err := svc.SepPtr("a.B")
		assert.NoError(t, err)
		assert.Equal(t, &StringPart{A: "a", B: "B", Sep: "."}, reply)
	}
}
