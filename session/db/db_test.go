package db

import "testing"

func TestRandomString(t *testing.T) {
	str, err := RandomString(37)
	if err != nil {
		t.Error(err)
	}
	if len(str) < 37 {
		t.Errorf("Unexpected random string length. Waiting for >= 47 got %d", len(str))
	}
}
