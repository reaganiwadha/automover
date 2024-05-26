package main

import (
	"fmt"
	"github.com/energye/systray"
	"github.com/energye/systray/icon"
)

func runSystray() {
	systray.Run(onReady, onExit)
}

func onExit() {
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome 超级棒")
	systray.SetOnClick(func(menu systray.IMenu) {
		fmt.Println("SetOnClick")
		mu.Lock()
		defer mu.Unlock()
		if !wndVisible {
			wndVisible = true
			go runGiu()
		}
	})
	systray.SetOnDClick(func(menu systray.IMenu) {
		fmt.Println("SetOnDClick")
	})
	systray.SetOnRClick(func(menu systray.IMenu) {
		menu.ShowMenu()
		fmt.Println("SetOnRClick")
	})
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	mQuit.SetIcon(icon.Data)
}
