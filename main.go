package main

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	wndVisible bool
	mu         sync.Mutex
)

func main() {
	logrus.Info("Starting Automover")
	if err := loadConfig(); err != nil {
		panic(err)
	}

	go startWatcher()

	go runSystray()

	for {
		time.Sleep(100 * time.Millisecond)
	}
}
