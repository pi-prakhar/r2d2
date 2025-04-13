package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/rivo/tview"
)

// NewApp initializes a new tview application.
func NewWatchTagsApp(namespace string) *App {
	app := &App{
		application: tview.NewApplication(),
		table:       tview.NewTable(),
		namespace:   namespace,
	}

	app.table.SetTitle("Namespace: " + namespace).
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(tcell.ColorPurple).
		SetBackgroundColor(tcell.ColorBlack).
		SetBorder(true)

	app.table.SetCell(0, 0, tview.NewTableCell("Deployment").
		SetAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorYellow).
		SetAttributes(tcell.AttrBold).
		SetBackgroundColor(tcell.ColorBlack))

	app.table.SetCell(0, 1, tview.NewTableCell("Tag").
		SetAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorYellow).
		SetAttributes(tcell.AttrBold).
		SetBackgroundColor(tcell.ColorBlack))

	flex := tview.NewFlex().AddItem(app.table, 0, 1, true)
	app.application.SetRoot(flex, true)

	return app
}

func (a *App) Run() error {
	return a.application.Run()
}

func (a *App) UpdateTable(data []k8s.TagInfo) {
	for i := 1; i < a.table.GetRowCount(); i++ {
		a.table.RemoveRow(i)
	}

	for i, item := range data {
		a.table.SetCell(i+1, 0, tview.NewTableCell(item.DeploymentName).SetAlign(tview.AlignLeft))
		a.table.SetCell(i+1, 1, tview.NewTableCell(item.Tag).SetAlign(tview.AlignLeft))
	}
	a.application.Draw()
}

func (a *App) Stop() {
	a.application.Stop()
}
