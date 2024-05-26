package main

import (
	"fmt"
	"sync"
	"time"

	g "github.com/AllenDang/giu"
	"github.com/energye/systray"
	"github.com/energye/systray/icon"
)

var (
	wndVisible bool
	mu         sync.Mutex
)

func onClickMe() {
	fmt.Println("Hello world!")
}

func onImSoCute() {
	fmt.Println("I'm sooooooo cute!!")
}

func loop() {
	g.SingleWindow().Layout(
		g.Label("Hello world from giu"),
		g.ListBox("ListBox1", []string{"item1", "item2", "item3"}).OnChange(func(ind int) {
			fmt.Println("Selected item:", ind)
		}),
		g.Table().Rows(buildRows()...),
		g.Row(
			g.Button("Click Me").OnClick(onClickMe),
			g.Button("I'm so cute").OnClick(onImSoCute),
		),
	)
}

var (
	names []string
)

func buildRows() []*g.TableRowWidget {
	rows := make([]*g.TableRowWidget, len(names))

	for i := range rows {
		rows[i] = g.TableRow(
			g.Label(fmt.Sprintf("%d", i)),
			g.Label(names[i]),
		)
	}

	//rows[0].BgColor(&(color.RGBA{200, 100, 100, 255}))

	return rows
}

func runGiu() {
	wnd := g.NewMasterWindow("Hello world", 400, 200, g.MasterWindowFlagsNotResizable)
	wnd.Run(loop)
	wndVisible = false
}

func runSystray() {
	systray.Run(onReady, onExit)
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

func onExit() {
	// Clean up here
}

func main() {
	go runSystray()

	names = make([]string, 10000)
	for i := range names {
		names[i] = fmt.Sprintf("Huge list name demo 范例 %d", i)
	}

	// Keep the main thread alive
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
