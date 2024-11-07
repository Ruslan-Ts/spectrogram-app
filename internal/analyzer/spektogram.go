package analyzer

//not sure if this belongs in analyzer, please help
import (
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"math"
	"os"
	"time"

	"github.com/mjibson/go-dsp/fft"
)

// tested by coping contents into main, how can I solve this without doing that step?
//testing did not work

func spektogram() {

	//Load the image -> use function that is also used in color.go = maybe put it into a service?
	file, err := os.Open("cmd/tomie.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	defer file.Close()

	_, err = jpeg.DecodeConfig(file)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	file.Seek(0, io.SeekStart)

	img, format, err := image.Decode(file)
	log.Printf("Image format: %s\n", format)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	start := time.Now()

	//Turn the image into a greyscale version of it by setting every single pixel in a grey version
	greyImg := image.NewGray(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			originalColor := color.GrayModel.Convert(img.At(x, y))
			greyImg.Set(x, y, originalColor)
		}
	}

	//	r := new(big.Int)
	//	fmt.Println(r.Binomial(1000, 10))

	elapsed := time.Since(start)
	log.Printf("Greyscale took %d", elapsed.Milliseconds())

	//save the image in the folder that spektogram.go is -> move it somewhere else later
	outFile, err := os.Create("greyscale.jpg")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outFile.Close()
	err = jpeg.Encode(outFile, greyImg, nil)
	if err != nil {
		log.Fatalf("failed to save greyscale image: %v", err)
	}

	//set the 2d FTT
	bounds := greyImg.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	data := make([][]float64, height)
	for y := range data {
		data[y] = make([]float64, width)
		for x := range data[y] {
			data[y][x] = float64(greyImg.GrayAt(x, y).Y)
		}
	}

	//use the FTT on every single pixel
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

	//save FFT magnitude visualization in the folder that spektogram.go is -> move it somewhere else later
	outFFTFile, err := os.Create("fft_visualization.jpg")
	if err != nil {
		log.Fatalf("failed to create FFT output file: %v", err)
	}
	defer outFFTFile.Close()

	//the fft magintude is being set into a greyscale
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
