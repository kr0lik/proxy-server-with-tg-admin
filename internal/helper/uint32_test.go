package helper_test

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
	"testing"
)

func TestStringToUint32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    uint32
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{3, "3"},
		{42, "42"},
		{300, "300"},
		{30, "30"},
		{1234567890, "1234567890"},
		{4294967295, "4294967295"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Uint32ToString(%d)", tt.input), func(t *testing.T) {
			t.Parallel()

			actual := helper.Uint32ToString(tt.input)
			if actual != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, actual)
			}
		})
	}
}

func TestUint32ToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected uint32
		hasError bool
	}{
		{"0", 0, false},
		{"1", 1, false},
		{"42", 42, false},
		{"1234567890", 1234567890, false},
		{"4294967295", 4294967295, false},
		{"-1", 0, true},
		{"123abc", 0, true},
		{"0000000001", 1, false},
		{"300", 300, false},
		{"30", 30, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("StringToUint32(%s)", tt.input), func(t *testing.T) {
			t.Parallel()

			actual, err := helper.StringToUint32(tt.input)
			if (err != nil) != tt.hasError {
				t.Errorf("Expected error: %v, got: %v", tt.hasError, err)
			}

			if !tt.hasError && actual != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, actual)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	t.Parallel()

	for i := range uint32(1000000) {
		s := helper.Uint32ToString(i)
		r, err := helper.StringToUint32(s)

		if err != nil || r != i {
			t.Fatalf("Round-trip failed for %d: got string %s, returned int %d", i, s, r)
		}
	}
}
