package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return []string{"bash", "zsh", "fish", "powershell"}, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Long: `To load completions:

Bash:
  $ source <(r2d2 completion bash)

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ r2d2 completion zsh > "${fpath[1]}/_r2d2"

fish:
  $ r2d2 completion fish | source

  # To load completions for each session, execute once:
  $ r2d2 completion fish > ~/.config/fish/completions/r2d2.fish

PowerShell:
  PS> r2d2 completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> r2d2 completion powershell > r2d2.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logFile, err := os.OpenFile("completion.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Failed to open log file: %v", err)
		}
		logger := log.New(logFile, "completion: ", log.LstdFlags)
		logger.Printf("Running completion for args: %v", args)
		switch args[0] {
		case "bash":
			logger.Println("Generating bash completion")
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			logger.Println("Generating zsh completion")
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			logger.Println("Generating fish completion")
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			logger.Println("Generating powershell completion")
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			logger.Printf("Unknown completion shell: %v", args[0])
		}
		if logFile != nil {
			logFile.Close()
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
