package view

import (
	"fmt"

	"github.com/pi-prakhar/r2d2/utils/helper"
	"github.com/rivo/tview"
)

type UI struct {
	app              *tview.Application
	textView         *tview.TextView
	owner, repo, tag string
}

// NewUI initializes the tview application and text view.
func NewUI(owner, repo, tag string) *UI {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() { app.Draw() })

	return &UI{
		app:      app,
		textView: textView,
		owner:    owner,
		repo:     repo,
		tag:      tag,
	}
}

// Run starts the tview application.
func (ui *UI) Run() error {
	return ui.app.SetRoot(ui.textView, true).Run()
}

// Stop stops the tview application.
func (ui *UI) Stop() {
	ui.app.Stop()
}

// Clear clears the text view and re-renders the header.
func (ui *UI) Clear() {
	ui.app.QueueUpdateDraw(func() {
		ui.textView.Clear()
		ui.renderHeader()
	})
}

// renderHeader prints the repository and tag info at the top.
func (ui *UI) renderHeader() {
	sep := helper.SeparatorLine()
	fmt.Fprintf(ui.textView, "[yellow]🚀 Monitoring Tag:[-] [white]%s\n", ui.tag)
	fmt.Fprintf(ui.textView, "[yellow]📦 Repository:[-] [white]%s/%s\n", ui.owner, ui.repo)
	fmt.Fprintf(ui.textView, "[yellow]─%s", sep)
	fmt.Fprintf(ui.textView, "\n[yellow]Checking ECR Push Status... [-]\n")
}

// PrintStatus prints a processing or warning message based on current status.
func (ui *UI) PrintStatus(isWaiting bool) {
	ui.app.QueueUpdateDraw(func() {
		if isWaiting {
			fmt.Fprintf(ui.textView, "\n[yellow]🔄 No workflow runs found yet for tag '%s'. Waiting...[-]\n", ui.tag)
		} else {
			fmt.Fprintf(ui.textView, "\n[blue]🔄 ECR Push for tag '%s' is still in progress...[-]\n", ui.tag)
		}
	})
}
