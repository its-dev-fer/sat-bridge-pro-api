package service

import "testing"

func TestParseRangeRejectsInverted(t *testing.T) {
	if _, _, err := parseRange("2020-01-01T00:00:00Z", "2019-01-01T00:00:00Z"); err == nil {
		t.Fatal("expected error")
	}
}
