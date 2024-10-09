package main

import (
	"log"
	"spectrogram-app/internal/handler"
	"spectrogram-app/internal/service"
)

func main() {
	svc := service.NewAnalyzerService()
	h := handler.NewHandler(svc)

	if err := h.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
