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

var (
	namespace string
	services  []string
	tag       string
	frequency int
)

// getNamespaces returns a list of available namespaces
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
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   false,
		DisableDescriptions: false,
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace (required)")
	rootCmd.MarkPersistentFlagRequired("namespace")
	rootCmd.RegisterFlagCompletionFunc("namespace", getNamespaces)
}

// getDeployments returns a list of deployments in the specified namespace
func getDeployments(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Open log file for debugging
	logFile, err := os.OpenFile("completion.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "completion: ", log.LstdFlags)

	var ns string

	// Get namespace from command flags
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

	var suggestions []string
	deployments, err := clientset.AppsV1().Deployments(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		logger.Printf("Error listing deployments: %v", err)
		return nil, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}
	logger.Printf("Found %d deployments in namespace %s", len(deployments.Items), ns)

	parts := strings.Split(toComplete, ",")
	lastPart := parts[len(parts)-1]

	for _, deploy := range deployments.Items {
		if lastPart == "" || (len(deploy.Name) >= len(lastPart) && deploy.Name[:len(lastPart)] == lastPart) {
			fullName := strings.Join(append(parts[:len(parts)-1], deploy.Name), ",")
			suggestions = append(suggestions, fullName)
		}
	}
	logger.Printf("Returning %d deployment suggestions: %v", len(suggestions), suggestions)

	return suggestions, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
