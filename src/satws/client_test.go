package satws

import "testing"

func TestNewCFDIClientRejectsGarbage(t *testing.T) {
	if _, err := NewCFDIClient([]byte("x"), []byte("y"), "z"); err == nil {
		t.Fatal("expected error")
	}
}
