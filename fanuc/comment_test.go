package fanuc

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestCommentToolUpdater(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("OK"))
		want := "/karel/ComSet?sComment=foo&sIndx=1&sFc=1"
		if req.URL.String() != want {
			t.Errorf("Bad URL. Got %q, want %q", req.URL, want)
		}
	}))
	defer server.Close()

	c := &CommentToolUpdater{100 * time.Millisecond}

	host := server.URL[7:] // ignore http://
	err := c.Update(Definition{Numreg, 1, "foo"}, host)
	if err != nil {
		t.Error(err)
	}
}

func TestMultiUpdater(t *testing.T) {
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

	c := &CommentToolUpdater{100 * time.Millisecond}

	hosts := []string{s1.URL[7:], s2.URL[7:]} // get rid of http://
	mu := NewMultiUpdater(hosts, c)

	err := mu.Update([]Definition{
		Definition{Numreg, 1, "foo"},
		Definition{Numreg, 1, "foo"},
	})
	if err != nil {
		t.Error(err)
	}

	if hfCallCount != 4 {
		t.Errorf("handlerFunc only called %d times. Want 2", hfCallCount)
	}
}
