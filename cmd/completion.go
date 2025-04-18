package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func getNamespaces(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	clientset, err := k8s.GetClientSet()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}

	var suggestions []string
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}
	for _, ns := range namespaces.Items {
		suggestions = append(suggestions, ns.Name)
	}

	return suggestions, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

// getDeployments returns a list of deployments in the specified namespace
func getDeployments(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	logFile, err := os.OpenFile("completion.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "completion: ", log.LstdFlags)

	var ns string
	if cmd.Flags().Changed("namespace") {
		ns, _ = cmd.Flags().GetString("namespace")
	} else {
		ns = namespace
	}
	logger.Printf("Resolved namespace: %s (from flags: %v)", ns, cmd.Flags().Changed("namespace"))

	if ns == "" {
		logger.Println("No namespace provided, returning empty suggestions")
		return nil, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}

	clientset, err := k8s.GetClientSet()
	if err != nil {
		logger.Printf("Error getting clientset: %v", err)
		return nil, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}

	// Split existing values
	parts := strings.Split(toComplete, ",")
	selected := make(map[string]bool)
	for _, part := range parts[:len(parts)-1] {
		if part != "" {
			selected[part] = true
		}
	}
	prefix := strings.Join(parts[:len(parts)-1], ",")
	if prefix != "" {
		prefix += ","
	}
	lastPart := parts[len(parts)-1]

	var suggestions []string
	deployments, err := clientset.AppsV1().Deployments(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		logger.Printf("Error listing deployments: %v", err)
		return nil, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}

	for _, deploy := range deployments.Items {
		name := deploy.Name
		if selected[name] {
			continue
		}
		if lastPart == "" || strings.HasPrefix(name, lastPart) {
			suggestions = append(suggestions, prefix+name)
		}
	}

	logger.Printf("Returning %d suggestions: %v", len(suggestions), suggestions)

	return suggestions, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
