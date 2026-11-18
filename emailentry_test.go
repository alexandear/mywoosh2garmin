package main

import "testing"

func TestMaskEmail(t *testing.T) {
	tests := map[string]string{
		"":                     "",
		"a@x.com":              "•@•.com",
		"ab@xy.com":            "a•@x•.com",
		"john.doe@example.com": "jo••••••@ex•••••.com",
		"jane@mail.example.co": "ja••@ma••••••••••.co",
		"олександр@пошта.укр":  "ол•••••••@по•••.укр",
		"user@localhost":       "us••@lo•••••••",
		"nodomain":             "no••••••",
	}
	for email, want := range tests {
		if got := maskEmail(email); got != want {
			t.Errorf("maskEmail(%q) = %q, want %q", email, got, want)
		}
	}
}
