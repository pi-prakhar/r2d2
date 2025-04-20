package table

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/pi-prakhar/r2d2/utils"
	"github.com/rivo/tview"
)

type WatchPodTagsApp struct {
	application *tview.Application
	table       *tview.Table
	namespace   string
}

// NewWatchPodTagsApp initializes a new tview application for pod-level information.
func NewWatchPodTagsApp(namespace string) *WatchPodTagsApp {
	app := &WatchPodTagsApp{
		application: tview.NewApplication(),
		table:       tview.NewTable(),
		namespace:   namespace,
	}

	app.table.SetTitle("Namespace: " + namespace + " (Pod Level)").
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(tcell.ColorPurple).
		SetBackgroundColor(tcell.ColorBlack).
		SetBorder(true)

	// Set up table headers
	app.table.SetCell(0, 0, tview.NewTableCell("Pod").
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

	app.table.SetCell(0, 3, tview.NewTableCell("Phase").
		SetAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorYellow).
		SetAttributes(tcell.AttrBold).
		SetBackgroundColor(tcell.ColorBlack))

	flex := tview.NewFlex().AddItem(app.table, 0, 1, true)
	app.application.SetRoot(flex, true)

	return app
}

func (a *WatchPodTagsApp) Run() error {
	return a.application.Run()
}

func (a *WatchPodTagsApp) UpdateTable(data []k8s.PodInfo) {
	// Clear existing rows (except header)
	for i := 1; i < a.table.GetRowCount(); i++ {
		a.table.RemoveRow(i)
	}

	for i, item := range data {
		a.table.SetCell(i+1, 0, tview.NewTableCell(item.PodName).SetAlign(tview.AlignLeft))
		a.table.SetCell(i+1, 1, tview.NewTableCell(item.Tag).SetAlign(tview.AlignLeft))

		// 🎨 Color status cell using a helper function
		statusColor, displayStatus := utils.GetPodStatusColor(item.Status)
		statusCell := tview.NewTableCell(displayStatus).
			SetAlign(tview.AlignLeft).
			SetTextColor(statusColor)

		a.table.SetCell(i+1, 2, statusCell)

		phaseColor, displayPhase := utils.GetPodPhaseColor(item.Phase)
		phaseCell := tview.NewTableCell(displayPhase).
			SetAlign(tview.AlignLeft).
			SetTextColor(phaseColor)

		a.table.SetCell(i+1, 3, phaseCell)
	}
	a.application.Draw()
}

func (a *WatchPodTagsApp) Stop() {
	a.application.Stop()
}
