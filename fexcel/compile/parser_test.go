package compile

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	filenames := []string{"test.ls"}

	for _, filename := range filenames {
		fpath := filepath.Join("testdata", filename)
		src, err := ioutil.ReadFile(fpath)
		if err != nil {
			t.Fatal(err)
		}

		_, err = Parse(filename, string(src))
		if err != nil {
			t.Errorf("Parse(%s): %s", filename, err)
		}
	}
}
