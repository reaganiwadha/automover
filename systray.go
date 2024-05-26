package main

import (
	"fmt"
	"github.com/energye/systray"
	"github.com/energye/systray/icon"
	"time"
)

func runSystray() {
	systray.Run(onReady, onExit)
}

func onExit() {
}

func updateSystray() {
	systray.SetIcon(scienceTinyFlashIcon)
	go func() {
		time.Sleep(500 * time.Millisecond)
		systray.SetIcon(scienceTinyIcon)
	}()

	systray.SetTooltip(fmt.Sprintf("Last Moved: %s\nMove count: %v", lastMoved, counter))
}

func onReady() {
	systray.SetIcon(scienceTinyIcon)
	systray.SetTitle("Automover")
	systray.SetTooltip(fmt.Sprintf("Last Moved: %s\nMove count: %v", "none", counter))
	systray.SetOnClick(func(menu systray.IMenu) {
		mu.Lock()
		defer mu.Unlock()
		if !wndVisible {
			wndVisible = true
			go runGiu()
		}
	})
	systray.SetOnDClick(func(menu systray.IMenu) {
	})
	systray.SetOnRClick(func(menu systray.IMenu) {
		menu.ShowMenu()
	})
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	mQuit.SetIcon(icon.Data)
}
