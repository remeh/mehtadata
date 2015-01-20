package common

import (
	"regexp"
	"strings"
)

// Removes everything between (), [], {}
// FIXME Is it doable in one regexp ?
func ClearName(name string) string {
	st := " *\\([^)]*\\) *"
	reg, _ := regexp.Compile(st)
	res := reg.ReplaceAllString(name, "")

	st = " *\\[[^)]*\\] *"
	reg, _ = regexp.Compile(st)
	res = reg.ReplaceAllString(res, "")

	st = " *\\{[^)]*\\} *"
	reg, _ = regexp.Compile(st)
	res = reg.ReplaceAllString(res, "")

	res = RemoveSpecialChars(res)

	return res
}

// Remove the special chars from the given string
// and returns the result.
func RemoveSpecialChars(s string) string {
	// List of special characters
	specialChars := ".-_\"':[](){}!%#"

	// Copy the string
	result := s

	for i := 0; i < len(specialChars); i++ {
		result = strings.Replace(result, string(specialChars[i]), " ", -1)
	}

	return result
}

// Code re-used from https://github.com/remeh/go-subtitles
// to compute the percentage of matching words into the given strings.
func CompareFilename(first string, second string) float32 {

	// clears both filename
	first = ClearName(first)
	second = ClearName(second)

	// We now have two cleared file names.
	// The idea is to split by spaces and look how much of the "words"
	// of the subtitles filename we can find in the video filename
	havingPercentage := computeHavingPercentage(first, second)
	missingPercentage := computeMissingWordPercentage(first, second)

	// TODO FIXME This isn't correcly working as the one
	// not having enough words but the good ones are better match
	// than ones having enough words.
	// Example:
	// For: Castlevania Aria Of Sorrow
	// 'Castlevania' is a better match than
	// 'Castlevania Aria of'

	return havingPercentage - missingPercentage
}

// Compute the percentage of word matching from subtitle to video filename
func computeHavingPercentage(filename string, subtitleName string) float32 {
	words := strings.Split(subtitleName, " ")
	found := 0
	for _, word := range words {
		if strings.Contains(filename, word) {
			found++
		}
	}
	havingPercentage := float32(found) / float32(len(words))
	return havingPercentage
}

// Compute the percentage of word missing from the subtitle filename.
func computeMissingWordPercentage(filename string, subtitleName string) float32 {
	words := strings.Split(filename, " ")
	found := 0
	for _, word := range words {
		if strings.Contains(subtitleName, word) {
			found++
		}
	}
	missingPercentage := 1.0 - (float32(found) / float32(len(words)))
	return missingPercentage
}
