package analyzer

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/jpeg"
	"math"
)

func CreateBrightnessMap(img image.Image) (string, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	brightnessMap := image.NewGray(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()

			brightness := calculateBrightness(r, g, b)

			brightnessMap.Set(x, y, color.Gray{Y: uint8(brightness)})
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, brightnessMap, nil); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func calculateBrightness(r, g, b uint32) uint8 {
	brightness := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	return uint8(math.Min(brightness/256, 255))
}
