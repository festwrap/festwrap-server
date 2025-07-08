package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevenshteinDistance(t *testing.T) {
	tests := map[string]struct {
		s1       string
		s2       string
		expected int
	}{
		"default case": {
			s1:       "house",
			s2:       "chair",
			expected: 5,
		},
		"empty strings": {
			s1:       "",
			s2:       "",
			expected: 0,
		},
		"first string empty": {
			s1:       "",
			s2:       "keyboard",
			expected: 8,
		},
		"second string empty": {
			s1:       "mouse",
			s2:       "",
			expected: 5,
		},
		"common prefix": {
			s1:       "mouse",
			s2:       "mousepad",
			expected: 3,
		},
		"common suffix": {
			s1:       "chair",
			s2:       "armchair",
			expected: 3,
		},
		"common middle": {
			s1:       "chair",
			s2:       "ai",
			expected: 3,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := LevenshteinDistance{}.Compute(test.s1, test.s2)

			assert.Equal(t, test.expected, actual)
		})
	}
}
