package satws

import "testing"

func TestNewCFDIClientRejectsGarbage(t *testing.T) {
	if _, err := NewCFDIClient([]byte("x"), []byte("y"), "z"); err == nil {
		t.Fatal("expected error")
	}
}

func TestParseKinds(t *testing.T) {
	if _, err := ParseServiceType("nope"); err == nil {
		t.Fatal("expected error")
	}
	if _, err := ParseDownloadType("received"); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseRequestType("metadata"); err != nil {
		t.Fatal(err)
	}
}
