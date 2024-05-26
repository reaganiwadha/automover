package main

import (
	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
)

func saveConfigFromGui() {
	tmpCurrentCfg := config
	config = guiConfigCopy
	err := saveConfig()
	if err != nil {
		statusText = "Error saving config"
		config = tmpCurrentCfg
	} else {
		statusText = "Config saved"
	}
}

var (
	guiConfigCopy AutomoverConfig
	selectedIndex int
	statusText    string
)

func guiLoop() {
	listBoxStr := make([]string, 0)

	for _, folder := range guiConfigCopy.Watchlist {
		listBoxStr = append(listBoxStr, folder.WatchPath+" -> "+folder.WatchPattern+" -> "+folder.DestinationPath)
	}

	var pathStr *string
	var regexStr *string
	var destinationStr *string

	if selectedIndex >= 0 && selectedIndex < len(guiConfigCopy.Watchlist) {
		pathStr = &guiConfigCopy.Watchlist[selectedIndex].WatchPath
		regexStr = &guiConfigCopy.Watchlist[selectedIndex].WatchPattern
		destinationStr = &guiConfigCopy.Watchlist[selectedIndex].DestinationPath
	}

	g.SingleWindow().Layout(
		g.Label("Watchlist"),
		g.ListBox("ListBox1", listBoxStr).OnChange(func(ind int) {
			selectedIndex = ind
		}).Size(-1, 100),
		g.Row(
			g.Button("Add").OnClick(func() {
				guiConfigCopy.Watchlist = append(guiConfigCopy.Watchlist, AutomoverWatchedFolder{
					WatchPath:       "C:\\Users\\User\\Downloads",
					WatchPattern:    "*",
					DestinationPath: "C:\\Users\\User\\Downloads\\test",
				})
			}),
			g.Button("Delete").OnClick(func() {
				if selectedIndex >= 0 && selectedIndex < len(guiConfigCopy.Watchlist) {
					guiConfigCopy.Watchlist = append(guiConfigCopy.Watchlist[:selectedIndex], guiConfigCopy.Watchlist[selectedIndex+1:]...)
				}
			}),
		),
		g.Spacing(),
		g.Spacing(),
		g.Spacing(),
		g.Row(
			g.Label("Path:"),
			g.Button("Browse to select").OnClick(func() {
				path, err := dialog.Directory().Browse()
				if err != nil {
					return
				}

				*pathStr = path
			}),
		),
		g.InputText(pathStr),
		g.Spacing(),
		g.Label("Regex:"),
		g.InputText(regexStr),
		g.Spacing(),
		g.Row(
			g.Label("Destination:"),
			g.Button("Browse to select").OnClick(func() {
				path, err := dialog.Directory().Browse()
				if err != nil {
					return
				}

				*destinationStr = path
			}),
		),
		g.InputText(destinationStr),
		g.Spacing(),
		g.Spacing(),
		g.Spacing(),
		g.Row(
			g.Button("Save").OnClick(saveConfigFromGui),
		),

		g.Separator(),
		g.Label(statusText),
	)
}

func runGiu() {
	statusText = ""
	guiConfigCopy = config

	guiConfigCopy.Watchlist = append(guiConfigCopy.Watchlist, AutomoverWatchedFolder{
		WatchPath:       "C:\\Users\\User\\Downloads",
		WatchPattern:    "*",
		DestinationPath: "C:\\Users\\User\\Downloads\\test",
	})

	wnd := g.NewMasterWindow("Automover", 700, 400, g.MasterWindowFlagsNotResizable)
	wnd.Run(guiLoop)
	wndVisible = false
}
