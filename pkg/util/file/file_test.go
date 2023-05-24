package file

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetBasePath(t *testing.T) {
	path := GetBasePath("s/Users/github.com/jeasonstudio/wasmer-k8s/pkg/cmd/create/test.yaml")
	if path != "test.yaml" {
		assert.Error(t, errors.New("GetBasePath error"), "")
	}
	path = GetBasePath("test.yaml")
	if path != "test.yaml" {
		assert.Error(t, errors.New("GetBasePath test.yaml error"), "")
	}
	path = GetBasePath("//test.yaml")
	if path != "test.yaml" {
		assert.Error(t, errors.New("GetBasePath //test.yaml error"), "")
	}
}
