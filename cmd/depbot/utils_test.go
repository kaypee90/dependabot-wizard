package main

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_createDirIfItDoesNotExit(t *testing.T) {
	dircreateDirIfItDoesNotExit("./test")
	assert.NilError(t, err)
}
