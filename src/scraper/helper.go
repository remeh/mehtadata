package scraper

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/nfnt/resize"
	"model"
)

func FillDefaults(inputDirectory, filename string, gameinfo *model.Gameinfo) {
	gameinfo.Title = RemoveSpecialChars(RemoveExtension(filename))
	if !strings.HasSuffix(inputDirectory, "/") && !strings.HasPrefix(filename, "/") {
		filename = "/" + filename
	}
	gameinfo.Filepath = inputDirectory + filename
}

func RemoveExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) == 1 {
		return filename
	}
	return strings.Join(parts[0:len(parts)-1], "")
}

// ResizeImage uses the given data as an image and resize it
// with the max width given (computing the height).
// The filename is needed to know the filetype of the image.
func ResizeAndWrite(filename string, data []byte, writer io.Writer, maxWidth uint) error {
	var img image.Image
	var err error

	t := detectContentType(data)

	// if the automatic detect conte type failed, try with the extension
	if t == "" {
		if strings.HasSuffix(filename, "jpg") || strings.HasSuffix(filename, "jpeg") {
			t = "jpg"
		} else if strings.HasSuffix(filename, "png") {
			t = "png"
		} else {
			return fmt.Errorf("Unknown image format for : %s\n", filename)
		}
	}

	dontResize := (maxWidth == 0)

	// decodes the data
	if t == "jpg" {
		img, err = loadJpeg(data)
		if err != nil {
			log.Printf("Cant read '%s': %s\n", filename, err.Error())
			dontResize = true
		}
	} else if t == "png" {
		img, err = loadPng(data)
		if err != nil {
			log.Printf("Cant read '%s': %s\n", filename, err.Error())
			dontResize = true
		}
	}

	// resize if needed
	if dontResize {
		_, err = writer.Write(data)
	} else {
		if maxWidth > 0 {
			img = resize.Resize(maxWidth, 0, img, resize.Lanczos3)
		}

		if t == "jpg" {
			err = jpeg.Encode(writer, img, nil)
			if err != nil {
				log.Printf("Cant write '%s': %s\n", filename, err.Error())
				return err
			}
		} else if t == "png" {
			err = png.Encode(writer, img)
			if err != nil {
				log.Printf("Cant write '%s': %s\n", filename, err.Error())
				return err
			}
		}
	}

	return err
}

// Download downlads the given url and saves it to
// name+prefix file of the current dir.
// If resizeWidth > 0, some resizing will be done before saving the file.
func Download(url string, name string, suffix string, outputDirectory string, resizeWidth uint) (string, error) {
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
	file, err := os.Create(outputDirectory + name + suffix)
	if err != nil {
		return "", err
	}

	// Do some resizing if needed
	err = ResizeAndWrite(url, data, file, resizeWidth)
	if err != nil {
		return "", err
	}

	// Finally, close the file writer.
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
	// Replaces the "'s" with ""
	s = strings.Replace(s, "'s", "", -1)

	s = strings.Replace(s, " & ", " and ", -1)

	// List of special characters
	specialChars := ",.-_\"':[](){}!%#"

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
	havingPercentage := computePercentage(first, second, false)
	missingPercentage := computePercentage(second, first, true)

	return havingPercentage - (missingPercentage / 10.0)
}

// ----------------------

func detectContentType(data []byte) string {
	httpType := http.DetectContentType(data)
	if httpType == "image/png" {
		return "png"
	} else if httpType == "image/jpeg" {
		return "jpg"
	}
	return ""
}

func loadJpeg(data []byte) (image.Image, error) {
	return jpeg.Decode(bytes.NewReader(data))
}

func loadPng(data []byte) (image.Image, error) {
	return png.Decode(bytes.NewReader(data))
}

// Compute the percentage of words matching : how many words are in second that exists in first.
func computePercentage(first string, second string, missing bool) float32 {
	words := strings.Split(first, " ")
	found := 0
	for _, word := range words {
		if missing {
			// compute missing
			if !strings.Contains(second, word) {
				found++
			}
		} else {
			// compute having
			if strings.Contains(second, word) {
				found++
			}
		}
	}
	percentage := float32(found) / float32(len(words))
	return percentage
}
