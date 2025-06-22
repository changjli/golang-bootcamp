package combination

import (
	"reflect" // For deep comparison of slices
	"sort"    // To sort combinations for consistent comparison
	"testing"
)

func sortCombinations(combinations [][]int) {
	for _, combo := range combinations {
		sort.Ints(combo)
	}

	sort.Slice(combinations, func(i, j int) bool {
		lenI, lenJ := len(combinations[i]), len(combinations[j])
		minLen := lenI
		if lenJ < minLen {
			minLen = lenJ
		}
		for x := 0; x < minLen; x++ {
			if combinations[i][x] != combinations[j][x] {
				return combinations[i][x] < combinations[j][x]
			}
		}
		return lenI < lenJ
	})
}

type testCase struct {
	name     string
	n        int
	k        int
	expected [][]int
}

func TestCombinations(t *testing.T) {
	testCases := []testCase{
		// 1. Basic case (n=4, k=2)
		{
			name: "Basic Case n=4 k=2",
			n:    4,
			k:    2,
			expected: [][]int{
				{1, 2}, {1, 3}, {1, 4},
				{2, 3}, {2, 4},
				{3, 4},
			},
		},
		// 2. Single elements (k=1)
		{
			name: "Single Elements k=1 n=3",
			n:    3,
			k:    1,
			expected: [][]int{
				{1}, {2}, {3},
			},
		},
		// 3. Full combination (n=k)
		{
			name: "Full Combination n=k n=3 k=3",
			n:    3,
			k:    3,
			expected: [][]int{
				{1, 2, 3},
			},
		},
		// 4. Empty result (k=0)
		{
			name: "Empty Result k=0 n=5",
			n:    5,
			k:    0,
			expected: [][]int{
				{}, // A single empty combination as per problem conventions
			},
		},
		// Edge case: k > n should result in an empty set of combinations.
		{
			name:     "Empty Result k>n n=3 k=4",
			n:        3,
			k:        4,
			expected: [][]int{}, // No combinations possible
		},
		// 5. Invalid input (n=0)
		{
			name: "Invalid Input n=0 k=0",
			n:    0,
			k:    0,
			expected: [][]int{
				{},
			},
		},
		{
			name:     "Invalid Input n=0 k>0",
			n:        0,
			k:        1,
			expected: [][]int{}, // No combinations possible if n=0 and you need to pick elements.
		},
		// 6. Larger case (n=5, k=3)
		{
			name: "Larger Case n=5 k=3",
			n:    5,
			k:    3,
			expected: [][]int{
				{1, 2, 3}, {1, 2, 4}, {1, 2, 5},
				{1, 3, 4}, {1, 3, 5},
				{1, 4, 5},
				{2, 3, 4}, {2, 3, 5},
				{2, 4, 5},
				{3, 4, 5},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Combine(tc.n, tc.k)

			// Sort both actual and expected results for consistent comparison
			sortCombinations(actual)
			sortCombinations(tc.expected)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("For n=%d, k=%d:\nExpected: %v\nActual:   %v", tc.n, tc.k, tc.expected, actual)
			}
		})
	}
}

// --- Benchmarking ---
func BenchmarkCombine(b *testing.B) {
	n, k := 20, 10 // A moderately sized problem for combinations

	// Run the combine function b.N times
	for i := 0; i < b.N; i++ {
		Combine(n, k)
	}
}
