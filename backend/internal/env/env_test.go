package env

import (
	"festwrap/internal/testtools"
	"os"
	"testing"
)

func TestGetEnvReturnsDefaultValueIfDoesntExist(t *testing.T) {
	t.Run("integer", func(t *testing.T) {
		t.Parallel()
		defaultValue := 2
		value, _ := GetEnvWithDefault("MY_KEY", defaultValue)

		testtools.AssertEqual(t, value, defaultValue)
	})

	t.Run("string", func(t *testing.T) {
		t.Parallel()
		defaultValue := "something"
		value, _ := GetEnvWithDefault("MY_KEY", defaultValue)

		testtools.AssertEqual(t, value, defaultValue)
	})
}

func TestGetEnvReturnsExistingEnvVariable(t *testing.T) {
	t.Run("integer", func(t *testing.T) {
		key := "myKey"
		value := "42"
		os.Setenv(key, value)
		actual, _ := GetEnvWithDefault(key, 0)

		testtools.AssertEqual(t, actual, 42)
	})

	t.Run("string", func(t *testing.T) {
		key := "myKey"
		value := "my_value"
		os.Setenv(key, value)
		actual, _ := GetEnvWithDefault(key, "")

		testtools.AssertEqual(t, actual, value)
	})
}

func TestGetEnvReturnsErrorOnInvalidEnvVar(t *testing.T) {
	key := "myKey"
	nonIntValue := "acd"
	os.Setenv(key, nonIntValue)

	_, err := GetEnvWithDefault(key, 0)

	testtools.AssertErrorIsNotNil(t, err)
}
