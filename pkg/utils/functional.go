package utils

import "fmt"

type Predicate[T any] func(T) bool

// MapSlice produces a new slice from a slice of elements and a mapping function
func MapSlice[T any, U any](ts []T, mapper func(T) U) []U {
	us := make([]U, 0, len(ts))
	for _, t := range ts {
		u := mapper(t)
		us = append(us, u)
	}
	return us
}

// MapFailableSlice maps a slice using a mapping function which may fail.
// Returns upon all elements are mapped or terminates upon the first mapping error.
func MapFailableSlice[T any, U any](ts []T, mapper func(T) (U, error)) ([]U, error) {
	us := make([]U, 0, len(ts))
	for i, t := range ts {
		u, err := mapper(t)
		if err != nil {
			return nil, fmt.Errorf("slice elem %v: %w", i, err)
		}
		us = append(us, u)
	}
	return us, nil
}

// FilterSlice iterates over slice and applies predicate to all items
// Returns the slice of elements that satisfies the predicate
func FilterSlice[T any](ts []T, predicate func(T) bool) []T {
	var filteredTs []T
	for _, t := range ts {
		if predicate(t) {
			filteredTs = append(filteredTs, t)
		}
	}
	return filteredTs
}

// PartitionSlice iterates over the slice and applies predicate to all items.
// Splits the slice into a slice of members that satisfy it and a slice of members that don't satisfy it
func PartitionSlice[T any](ts []T, predicate func(T) bool) (accepeted []T, rejected []T) {
	for _, t := range ts {
		if predicate(t) {
			accepeted = append(accepeted, t)
		} else {
			rejected = append(rejected, t)
		}
	}
	return
}

// MapFilterSlice produces a new slice of elements which satisfy the predicate and consequently mapped
func MapFilterSlice[T any, U any](ts []T, predicate func(T) bool, mapper func(T) U) []U {
	us := make([]U, 0, len(ts))
	for _, t := range ts {
		if predicate(t) {
			us = append(us, mapper(t))
		}
	}
	return us
}
