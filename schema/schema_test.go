package schema_test

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/dcrichards/go-to-openapi/schema"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

func TestSimple(t *testing.T) {
	type Tag struct {
		Name   string `json:"name"`
		Active bool   `json:"active"`
	}

	type Example struct {
		ID         string            `json:"id"`
		Email      string            `json:"email"`
		Tags       []Tag             `json:"tags"`
		Properties map[string]string `json:"props"`
	}

	actual, err := schema.Generate(Example{})
	if err != nil {
		t.Error(err)
	}

	golden := filepath.Join("testdata", fmt.Sprintf("%s.golden.yml", t.Name()))
	expected, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(expected, []byte(actual)) {
		t.Errorf("\n[EXPECTED]\n %s \n[ACTUAL]\n %s", string(expected), actual)
	}

	if *update {
		if err := ioutil.WriteFile(golden, []byte(actual), 0644); err != nil {
			t.Error(err)
		}
	}
}
