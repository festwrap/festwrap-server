package testtools

import "testing"

func SkipOnShortRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}
