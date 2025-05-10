package utils

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/pi-prakhar/r2d2/constants"
)

/*
GetDeploymentStatusColor - Updated with:
	🟢 Healthy/Complete statuses (green)
	🟠 Updating/Progressing statuses (orange)
	🟡 Scaling/Starting statuses (light yellow)
	🔴 Failed statuses (red)
	⚪️ Unknown status (gray)
GetPodStatusColor - Updated with:
	🟢 Healthy/Successful statuses (green)
	🟠 Transitioning/Starting status (yellow)
	🟡 Waiting/Pending statuses (light yellow)
	🔴 Failed/Terminated statuses (red)
	⚪️ Unknown status (gray)
GetPodPhaseColor - Updated with:
	🟡 Pending phase (light yellow)
	🟢 Running phase (green)
	✅ Succeeded phase (light green)
	🔴 Failed phase (red)
	⚪️ Unknown phase (gray)
*/

// GetDeploymentStatusColor returns appropriate color and display text for deployment status
func GetDeploymentStatusColor(status string) (tcell.Color, string) {
	displayText := status
	baseStatus := status

	// Handle statuses with reason (format "Status:Reason")
	if strings.Contains(status, ":") {
		parts := strings.SplitN(status, ":", 2)
		baseStatus = parts[0]
		reason := parts[1]
		displayText = fmt.Sprintf("%s: %s", baseStatus, reason)

		// Failures with reasons are always red
		if baseStatus == constants.DeploymentStatusFailed {
			return tcell.ColorRed, displayText
		}
	}

	switch baseStatus {
	// 🟢 Healthy/Complete
	case constants.DeploymentStatusComplete, constants.DeploymentStatusAvailable:
		return tcell.ColorGreen, displayText

	// 🟠 Updating/Progressing
	case constants.DeploymentStatusUpdating, constants.DeploymentStatusProgressing:
		return tcell.ColorOrange, displayText

	// 🟡 Scaling/Starting
	case constants.DeploymentStatusScaling, constants.DeploymentStatusStarting:
		return tcell.ColorLightYellow, displayText

	// 🔴 Failed
	case constants.DeploymentStatusFailed, constants.DeploymentStatusReplicaFailure:
		return tcell.ColorRed, displayText

	// ⚪️ Unknown
	case constants.DeploymentStatusUnknown:
		return tcell.ColorGray, displayText

	// Default if not recognized
	default:
		return tcell.ColorWhite, displayText
	}
}

// GetPodStatusColor returns appropriate color and display text for pod status
func GetPodStatusColor(status string) (tcell.Color, string) {
	// Check if the status is in the "Status:Reason" format
	if strings.Contains(status, ":") {
		parts := strings.SplitN(status, ":", 2)
		baseStatus := parts[0]
		reason := parts[1]

		// Handle different status types with reasons
		switch baseStatus {
		case constants.PodStatusPending, constants.PodStatusWaiting:
			return tcell.ColorLightYellow, baseStatus + ": " + reason
		case constants.PodStatusTerminated, constants.PodStatusFailed:
			return tcell.ColorRed, baseStatus + ": " + reason
		}
	}

	// 🟢 Healthy/Successful
	if status == constants.PodStatusRunning || status == constants.PodStatusSucceeded {
		return tcell.ColorGreen, status
	}

	// 🟠 Transitioning/Starting
	if status == constants.PodStatusStarting {
		return tcell.ColorYellow, status
	}

	// 🟡 Waiting/Pending
	if status == constants.PodStatusPending || status == constants.PodStatusWaiting {
		return tcell.ColorLightYellow, status
	}

	// 🔴 Failed/Terminated
	if status == constants.PodStatusTerminated || status == constants.PodStatusFailed {
		return tcell.ColorRed, status
	}

	// ⚪️ Unknown
	if status == constants.PodStatusUnknown {
		return tcell.ColorGray, status
	}

	// Default case
	return tcell.ColorWhite, status
}

// GetPodPhaseColor returns appropriate color and display text for pod phase
func GetPodPhaseColor(phase string) (tcell.Color, string) {
	switch phase {
	// 🟡 Pending - Pod has been accepted but not all containers are ready
	case constants.PodPhasePending:
		return tcell.ColorLightYellow, phase

	// 🟢 Running - Pod has been bound to a node and all containers are running
	case constants.PodPhaseRunning:
		return tcell.ColorGreen, phase

	// ✅ Succeeded - All containers have terminated successfully
	case constants.PodPhaseSucceeded:
		return tcell.ColorLightGreen, phase

	// 🔴 Failed - All containers have terminated and at least one failed
	case constants.PodPhaseFailed:
		return tcell.ColorRed, phase

	// ⚪️ Unknown - State of the pod could not be determined
	case constants.PodPhaseUnknown:
		return tcell.ColorGray, phase

	// Default case for any unexpected phases
	default:
		return tcell.ColorWhite, phase
	}
}

func GetJobStatusColorTag(status string) string {
	switch status {
	// 🟢 Success - Job completed successfully
	case constants.JobStatusSuccess:
		return "[green]"

	// 🔴 Failed - Job failed
	case constants.JobStatusFailed:
		return "[red]"

	// ⚪️ In Progress - Job is still running
	case constants.JobStatusInProgress:
		return "[gray]"

	// ⚪️ Default - Unknown or unhandled status
	default:
		return "[gray]"
	}
}
