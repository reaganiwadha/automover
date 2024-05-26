package main

import (
	"github.com/sirupsen/logrus"
)

var (
	logFacet *MultiWriter
)

func init() {
	writer, err := NewMultiWriter("log.txt")
	if err != nil {
		logrus.Fatalf("Failed to create multiwriter: %v", err)
	}

	logFacet = writer

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		DisableQuote:  true,
	})
	logrus.SetOutput(writer)
}
