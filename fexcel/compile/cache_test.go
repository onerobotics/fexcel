package compile

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/onerobotics/fexcel/fexcel"
)

func removeCacheFile() error {
	return os.Remove(CACHE_FILE_NAME)
}

func TestNewCache(t *testing.T) {
	removeCacheFile()
	defer removeCacheFile()

	fpath := "testdata/test.xlsx"
	cfg := fexcel.FileConfig{
		Constants: "G2",
		Numregs:   "A2",
		Posregs:   "D2",
		Sheet:     "Data",
		Offset:    1,
	}

	c, err := newCache(fpath, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if c.Path != fpath {
		t.Errorf("bad cache path. Got %q, want %q", c.Path, fpath)
	}
	if c.FileConfig != cfg {
		t.Errorf("bad cache cfg. Got %v, want %v", c.FileConfig, cfg)
	}
	if c.Version != fexcel.Version {
		t.Errorf("bad version. Got %q, want %q", c.Version, fexcel.Version)
	}
	if len(c.Definitions["R"]) != 3 {
		t.Errorf("bad numreg defs. Got %d, want %d", len(c.Definitions["R"]), 3)
	}
	if len(c.Definitions["PR"]) != 3 {
		t.Errorf("bad posreg defs. Got %d, want %d", len(c.Definitions["PR"]), 3)
	}
	if len(c.Constants) != 2 {
		t.Errorf("bad constant data. Got %d, want %d", len(c.Constants), 2)
	}

	// test cache was saved
	_, err = os.Stat(CACHE_FILE_NAME)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCacheSave(t *testing.T) {
	removeCacheFile()
	defer removeCacheFile()

	modified := time.Now()
	cfg := fexcel.FileConfig{
		Constants: "G2",
		Numregs:   "A2",
		Posregs:   "D2",
		Sheet:     "Data",
		Offset:    1,
	}
	c := &cache{
		Path:       "testdata/test.xlsx",
		FileConfig: cfg,
		Version:    fexcel.Version,
		ModifiedAt: modified,
	}
	err := c.save()
	if err != nil {
		t.Fatal(err)
	}

	// test cache was saved
	_, err = os.Stat(CACHE_FILE_NAME)
	if err != nil {
		t.Fatal(err)
	}

	// loadCached should work then
	err = c.loadCached()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCacheValidate(t *testing.T) {
	removeCacheFile()
	defer removeCacheFile()

	modified := time.Now()
	cfg := fexcel.FileConfig{
		Constants: "G2",
		Numregs:   "A2",
		Posregs:   "D2",
		Sheet:     "Data",
		Offset:    1,
	}
	c := &cache{
		Path:       "testdata/test.xlsx",
		FileConfig: cfg,
		Version:    fexcel.Version,
		ModifiedAt: modified,
	}

	err := c.validate("testdata/test.xlsx", cfg)
	if err != nil {
		t.Fatalf("invalid cache: %v", err)
	}

	c.Path = "some other path"
	err = c.validate("testdata/test.xlsx", cfg)
	if !errors.Is(err, ErrCacheInvalid) {
		t.Errorf("expected ErrCacheInvalid, got %v", err)
	}
	c.Path = "testdata/test.xlsx"

	c.ModifiedAt = time.Time{}
	err = c.validate("testdata/test.xlsx", cfg)
	if !errors.Is(err, ErrCacheInvalid) {
		t.Errorf("expected ErrCacheInvalid, got %v", err)
	}
	c.ModifiedAt = time.Now()

	c.Version = "foo"
	err = c.validate("testdata/test.xlsx", cfg)
	if !errors.Is(err, ErrCacheInvalid) {
		t.Errorf("expected ErrCacheInvalid, got %v", err)
	}
	c.Version = fexcel.Version

	c.FileConfig = fexcel.FileConfig{}
	err = c.validate("testdata/test.xlsx", cfg)
	if !errors.Is(err, ErrCacheInvalid) {
		t.Errorf("expected ErrCacheInvalid, got %v", err)
	}
}

func TestLoadCachedReturnsIfCacheFileDoesNotExist(t *testing.T) {
	removeCacheFile()
	defer removeCacheFile()

	c := &cache{}
	err := c.loadCached()
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("cache file should not exist")
	}
}

func TestLoadCachedReturnsErrorOnInvalidCache(t *testing.T) {
	removeCacheFile()
	defer removeCacheFile()

	modified := time.Now()
	cfg := fexcel.FileConfig{
		Constants: "G2",
		Numregs:   "A2",
		Posregs:   "D2",
		Sheet:     "Data",
		Offset:    1,
	}
	c := &cache{
		Path:       "testdata/test.xlsx",
		FileConfig: cfg,
		Version:    fexcel.Version,
		ModifiedAt: modified,
	}
	err := c.save()
	if err != nil {
		t.Fatal(err)
	}
	c.FileConfig.Constants = "H2"
	err = c.loadCached()
	if !errors.Is(err, ErrCacheInvalid) {
		t.Errorf("Expected ErrCacheInvalid got %v", err)
	}
}
