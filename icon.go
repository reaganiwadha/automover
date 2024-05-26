package main

import (
	"embed"
	"image"
	_ "image/png"
	"log"
)

//go:embed assets/*
var iconEmbed embed.FS

var (
	scienceTinyIconPng   image.Image
	scienceTinyIcon      []byte
	scienceTinyFlashIcon []byte
)

func init() {
	scienceTinyIcon, _ = iconEmbed.ReadFile("assets/science_tiny.ico")
	scienceTinyFlashIcon, _ = iconEmbed.ReadFile("assets/science_tiny_flash.ico")

	file, err := iconEmbed.Open("assets/science_tiny.png")
	if err != nil {
		log.Fatalf("Failed to open embedded file: %v", err)
	}
	defer file.Close()

	// Decode the image
	scienceTinyIconPng, _, err = image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}
}
