package dither

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/makeworld-the-better-one/dither/v2"
)

func raiseIfError(err error) {
	if err != nil {
		panic(err)
	}
}

/*
Returns an array of `color.Color` objects based on the input `theme` string. The array can then be
used as a color palette when dithering images. The `compression` input allows you to specify the
compression required:
  - compression = 0% - all 6 colours used
  - compression = 25% - 5 colours used
  - compression = 50% - 4 colours used
  - compression = 75% - 3 colours used
  - compression = 100% - 2 colours used
*/
func getColorPalette(theme string, compression string) (p []color.Color, err error) {
	// Build the palette
	switch theme {
	case "low-tech":
		p = []color.Color{
			color.RGBA{R: 30, G: 32, B: 40, A: 255},
			color.RGBA{R: 11, G: 21, B: 71, A: 255},
			color.RGBA{R: 57, G: 77, B: 174, A: 255},
			color.RGBA{R: 158, G: 168, B: 218, A: 255},
			color.RGBA{R: 187, G: 196, B: 230, A: 255},
			color.RGBA{R: 243, G: 244, B: 250, A: 255},
		}
	case "obsolete":
		p = []color.Color{
			color.RGBA{R: 9, G: 74, B: 58, A: 255},
			color.RGBA{R: 58, G: 136, B: 118, A: 255},
			color.RGBA{R: 101, G: 163, B: 148, A: 255},
			color.RGBA{R: 144, G: 189, B: 179, A: 255},
			color.RGBA{R: 169, G: 204, B: 195, A: 255},
			color.RGBA{R: 242, G: 247, B: 246, A: 255},
		}
	case "high-tech":
		p = []color.Color{
			color.RGBA{R: 86, G: 9, B: 6, A: 255},
			color.RGBA{R: 197, G: 49, B: 45, A: 255},
			color.RGBA{R: 228, G: 130, B: 124, A: 255},
			color.RGBA{R: 233, G: 155, B: 151, A: 255},
			color.RGBA{R: 242, G: 193, B: 190, A: 255},
			color.RGBA{R: 252, G: 241, B: 240, A: 255},
		}
	case "grayscale":
		p = []color.Color{
			color.RGBA{R: 25, G: 25, B: 25, A: 255},
			color.RGBA{R: 75, G: 75, B: 75, A: 255},
			color.RGBA{R: 125, G: 125, B: 125, A: 255},
			color.RGBA{R: 175, G: 175, B: 175, A: 255},
			color.RGBA{R: 225, G: 225, B: 225, A: 255},
			color.RGBA{R: 250, G: 250, B: 250, A: 255},
		}
	default:
		return p, fmt.Errorf("theme '%s' not supported", theme)
	}

	// Apply compression. Lower number of colours means higher compression
	switch compression {
	case "0%":
	case "25%":
		p = p[1:]
	case "50%":
		p = []color.Color{
			p[0],
			p[2],
			p[3],
			p[5],
		}
	case "75%":
		p = []color.Color{
			p[0],
			p[3],
			p[5],
		}
	case "100%":
		p = []color.Color{
			p[0],
			p[5],
		}
	default:
		return p, fmt.Errorf("compression '%s' not supported", compression)
	}
	return p, nil
}

func main() {
	inputFilename := "IMG_9916.png"

	filenameParts := strings.Split(inputFilename, ".")
	fileBasename := filenameParts[0]
	fileExtension := filenameParts[1]

	// open input file and decode based on image type
	fi, err := os.Open(inputFilename)
	raiseIfError(err)

	defer fi.Close()

	var image image.Image

	switch fileExtension {
	case "png":
		img, err := png.Decode(fi)
		raiseIfError(err)
		image = img

	case "jpg":
		img, err := jpeg.Decode(fi)
		raiseIfError(err)
		image = img
	default:
		panic(fmt.Sprintf("Image type '%s' not supported.", fileExtension))
	}

	colorPalettes := [4]string{"low-tech", "obsolete", "high-tech", "grayscale"}
	compressionValues := [5]string{"0%", "25%", "50%", "75%", "100%"}

	for i := 0; i < len(colorPalettes); i++ {
		for j := 0; j < len(compressionValues); j++ {
			// These are the colors we want in our output image
			palette, err := getColorPalette(colorPalettes[i], compressionValues[j])
			raiseIfError(err)

			fmt.Print(compressionValues[j], palette, "\n")
			// Create ditherer
			d := dither.NewDitherer(palette)
			d.Mapper = dither.Bayer(8, 8, 1.0)
			// d.Matrix = dither.FloydSteinberg
			// d.Matrix = dither.ErrorDiffusionStrength(dither.FloydSteinberg, 0.8)
			d.Serpentine = true

			// Dither the image, attempting to modify the existing image
			// If it can't then a dithered copy will be returned.
			imageDithered := d.DitherCopy(image)

			// create output filename
			outputFilename := fmt.Sprintf("%s_dithered_%s_%s.png", fileBasename, colorPalettes[i], compressionValues[j])

			// open output file
			fo, err := os.Create(outputFilename)
			raiseIfError(err)

			defer fo.Close()

			if err = png.Encode(fo, imageDithered); err != nil {
				panic(err)
			}
		}
	}
}
