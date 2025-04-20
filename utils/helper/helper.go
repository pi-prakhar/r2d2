package helper

import (
	"strings"

	"github.com/pi-prakhar/r2d2/constants"
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
