package internal_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"testing"
)

var updateGolden = flag.Bool("update-golden", false, "update golden files")

func GoldenBytes(t *testing.T, name string, replacementContent []byte) []byte {
	t.Helper()

	golden := path.Join("test-fixtures", fmt.Sprintf("%s.golden", name))

	if *updateGolden {
		ioutil.WriteFile(golden, replacementContent, 0644)
	}

	content, _ := ioutil.ReadFile(golden)

	return content
}
