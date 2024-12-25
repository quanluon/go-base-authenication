package utils

import "testing"

func TestHashPassword(t *testing.T) {
	hashedPassword, err := HashPassword("test")
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	if hashedPassword == "" {
		t.Fatalf("Hashed password is empty")
	}
}

func TestComparePassword(t *testing.T) {
	hashedPassword, err := HashPassword("test")
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	if !ComparePassword("test", hashedPassword) {
		t.Fatalf("Password comparison failed")
	}
}
