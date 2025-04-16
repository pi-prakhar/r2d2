package cmd

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

var (
	namespace string
	services  []string
	tag       string
	frequency int
	path      string
)

var rootCmd = &cobra.Command{
	Use:   "r2d2",
	Short: "Your Kubernetes protocol droid.",
	Long: `
            ___
          ,-'___'-.
        ,'  [(_)]  '.
       |_]||[][O]o[][|
     _ |_____________| _
    | []   _______   [] |
    | []   _______   [] |
   [| ||      _      || |]
    |_|| =   [=]     ||_|
    | || =   [|]     || |
    | ||      _      || |
    | ||||   (+)    (|| |
    | ||_____________|| |
    |_| \___________/ |_|
    / \      | |      / \
   /___\    /___\SSt /___\

Beep-boop! I'm R2-D2 — your loyal CLI droid, helping you monitor and 
manage Kubernetes deployments like a true Jedi.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
