package table

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/pi-prakhar/r2d2/utils"
	"github.com/rivo/tview"
)

type WatchTagsApp struct {
	application *tview.Application
	table       *tview.Table
	namespace   string
}

// NewApp initializes a new tview application.
func NewWatchTagsApp(namespace string) *WatchTagsApp {
	app := &WatchTagsApp{
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

	app.table.SetCell(0, 2, tview.NewTableCell("Status").
		SetAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorYellow).
		SetAttributes(tcell.AttrBold).
		SetBackgroundColor(tcell.ColorBlack))

	flex := tview.NewFlex().AddItem(app.table, 0, 1, true)
	app.application.SetRoot(flex, true)

	return app
}

func (a *WatchTagsApp) Run() error {
	return a.application.Run()
}

func (a *WatchTagsApp) UpdateTable(data []k8s.Info) {
	for i := 1; i < a.table.GetRowCount(); i++ {
		a.table.RemoveRow(i)
	}

	for i, item := range data {
		a.table.SetCell(i+1, 0, tview.NewTableCell(item.DeploymentName).SetAlign(tview.AlignLeft))
		a.table.SetCell(i+1, 1, tview.NewTableCell(item.Tag).SetAlign(tview.AlignLeft))

		// 🎨 Color status cell using a helper function
		statusColor, displayStatus := utils.GetDeploymentStatusColor(item.Status)
		statusCell := tview.NewTableCell(displayStatus).
			SetAlign(tview.AlignLeft).
			SetTextColor(statusColor)

		a.table.SetCell(i+1, 2, statusCell)
	}
	a.application.Draw()
}

func (a *WatchTagsApp) Stop() {
	a.application.Stop()
}
