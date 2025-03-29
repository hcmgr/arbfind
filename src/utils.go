package main

import (
	"os"

	"golang.org/x/exp/constraints"
)

/////////////////////////////////////////
// Miscellaneous helper functions
/////////////////////////////////////////

// Finds and returns the key with the max value
func findMaxKey[K comparable, V constraints.Ordered](m map[K]V) K {
	var maxKey K
	var maxValue V
	var first bool = true // Ensures we initialize with the first entry

	for key, value := range m {
		if first || value > maxValue {
			maxKey = key
			maxValue = value
			first = false
		}
	}

	return maxKey
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || !os.IsNotExist(err)
}
