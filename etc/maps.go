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

// checks if a key is in a map
func MapContains[K comparable, V any, M ~map[K]V](key K, m M) bool {
	_, ok := m[key]
	return ok
}

//////////////////////////////////////////////////////////////////////
// Sets are just maps with no values

type Set[T comparable] map[T]struct{}

// returns a map that acts like a set
func MakeSet[T comparable](elements ...T) Set[T] {
	set := make(Set[T])
	for _, e := range elements {
		set[e] = struct{}{}
	}
	return set
}

func (s Set[T]) Includes(elem T) bool {
	_, ok := s[elem]
	return ok
}
