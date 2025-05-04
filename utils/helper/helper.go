package helper

import (
	"fmt"
	"os"
	"strings"

	"github.com/pi-prakhar/r2d2/config"
	"github.com/pi-prakhar/r2d2/constants"
	"golang.org/x/term"
	corev1 "k8s.io/api/core/v1"
)

func MapPodPhaseToConstant(phase corev1.PodPhase) string {
	switch phase {
	case corev1.PodPending:
		return constants.PodPhasePending
	case corev1.PodRunning:
		return constants.PodPhaseRunning
	case corev1.PodSucceeded:
		return constants.PodPhaseSucceeded
	case corev1.PodFailed:
		return constants.PodPhaseFailed
	case corev1.PodUnknown:
		return constants.PodPhaseUnknown
	default:
		return constants.PodPhaseUnknown
	}
}

func ParseImageTag(image string) (string, string) {
	parts := strings.Split(image, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return image, "latest"
}

func GetGitHubTokenFromEnv() (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return "", fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}
	return token, nil
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width == 0 {
		// Default width if unable to get terminal size
		return 80
	}
	return width
}

func SeparatorLine() string {
	width := getTerminalWidth()
	return strings.Repeat("─", width)
}

func BuildFinalSummary(ecrSuccess, deploySuccess bool, config config.AutoDeployConfig, errors []error) string {
	var sb strings.Builder
	// sep := helper.SeparatorLine()
	sb.WriteString(fmt.Sprintf("🚀 [Tag] %s\n", config.Tag))
	sb.WriteString(fmt.Sprintf("📦 [Repository] %s/%s\n", config.Owner, config.Repo))
	sb.WriteString("🏁 Final Summary:\n\n")

	if !ecrSuccess {
		sb.WriteString("❌ [ECR Push] - Failed\n")
	} else {
		sb.WriteString("✅ [ECR Push] - Success\n")
		if deploySuccess {
			sb.WriteString("✅ [Deployment] - All deployments updated successfully!\n")
		} else {
			sb.WriteString("⚠️ [Deployment] - Some deployments failed.\n")
		}
	}

	if len(errors) > 0 {
		sb.WriteString("\nErrors:\n")
		for _, err := range errors {
			if err != nil {
				sb.WriteString(fmt.Sprintf("- %s\n", err.Error()))
			}
		}
	}

	sb.WriteString(fmt.Sprintf("\n🔗 GitHub Actions: https://github.com/%s/%s/actions\n", config.Owner, config.Repo))

	return sb.String()
}
