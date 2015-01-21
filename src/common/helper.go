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

	return strings.Trim(res, " ")
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

	// Replace the double spaces that the previous commands could have been created.
	st := "\\s{2,}"
	reg, _ := regexp.Compile(st)
	result = reg.ReplaceAllString(result, " ")

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

	return havingPercentage
}

// Compute the percentage of words matching : how many words are in second that exists in first.
func computeHavingPercentage(first string, second string) float32 {
	words := strings.Split(second, " ")
	found := 0
	for _, word := range words {
		if strings.Contains(first, word) {
			found++
		}
	}
	havingPercentage := float32(found) / float32(len(words))
	return havingPercentage
}
