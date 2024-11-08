package etc

import (
	"sort"
	"strings"
)

// returns the keys from a map (in unpredictable order)
func Keys[K comparable, V any](m map[K]V) (results []K) {
	for k := range m {
		results = append(results, k)
	}
	return
}

// returns the keys from a map (in sorted order)
func SortedKeys[K ~string, V any](m map[K]V) (results []K) {
	for k := range m {
		results = append(results, k)
	}
	sort.Slice(results, func(i, j int) bool { return strings.Compare(string(results[i]), string(results[j])) == -1 })
	return
}
