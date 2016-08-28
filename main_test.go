package main

import (
	"testing"
	"os"

	"github.com/stretchr/testify/assert"
)

func TestSetIntEnvvar(t *testing.T) {
	v := 1
	os.Setenv("TEST_SET_INT_ENVVAR", "2")
	setIntEnvvar(&v, "TEST_SET_INT_ENVVAR")
	assert.Equal(t, 2, v, "setIntEnvvar don't override v value")

	/*v = 1
	os.Setenv("TEST_SET_INT_ENVVAR", "not_int")
	setIntEnvvar(&v, "TEST_SET_INT_ENVVAR")
	assert.Equal(t, 1, v, "setIntEnvvar not_int overrided v value")*/

	os.Unsetenv("TEST_SET_INT_ENVVAR")
}