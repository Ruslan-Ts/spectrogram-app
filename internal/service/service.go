package service

import (
	"image"
	"image/jpeg"
	"net/http"
	"spectrogram-app/internal/analyzer"
)

type AnalyzerService struct{}

func NewAnalyzerService() *AnalyzerService {
	return &AnalyzerService{}
}

type AnalysisResult struct {
	DominantColor string         `json:"dominant_color"`
	Spectrum      []SpectralData `json:"spectrum"`
	BrightnessMap string         `json:"brightness_map"` // base64 encoded image
}

type SpectralData struct {
	Wavelength float64 `json:"wavelength"`
	Intensity  float64 `json:"intensity"`
}

func (s *AnalyzerService) AnalyzeImage(url string) (*AnalysisResult, error) {
	// 1. Hochladen eines Bildes
	img, err := s.downloadImage(url)
	if err != nil {
		return nil, err
	}

	// 2. Analyse (Stopfen für Methoden, außer für die Helligkeitskarte)
	dominantColor := s.findDominantColor(img)
	spectrum := s.createSpectrum(img)
	brightnessMap, err := analyzer.CreateBrightnessMap(img)
	if err != nil {
		return nil, err
	}

	return &AnalysisResult{
		DominantColor: dominantColor,
		Spectrum:      spectrum,
		BrightnessMap: brightnessMap,
	}, nil
}

func (s *AnalyzerService) downloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (s *AnalyzerService) findDominantColor(img image.Image) string {
	// TODO: Durchführung der Suche nach der dominanten Farbe
	return "#000000"
}

func (s *AnalyzerService) createSpectrum(img image.Image) []SpectralData {
	// TODO: Die Erstellung eines Spektrogramms realisieren
	return []SpectralData{}
}

func (s *AnalyzerService) createBrightnessMap(img image.Image) string {
	// TODO: Implementierung der Erstellung einer Helligkeitskarte
	return ""
}
