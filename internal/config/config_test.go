package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("DECENT_EXPORTER_URL", "")
	got, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.ListenAddress != ":8080" {
		t.Fatalf("listen address = %q", got.ListenAddress)
	}
	if got.DecentURL != "http://127.0.0.1:8080" {
		t.Fatalf("decent url = %q", got.DecentURL)
	}
}
