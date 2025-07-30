package config

import (
	"path/filepath"
	"testing"
)

func TestLoad_Valid(t *testing.T) {
	t.Parallel()
	path := filepath.Join("testdata", "valid.yaml")
	c, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c.Listen != ":9000" {
		t.Errorf("unexpected listen: %s", c.Listen)
	}
	if c.JWTSecret != "testsecret" {
		t.Errorf("unexpected jwt_secret: %s", c.JWTSecret)
	}
	if len(c.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(c.Policies))
	}
}

func TestLoad_Invalid(t *testing.T) {
	t.Parallel()
	path := filepath.Join("testdata", "invalid.yaml")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
