package auth

import (
	"strings"
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	t.Parallel()
	s := NewSigner("testkey")
	token, err := s.Sign("1.2.3.4", "teamA")
	if err != nil {
		t.Fatalf("sign error: %v", err)
	}
	ip, team, err := s.Verify(token)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}
	if ip != "1.2.3.4" || team != "teamA" {
		t.Errorf("unexpected claims: %s %s", ip, team)
	}
}

func FuzzVerifyRejectsAltered(f *testing.F) {
	s := NewSigner("testkey")
	token, _ := s.Sign("1.2.3.4", "teamA")
	f.Add(token)
	f.Fuzz(func(t *testing.T, orig string) {
		if len(orig) < 10 {
			return
		}
		altered := orig[:len(orig)/2] + "x" + orig[len(orig)/2+1:]
		ip, team, err := s.Verify(altered)
		if err == nil && (ip != "1.2.3.4" || team != "teamA") {
			t.Errorf("expected error for altered token")
		}
		if strings.HasSuffix(orig, "=") {
			_, _, err = s.Verify(orig[:len(orig)-1])
			if err == nil {
				t.Errorf("expected error for truncated token")
			}
		}
	})
}
