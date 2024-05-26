package main

import (
	"bytes"
	"io"
	"os"
	"sync"
)

// MultiWriter is a struct that writes logs to both a file and memory buffer.
type MultiWriter struct {
	fileWriter io.Writer
	memBuffer  *bytes.Buffer
	mutex      sync.Mutex
}

// NewMultiWriter creates a new instance of MultiWriter.
func NewMultiWriter(filePath string) (*MultiWriter, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &MultiWriter{
		fileWriter: file,
		memBuffer:  bytes.NewBuffer(nil),
	}, nil
}

// Write writes the provided data to both file and memory buffer.
func (mw *MultiWriter) Write(p []byte) (n int, err error) {
	mw.mutex.Lock()
	defer mw.mutex.Unlock()

	// Write to file
	_, err = mw.fileWriter.Write(p)
	if err != nil {
		return 0, err
	}

	// Write to memory buffer
	n, err = mw.memBuffer.Write(p)
	if err != nil {
		return 0, err
	}

	return n, nil
}

// ReadLastLines reads the last n lines from the memory buffer.
func (mw *MultiWriter) ReadLastLines(n int) []byte {
	mw.mutex.Lock()
	defer mw.mutex.Unlock()

	buffer := mw.memBuffer.Bytes()
	lines := bytes.Split(buffer, []byte("\n"))

	start := len(lines) - n
	if start < 0 {
		start = 0
	}

	return bytes.Join(lines[start:], []byte("\n"))
}
