package ordoe2e

import "testing"

func TestParityLabel(t *testing.T) {
	for value, want := range map[int]string{4: "even", 7: "odd"} {
		if got := ParityLabel(value); got != want {
			t.Errorf("ParityLabel(%d) = %q, want %q", value, got, want)
		}
	}
}
