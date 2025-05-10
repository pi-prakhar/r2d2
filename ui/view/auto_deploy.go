package view

import (
	"fmt"

	"github.com/pi-prakhar/r2d2/constants"
	ghservice "github.com/pi-prakhar/r2d2/internal/github"
	"github.com/pi-prakhar/r2d2/utils"
	"github.com/pi-prakhar/r2d2/utils/helper"
	"github.com/rivo/tview"
)

type UI struct {
	app              *tview.Application
	textView         *tview.TextView
	owner, repo, tag string
	deploymentNames  []string
}

// NewUI initializes the tview application and text view.
func NewUI(owner, repo, tag string, deploymentNames []string) *UI {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() { app.Draw() })

	return &UI{
		app:             app,
		textView:        textView,
		owner:           owner,
		repo:            repo,
		tag:             tag,
		deploymentNames: deploymentNames,
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

	fmt.Fprintf(ui.textView, "[yellow]🚀 Monitoring Tag:[-] [white]%s[-]\n", ui.tag)
	fmt.Fprintf(ui.textView, "[yellow]📦 Repository:[-] [white]%s/%s[-]\n", ui.owner, ui.repo)
	fmt.Fprintf(ui.textView, "[yellow]🧩 Deployments:[-] [white]%s[-]\n", helper.FormatList(ui.deploymentNames, "None"))
	fmt.Fprintf(ui.textView, "[yellow]%s[-]\n", sep)
	fmt.Fprintf(ui.textView, "\n[white]🔍 Checking ECR Push Status[-]\n")
}

// PrintStatus prints a processing or warning message based on current status.
func (ui *UI) PrintStatus(isWaiting bool) {
	ui.app.QueueUpdateDraw(func() {
		if isWaiting {
			fmt.Fprintf(ui.textView, "[yellow]🔄 No workflow runs found yet for tag '%s'. Waiting...[-]\n", ui.tag)
		} else {
			fmt.Fprintf(ui.textView, "[blue]🔄 ECR Push for tag '%s' is still in progress...[-]\n", ui.tag)
		}
	})
}

func (ui *UI) PrintJobStatuses(jobs []ghservice.WorkflowJobStatus) {
	ui.app.QueueUpdateDraw(func() {
		fmt.Fprintf(ui.textView, "\n[yellow]%s[-]\n", helper.SeparatorLine())
		fmt.Fprintf(ui.textView, "[yellow]%-30s  %-10s[-]\n", "Job Name", "Status")

		for _, job := range jobs {
			jobName := job.Name
			status := job.Status

			if job.Status == constants.JobStatusCompleted {
				status = job.Conclusion
			}

			colorTag := utils.GetJobStatusColorTag(status)

			fmt.Fprintf(ui.textView, "[white]%-30s  %s%-10s[-]\n", jobName, colorTag, status)
		}
	})
}
