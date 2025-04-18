package utils

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/pi-prakhar/r2d2/constants"
)

func GetDeploymentStatusColor(status string) (tcell.Color, string) {
	// Check if the status is in the "Failed: <reason>" format
	if strings.HasPrefix(status, constants.DeploymentStatusFailed+":") {
		reason := status[len(constants.DeploymentStatusFailed)+1:]             // Extract the reason part
		return tcell.ColorRed, constants.DeploymentStatusFailed + ":" + reason // Keep the reason in the status
	}

	// Handle other statuses
	switch status {
	case constants.DeploymentStatusReplicaFailure:
		return tcell.ColorRed, constants.DeploymentStatusReplicaFailure
	case constants.DeploymentStatusFailed:
		return tcell.ColorRed, constants.DeploymentStatusFailed
	case constants.DeploymentStatusComplete:
		return tcell.ColorGreen, constants.DeploymentStatusComplete
	case constants.DeploymentStatusAvailable:
		return tcell.ColorYellow, constants.DeploymentStatusAvailable
	case constants.DeploymentStatusScaling:
		return tcell.ColorLightYellow, constants.DeploymentStatusScaling
	case constants.DeploymentStatusStarting:
		return tcell.ColorLightYellow, constants.DeploymentStatusStarting
	case constants.DeploymentStatusUpdating:
		return tcell.ColorOrange, constants.DeploymentStatusUpdating
	case constants.DeploymentStatusProgressing:
		return tcell.ColorOrange, constants.DeploymentStatusProgressing
	case constants.DeploymentStatusUnknown:
		return tcell.ColorGray, constants.DeploymentStatusUnknown
	default:
		return tcell.ColorWhite, status // Default if not recognized
	}
}
