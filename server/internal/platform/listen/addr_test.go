package listen

import (
	"testing"
)

func TestAddr(t *testing.T) {
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("PORT", "")

	if got := Addr(); got != ":8080" {
		t.Fatalf("default = %q, want :8080", got)
	}

	t.Setenv("PORT", "3000")
	if got := Addr(); got != ":3000" {
		t.Fatalf("PORT = %q, want :3000", got)
	}

	t.Setenv("HTTP_ADDR", ":9090")
	if got := Addr(); got != ":9090" {
		t.Fatalf("HTTP_ADDR override = %q, want :9090", got)
	}
}
