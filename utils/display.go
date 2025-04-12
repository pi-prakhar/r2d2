package utils

import (
	"fmt"
	"os"

	"github.com/pi-prakhar/r2d2/k8s"

	"github.com/olekukonko/tablewriter"
)

// ClearTerminal clears the terminal screen (works on Unix-like systems).
func ClearTerminal() {
	fmt.Print("\033[H\033[2J")
}

// PrintTable prints the list of TagInfo in a table format.
func PrintTable(data []k8s.TagInfo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Deployment", "Container", "Image", "Tag"})

	for _, item := range data {
		table.Append([]string{
			item.DeploymentName,
			item.ContainerName,
			item.Image,
			item.Tag,
		})
	}

	table.Render()
}
