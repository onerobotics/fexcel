package fexcel

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	fanuc "github.com/onerobotics/go-fanuc"
)

func TestCommentToolSetter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("OK"))
		want := "/karel/ComSet?sComment=foo&sIndx=1&sFc=1"
		if req.URL.String() != want {
			t.Errorf("Bad URL. Got %q, want %q", req.URL, want)
		}
	}))
	defer server.Close()

	c := &CommentToolSetter{100 * time.Millisecond}

	host := server.URL[7:] // ignore http://
	err := c.Set(Definition{fanuc.Numreg, 1, "foo"}, host)
	if err != nil {
		t.Error(err)
	}
}

func TestMultiSetter(t *testing.T) {
	var hfCallCount uint32
	hf := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("OK"))
		want := "/karel/ComSet?sComment=foo&sIndx=1&sFc=1"
		if req.URL.String() != want {
			t.Errorf("Bad URL. Got %q, want %q", req.URL, want)
		}
		atomic.AddUint32(&hfCallCount, 1)
	})

	s1 := httptest.NewServer(hf)
	defer s1.Close()
	s2 := httptest.NewServer(hf)
	defer s2.Close()

	c := &CommentToolSetter{100 * time.Millisecond}

	hosts := []string{s1.URL[7:], s2.URL[7:]} // get rid of http://
	ms := NewMultiSetter(hosts, c)

	defs := make(map[string][]Definition)
	defs[hosts[0]] = []Definition{Definition{fanuc.Numreg, 1, "foo"}}
	defs[hosts[1]] = []Definition{Definition{fanuc.Numreg, 1, "foo"}}
	err := ms.Set(defs)
	if err != nil {
		t.Error(err)
	}

	if hfCallCount != 2 {
		t.Errorf("handlerFunc only called %d times. Want 2", hfCallCount)
	}
}
