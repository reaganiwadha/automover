package main

import (
	g "github.com/AllenDang/giu"
	"github.com/sirupsen/logrus"
	"github.com/sqweek/dialog"
	"time"
)

var (
	guiConfigCopy AutomoverConfig
	selectedIndex int
)

var (
	sashPos1 float32 = 500
	sashPos2 float32 = 500
)

func guiLoop() {
	listBoxStr := make([]string, 0)

	for _, folder := range guiConfigCopy.Watchlist {
		listBoxStr = append(listBoxStr, folder.WatchPath+"\n -> "+folder.WatchPattern+" -> \n"+folder.DestinationPath)
	}

	var pathStr *string
	var regexStr *string
	var destinationStr *string

	if selectedIndex >= 0 && selectedIndex < len(guiConfigCopy.Watchlist) {
		pathStr = &guiConfigCopy.Watchlist[selectedIndex].WatchPath
		regexStr = &guiConfigCopy.Watchlist[selectedIndex].WatchPattern
		destinationStr = &guiConfigCopy.Watchlist[selectedIndex].DestinationPath
	} else {
		pathStr = new(string)
		regexStr = new(string)
		destinationStr = new(string)
	}

	log := string(logFacet.ReadLastLines(100))

	g.SingleWindow().Layout(
		g.SplitLayout(g.DirectionHorizontal, &sashPos1, g.Layout{
			g.SplitLayout(g.DirectionVertical, &sashPos2, g.Layout{
				g.Label("Watchlist"),
				g.ListBox("ListBox1", listBoxStr).OnChange(func(ind int) {
					selectedIndex = ind
				}).Size(-1, 400),
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
							selectedIndex = -1
						}
					}),
				),
			}, g.Layout{
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
				g.InputText(pathStr).Size(-1),
				g.Spacing(),
				g.Label("Regex:"),
				g.InputText(regexStr).Size(-1),
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
				g.InputText(destinationStr).Size(-1),
				g.Spacing(),
				g.Spacing(),
				g.Spacing(),
				g.Row(
					g.Button("Save").OnClick(func() {
						tmpCurrentCfg := config
						config = guiConfigCopy
						err := saveConfig()
						if err != nil {
							logrus.Info("Error saving config: ", err)
							config = tmpCurrentCfg
						} else {
							stopWatcher()
							startWatcher()
							logrus.Info("Config saved, restarting watcher...")
						}
					}).Size(-1, 25),
				),
			}),
		}, g.Layout{
			g.Label("Logs"),
			g.InputTextMultiline(&log).
				Size(-1, -1).
				Flags(g.InputTextFlagsReadOnly).
				AutoScrollToBottom(true),
		}),
	)
}

func runGiu() {
	guiConfigCopy = config

	wnd := g.NewMasterWindow("Automover", 1000, 700, 0)

	done := make(chan bool)

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				g.Update()
			}
		}
	}()

	wnd.Run(guiLoop)
	close(done)

	wndVisible = false
}
