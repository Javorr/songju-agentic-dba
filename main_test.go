package main

import "testing"

func TestMaskPII(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "email masking",
			input:    "SELECT * FROM users WHERE email = 'john@example.com'",
			expected: "SELECT * FROM users WHERE email = [REDACTED_EMAIL]",
		},
		{
			name:     "numeric ID in WHERE",
			input:    "SELECT * FROM orders WHERE user_id = 123",
			expected: "SELECT * FROM orders WHERE user_id = [REDACTED_ID]",
		},
		{
			name:     "multiple IDs with IN",
			input:    "SELECT * FROM orders WHERE id IN (1, 2, 3)",
			expected: "SELECT * FROM orders WHERE id IN [REDACTED_ID]",
		},
		{
			name:     "SSN format",
			input:    "SELECT * FROM records WHERE ssn = '123-45-6789'",
			expected: "SELECT * FROM records WHERE ssn = [REDACTED_SSN]",
		},
		{
			name:     "credit card",
			input:    "SELECT * FROM payments WHERE card = '4111-1111-1111-1111'",
			expected: "SELECT * FROM payments WHERE card = [REDACTED_CARD]",
		},
		{
			name:     "UUID",
			input:    "SELECT * FROM users WHERE id = '550e8400-e29b-41d4-a716-446655440000'",
			expected: "SELECT * FROM users WHERE id = [REDACTED_UUID]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskPII(tt.input)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func BenchmarkMaskPIITypes(b *testing.B) {
	benchmarks := []struct {
		name  string
		input string
	}{
		{"Email", "WHERE email = 'john@example.com'"},
		{"UUID", "WHERE id = '550e8400-e29b-41d4-a716-446655440000'"},
		{"FullQuery", "SELECT * FROM users WHERE email = 'john@example.com' AND id = 12345 AND ssn = '123-45-6789'"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = maskPII(bm.input)
			}
		})
	}
}
