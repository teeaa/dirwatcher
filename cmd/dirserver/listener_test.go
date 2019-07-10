package main

import (
	"testing"
)

func TestCompare(t *testing.T) {
	comparisons := []string{"CREATE", "RENAME", "REMOVE"}
	matches := [][]string{
		{"CREATE", "CREATE|WRITE", "CHMOD|CREATE", "CREATE|STAT", "REMOVE|RENAME|CREATE"},
		{"RENAME", "RENAME|WRITE", "CHMOD|RENAME", "RENAME|STAT", "REMOVE|RENAME|CREATE"},
		{"REMOVE", "REMOVE|WRITE", "CHMOD|REMOVE", "REMOVE|STAT", "REMOVE|RENAME|CREATE"},
	}
	mismatches := []string{
		"", "WRITE", "CHMOD", "WRITE|CHMOD", "BANANA",
	}

	for i := 0; i < len(comparisons); i++ {
		for _, match := range matches[i] {
			if !compare(match, comparisons[i]) {
				t.Errorf("%s doesn't match %s", comparisons[i], match)
			}
		}
	}
	for _, mismatch := range mismatches {
		for i := 0; i < len(comparisons); i++ {
			if compare(mismatch, comparisons[i]) {
				t.Errorf("%s matches %s when it shouldn't", comparisons[i], mismatch)
			}
		}
	}

}
