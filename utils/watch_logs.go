package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/rivo/tview"
)

type WatchLogsApp struct {
	application *tview.Application
	table       *tview.Table
	namespace   string
}

// NewApp initializes a new tview application.
func NewWatchLogsApp(namespace string) *WatchLogsApp {
	app := &WatchLogsApp{
		application: tview.NewApplication(),
		table:       tview.NewTable(),
		namespace:   namespace,
	}

	app.table.SetTitle("Namespace: " + namespace).
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(tcell.ColorPurple).
		SetBackgroundColor(tcell.ColorBlack).
		SetBorder(true)

	app.table.SetCell(0, 0, tview.NewTableCell("Pod").
		SetAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorYellow).
		SetAttributes(tcell.AttrBold).
		SetBackgroundColor(tcell.ColorBlack))

	app.table.SetCell(0, 1, tview.NewTableCell("File").
		SetAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorYellow).
		SetAttributes(tcell.AttrBold).
		SetBackgroundColor(tcell.ColorBlack))

	flex := tview.NewFlex().AddItem(app.table, 0, 1, true)
	app.application.SetRoot(flex, true)

	return app
}

func (a *WatchLogsApp) Run() error {
	return a.application.Run()
}

func (a *WatchLogsApp) UpdateTable(data []k8s.Info) {
	for i := 1; i < a.table.GetRowCount(); i++ {
		a.table.RemoveRow(i)
	}

	for i, item := range data {
		a.table.SetCell(i+1, 0, tview.NewTableCell(item.PodName).SetAlign(tview.AlignLeft))
		a.table.SetCell(i+1, 1, tview.NewTableCell(item.Path).SetAlign(tview.AlignLeft))
	}
	a.application.Draw()
}

func (a *WatchLogsApp) Stop() {
	a.application.Stop()
}
