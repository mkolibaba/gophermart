package validation

import (
	"fmt"
	"testing"
)

func TestLuhn(t *testing.T) {
	cases := map[string]bool{
		"8532":     true,
		"12345674": true,
		"518191":   true,
		"6291911":  false,
		"62333":    false,
		"qwe":      false,
	}

	for num, want := range cases {
		t.Run(fmt.Sprintf("%d_should_be_%v", num, want), func(t *testing.T) {
			if Luhn(num) != want {
				t.Errorf("want %d to be %v", num, want)
			}
		})
	}
}
