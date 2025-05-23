package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pi-prakhar/r2d2/version"
)

const githubLatestReleaseURL = "https://api.github.com/repos/pi-prakhar/r2d2/releases/latest"

func CheckForUpdate() {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(githubLatestReleaseURL)
	if err != nil {
		return // silently skip if network fails
	}
	defer resp.Body.Close()

	var data struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	latest := strings.TrimPrefix(data.TagName, "v")
	current := strings.TrimPrefix(version.Version, "v")
	if current == "dev" || latest != current && latest > current {
		fmt.Printf("\n\x1b[33m"+
			"╭──────────────────────────────────── Update Available ────────────────────────────────────╮\n"+
			"│ A new version of the CLI is available!                                                   │\n"+
			"│ Current version: %-72s│\n"+
			"│ Latest version:  %-72s│\n"+
			"│                                                                                          │\n"+
			"│ To update,                                                                               │\n"+
			"│ Run (recommended):                                                                       │\n"+
			"│   curl -sSL https://raw.githubusercontent.com/pi-prakhar/r2d2/master/install.sh | bash   │\n"+
			"│ Or visit:                                                                                │\n"+
			"│   https://github.com/pi-prakhar/r2d2/releases                                            │\n"+
			"╰──────────────────────────────────────────────────────────────────────────────────────────╯\n"+
			"\x1b[0m",
			current,
			latest,
		)
	}
}
