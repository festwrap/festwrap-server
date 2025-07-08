package str

type Distance interface {
	compute(s1, s2 string) int
}

type LevenshteinDistance struct{}

// Computes the Levenshtein distance between two strings using definition
// Followed definition in https://en.wikipedia.org/wiki/Levenshtein_distance
func (d LevenshteinDistance) Compute(s1, s2 string) int {
	m := len(s1)
	n := len(s2)
	// Need to use an empty position to represent empty prefixes
	distances := make([][]int, m+1)
	for i := range distances {
		distances[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		distances[i][0] = i
	}
	for j := 0; j <= n; j++ {
		distances[0][j] = j
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			var cost int
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			distances[i][j] = min(
				min(distances[i-1][j]+1, distances[i][j-1]+1), // Deletion or insertion
				distances[i-1][j-1]+cost,                      // Substitution
			)
		}
	}
	return distances[m][n]
}
