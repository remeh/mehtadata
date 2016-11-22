package scraper

import (
	"testing"
)

func TestCompareFilename(t *testing.T) {
	original := "Pocman Forever"
	found := "Pocman Forever"

	res := CompareFilename(original, found)
	if res != 1.0 {
		t.Errorf("Should be a perfect match. Had %f instead of %f.", res, 1.0)
	}

	original = "Pocman Forever"
	found = "Pocman 2"
	res = CompareFilename(original, found)
	if res != 0.45 {
		t.Errorf("Should be a half match. Had %f instead of %f.", res, 0.5)
	}

	original = "Pocman Forever Machin"
	found = "Pocman 2"
	res = CompareFilename(original, found)
	if res != 0.283333333 {
		t.Errorf("Should be bad match. Had %f instead of %f.", res, 0.333333)
	}

	original = "Hache Dargent II"
	found = "Hache Dargent"
	res = CompareFilename(original, found)
	if res != 0.6666666666 {
		t.Errorf("Should be an approximate match. Had %f instead of %f.", res, 0.666666)
	}

	original = "Hache Dargent II"
	found = "Hache Dargent II"
	res = CompareFilename(original, found)
	if res != 1.0 {
		t.Errorf("Should be a perfect match. Had %f instead of %f.", res, 1.0)
	}
}
