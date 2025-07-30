package router

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRouter_Lookup(t *testing.T) {
	t.Parallel()
	file := filepath.Join("testdata", "policies.yaml")
	f, err := os.CreateTemp("", "policies-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	data, _ := os.ReadFile(file)
	f.Write(data)
	f.Close()
	r, err := NewRouter(f.Name())
	if err != nil {
		t.Fatalf("router init: %v", err)
	}
	backend, ok := r.Lookup("teamA")
	if !ok || backend != "127.0.0.1:10001" {
		t.Errorf("lookup failed: %v %v", backend, ok)
	}
	backend, ok = r.Lookup("teamB")
	if !ok || backend != "127.0.0.1:10002" {
		t.Errorf("lookup failed: %v %v", backend, ok)
	}
	backend, ok = r.Lookup("nope")
	if ok {
		t.Errorf("expected miss for unknown team")
	}
}

func BenchmarkRouter_Lookup(b *testing.B) {
	file := filepath.Join("testdata", "policies.yaml")
	r, err := NewRouter(file)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.Lookup("teamA")
	}
}
