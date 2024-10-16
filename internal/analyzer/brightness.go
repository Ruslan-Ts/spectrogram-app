package analyzer

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/jpeg"
	"math"
)

// CreateBrightnessMap erzeugt eine Helligkeitskarte aus dem Eingabebild
func CreateBrightnessMap(img image.Image) (string, error) {
	// Schritt 1: Ermittlung der Abmessungen des Bildes
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Schritt 2: Erstellen Sie ein neues Bild für die Helligkeitskarte
	brightnessMap := image.NewGray(bounds)

	// Schritt 3: Gehen Sie jedes Pixel des Quellbildes durch
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Ermittelt die Farbe des Pixels
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			// Schritt 4: Berechnung der Helligkeit
			brightness := calculateBrightness(r, g, b)
			// Schritt 5: Einstellen des Helligkeitswertes im neuen Bild
			brightnessMap.Set(x, y, color.Gray{Y: uint8(brightness)})
		}
	}
	// Schritt 6: Kodierung der Leuchtdichtekarte in ein JPEG
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, brightnessMap, nil); err != nil {
		return "", err
	}
	// Schritt 7: Umwandlung in base64
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// calculateBrightness berechnet die Helligkeit der Pixel
func calculateBrightness(r, g, b uint32) uint8 {
	// Formel zur Berechnung der Helligkeit
	brightness := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	// Normalisieren auf den Bereich 0-255
	return uint8(math.Min(brightness/256, 255))
}
