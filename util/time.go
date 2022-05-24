package util

import (
	"log"
	"time"
)

// PrintTimeElapsed print the execution time need by caller block code
func PrintTimeElapsed(what string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n", what, time.Since(start))
	}
}
