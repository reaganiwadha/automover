package main

import (
	"sync"
	"time"
)

var (
	wndVisible bool
	mu         sync.Mutex
)

func main() {
	if err := loadConfig(); err != nil {
		panic(err)
	}

	go runSystray()

	for {
		time.Sleep(100 * time.Millisecond)
	}
}
