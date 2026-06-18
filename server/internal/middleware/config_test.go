package middleware

import "testing"

func TestParseCORSOrigins(t *testing.T) {
	got := parseCORSOrigins(`"https://www.woragis.me","https://management.woragis.me"`)
	want := []string{"https://www.woragis.me", "https://management.woragis.me"}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestLoadConfigFromEnvStripsQuotes(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", `"https://management.woragis.me",https://www.woragis.me`)
	cfg := LoadConfigFromEnv()
	if len(cfg.CORSOrigins) != 2 {
		t.Fatalf("origins = %v", cfg.CORSOrigins)
	}
	if cfg.CORSOrigins[0] != "https://management.woragis.me" {
		t.Fatalf("first origin = %q", cfg.CORSOrigins[0])
	}
}
