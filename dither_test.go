package dither

import (
	"reflect"
	"strings"
	"testing"
)

/*
`getColourPalette` should return a valid array if input `theme` is one of the following:
  - "grayscale"
  - "high-tech"
  - "low-tech"
  - "obsolete"
*/
func TestGetColorPaletteReturnsArrayIfGoodTheme(t *testing.T) {
	validThemes := []string{"grayscale", "high-tech", "low-tech", "obsolete"}

	for i := 0; i < len(validThemes); i++ {
		palette, err := getColorPalette(validThemes[i], "0%")
		paletteType := reflect.TypeOf(palette)

		// Valid response is `palette` as array of `Color.color` and nil error
		if paletteType.Kind() == reflect.Array || err != nil {
			t.Fatalf("getColorPalette raised error: %s", err)
		}
	}
}

/*
`getColourPalette` should raise an error if input `theme` is not one of the following:
  - "grayscale"
  - "high-tech"
  - "low-tech"
  - "obsolete"
*/
func TestGetColorPaletteRaisesErrorIfBadTheme(t *testing.T) {
	invalidTheme := "my-bad-theme"

	_, err := getColorPalette(invalidTheme, "0%")

	// function should raise correct error
	if err == nil {
		t.Fatalf("getColorPalette did not raise error")
	}
	if !strings.Contains(err.Error(), "theme") {
		t.Fatalf("getColorPalette raised incorrect error: %s", err.Error())
	}
}

/*
`getColourPalette` should return a valid array if input `compression` is one of the following:
  - "0%"
  - "25%"
  - "50%"
  - "75%"
  - "100%"
*/
func TestGetColorPaletteReturnsArrayIfGoodCompression(t *testing.T) {
	validCompression := []string{"0%", "25%", "50%", "75%", "100%"}

	for i := 0; i < len(validCompression); i++ {
		palette, err := getColorPalette("grayscale", validCompression[i])
		paletteType := reflect.TypeOf(palette)

		// Valid response is `palette` as array of `Color.color` and nil error
		if paletteType.Kind() == reflect.Array || err != nil {
			t.Fatalf("getColorPalette raised error: %s", err)
		}
	}
}

/*
`getColourPalette` should raise an error if input `compression` is not one of the following:
  - "0%"
  - "25%"
  - "50%"
  - "75%"
  - "100%"
*/
func TestGetColorPaletteRaisesErrorIfBadCompression(t *testing.T) {
	invalidCompression := "1%"

	_, err := getColorPalette("grayscale", invalidCompression)

	// function should raise correct error
	if err == nil {
		t.Fatalf("getColorPalette did not raise error")
	}
	if !strings.Contains(err.Error(), "compression") {
		t.Fatalf("getColorPalette raised incorrect error: %s", err.Error())
	}
}
