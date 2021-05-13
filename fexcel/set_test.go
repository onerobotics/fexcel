package fexcel

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestNewSetCommandErrors(t *testing.T) {
	tests := []struct {
		fpath   string
		cfg     Config
		targets []string
		err     string
	}{
		{"./testdata/test.xlsx", Config{}, []string{}, "Need at least one target"},
		{"./testdata/test.xlsx", Config{}, []string{"foo"}, "no cell locations defined"},
		{"./testdata/test.xlsx", Config{FileConfig: FileConfig{Numregs: "A2"}}, []string{"./testdata"}, "offset must be nonzero"},
		{"./testdata/test.xlsx", Config{FileConfig: FileConfig{Numregs: "A2", Offset: 1}}, []string{"./testdata"}, "\"./testdata\" is not a valid remote host"},
	}

	for id, test := range tests {
		_, err := NewSetCommand(test.fpath, test.cfg, test.targets...)
		if err == nil {
			t.Errorf("case(%d): expected an error", id)
			continue
		}

		if err.Error() != test.err {
			t.Errorf("bad error. Got %q, want %q", err.Error(), test.err)
		}
	}

	// valid
	_, err := NewSetCommand("./testdata/test.xlsx", Config{FileConfig: FileConfig{Numregs: "Data:A2", Offset: 1}}, "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetCommand(t *testing.T) {
	var commentCount uint32
	hf := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/KAREL/ComSet":
			rw.Write([]byte("OK"))
			atomic.AddUint32(&commentCount, 1)
		case "/MD/numreg.va":
			rw.Write([]byte(" [1] = 0  'this is an extre'\n"))
		default:
			http.Error(rw, "Not implemented", http.StatusNotImplemented)
			t.Fatalf("Unexpected request: %q", req.URL.Path)
		}
	})

	s1 := httptest.NewServer(hf)
	defer s1.Close()
	s2 := httptest.NewServer(hf)
	defer s2.Close()

	hosts := []string{s1.URL, s2.URL}

	s, err := NewSetCommand("./testdata/test.xlsx", Config{FileConfig: FileConfig{Numregs: "Data:A2", Offset: 1}}, hosts...)
	if err != nil {
		t.Fatal(err)
	}

	result, err := s.Execute()
	if err != nil {
		t.Fatal(err)
	}

	for _, host := range hosts {
		if count := result.Counts[host][Numreg]; count != 4 {
			t.Errorf("Result.Counts[%s][%s]: Got %d, want 4", host, Numreg, count)
		}
	}

	// (5 defs - 1 accurate) * 2 hosts = 8
	if commentCount != 8 {
		t.Errorf("comment request called %d times. Want 8", commentCount)
	}
}
