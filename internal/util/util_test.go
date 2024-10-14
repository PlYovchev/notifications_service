package util_test

import (
	"fmt"
	"testing"

	"github.com/plyovchev/sumup-assignment-notifications/internal/util"
)

type DevModeTestCase struct {
	input  string
	result bool
}

func TestIsDevMode(t *testing.T) {
	cases := []DevModeTestCase{
		{"dev", true},
		{"development ", true},
		{"test", false},
		{"stage", false},
		{"production", false},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s=%t", tc.input, tc.result), func(t *testing.T) {
			got := util.IsDevMode(tc.input)
			if tc.result != got {
				t.Errorf("Expected '%t', but got '%t'", tc.result, got)
			}
		})
	}
}
