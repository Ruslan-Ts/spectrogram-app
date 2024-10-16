package main

//später ändern zum Testen aber main erstmal
import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
)

func main() {
	img, err := urlToImage()
	if err != nil {
		fmt.Println("Error loading image:", err)
		return
	}

	getDominantColor(img)
}

func urlToImage() (image.Image, error) {
	//noch nicht dynamisch
	url := "https://upload.wikimedia.org/wikipedia/commons/b/be/Random_pyramids.jpg"

	//HTTP Request um das Bild zu laden, bei Error wird es gecatched
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching image:", err)
		return nil, fmt.Errorf("error fetching image: %w", err)
	}

	defer resp.Body.Close()
	//Bild decoden, Error wird gecatched
	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		fmt.Println("Error decoding JPEG image:", err)
		return nil, fmt.Errorf("error decoding JPEG image: %w", err)
	}

	return img, nil
}

func getDominantColor(img image.Image) string {
	//Map erstellen zum Farben speichern
	colorCount := make(map[color.Color]int)
	//Bounds vom Img festlegen
	bounds := img.Bounds()
	maxCount := 0
	var dominantColor color.Color

	//image Bounds in zwei Schleifen festhalten und dann jede Pixelfarbe abzählen
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixelColor := img.At(x, y)
			colorCount[pixelColor]++
		}
	}

	//meist genutzte Farbe zählen in dem man jede Farbe durchgeht und die mit den meisten Encountern in Var speichert
	for col, count := range colorCount {
		if count > maxCount {
			maxCount = count
			dominantColor = col
		}
	}

	//Konvertieren in Hex
	r, g, b, _ := dominantColor.RGBA()
	hexColor := fmt.Sprintf("#%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8))

	fmt.Println(hexColor)
	return hexColor
}
