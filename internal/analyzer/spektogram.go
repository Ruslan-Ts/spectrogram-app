package analyzer

//not sure if this belongs in analyzer, please help
import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"

	"github.com/mjibson/go-dsp/fft"
)

// NEEDS TESTING
func spektogram() {
	// Step 1: Load and convert the image to greyscale
	file, err := os.Open("input.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	// Convert to greyscale
	greyImg := image.NewGray(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			originalColor := color.GrayModel.Convert(img.At(x, y))
			greyImg.Set(x, y, originalColor)
		}
	}

	// Save greyscale image
	outFile, err := os.Create("greyscale.jpg")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outFile.Close()
	err = jpeg.Encode(outFile, greyImg, nil)
	if err != nil {
		log.Fatalf("failed to save greyscale image: %v", err)
	}

	// Step 2: Apply 2D FFT
	bounds := greyImg.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	data := make([][]float64, height)
	for y := range data {
		data[y] = make([]float64, width)
		for x := range data[y] {
			data[y][x] = float64(greyImg.GrayAt(x, y).Y)
		}
	}

	// Perform FFT row by row and column by column
	for y := 0; y < height; y++ {
		row := make([]complex128, width)
		for x := 0; x < width; x++ {
			row[x] = complex(data[y][x], 0)
		}
		row = fft.FFT(row)
		for x := 0; x < width; x++ {
			data[y][x] = cmplxAbs(row[x])
		}
	}

	for x := 0; x < width; x++ {
		col := make([]complex128, height)
		for y := 0; y < height; y++ {
			col[y] = complex(data[y][x], 0)
		}
		col = fft.FFT(col)
		for y := 0; y < height; y++ {
			data[y][x] = cmplxAbs(col[y])
		}
	}

	// Save FFT magnitude visualization
	outFFTFile, err := os.Create("fft_visualization.jpg")
	if err != nil {
		log.Fatalf("failed to create FFT output file: %v", err)
	}
	defer outFFTFile.Close()

	// Visualize the FFT magnitude as a grayscale image
	fftImg := image.NewGray(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			greyValue := uint8(math.Min(255, data[y][x]/1000)) // Scale down for visualization
			fftImg.Set(x, y, color.Gray{Y: greyValue})
		}
	}
	err = jpeg.Encode(outFFTFile, fftImg, nil)
	if err != nil {
		log.Fatalf("failed to save FFT visualization: %v", err)
	}
}

// cmplxAbs returns the magnitude of a complex number
func cmplxAbs(c complex128) float64 {
	return math.Sqrt(real(c)*real(c) + imag(c)*imag(c))
}
