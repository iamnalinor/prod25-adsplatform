package repo

import (
	"os"
	"testing"
)

func TestIsNotFound(t *testing.T) {
	if !IsNotFound(ErrNotFound) {
		t.Errorf("IsNotFound() = false, want true")
	}
	if IsNotFound(nil) {
		t.Errorf("IsNotFound() = true, want false")
	}
	if IsNotFound(os.ErrNotExist) {
		t.Errorf("IsNotFound() = true, want false")
	}
}
