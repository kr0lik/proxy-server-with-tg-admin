package helper_test

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
	"testing"
)

func TestPasswordGenerate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		length int
	}{
		{1},
		{3},
		{5},
		{10},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("TestPasswordGenerate(%d)", tt.length), func(t *testing.T) {
			t.Parallel()

			got := helper.PasswordGenerate(tt.length)
			length := len([]rune(got))

			if length != tt.length {
				t.Errorf("PasswordGenerate(%d) = %s (%d)", tt.length, got, length)
			}
		})
	}
}
