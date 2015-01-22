package common

import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Download downlads the given url and saves it to
// name+prefix file of the current dir.
func Download(url string, name string, suffix string) (string, error) {
	// http call
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// read the content
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Opens the file
	file, err := os.Create(name + suffix)
	if err != nil {
		return "", err
	}

	// Writes and closes the file
	_, err = file.Write(data)
	if err != nil {
		return "", err
	}
	err = file.Close()
	if err != nil {
		return "", err
	}

	return name + suffix, nil
}

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
	res = strings.ToLower(res)

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
	missingPercentage := computeHavingPercentage(second, first)

	return missingPercentage + havingPercentage
}

// Compute the percentage of words matching : how many words are in second that exists in first.
func computeHavingPercentage(first string, second string) float32 {
	words := strings.Split(first, " ")
	found := 0
	for _, word := range words {
		if strings.Contains(second, word) {
			found++
		}
	}
	havingPercentage := float32(found) / float32(len(words))
	return havingPercentage
}
