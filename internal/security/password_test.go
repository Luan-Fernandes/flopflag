package security

import "testing"

func TestHashPasswordAndCompare(t *testing.T) {
	hash, err := HashPassword("12345678")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}

	if err := ComparePasswordHash("12345678", hash); err != nil {
		t.Fatalf("ComparePasswordHash should pass: %v", err)
	}

	if err := ComparePasswordHash("wrong", hash); err == nil {
		t.Fatalf("ComparePasswordHash should fail for wrong password")
	}
}
