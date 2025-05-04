package k8s

import (
	"context"
	"fmt"

	"github.com/pi-prakhar/r2d2/constants"
	"github.com/pi-prakhar/r2d2/utils/helper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodInfo struct {
	DeploymentName string
	PodName        string
	ContainerName  string
	Image          string
	Tag            string
	Status         string
	Phase          string
}

// FetchPodInfoForDeployments retrieves pods for given deployment names and extracts container image tags and status
func FetchPodInfoForDeployments(clientset *kubernetes.Clientset, namespace string, deploymentNames []string) ([]PodInfo, error) {
	var results []PodInfo
	ctx := context.Background()

	for _, deployName := range deploymentNames {
		deploy, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deployName, metav1.GetOptions{})
		if err != nil {
			results = append(results, PodInfo{
				DeploymentName: deployName,
				PodName:        "-",
				ContainerName:  "-",
				Image:          "Not Found",
				Tag:            "-",
				Status:         "Not Found",
				Phase:          "Not Found",
			})
			continue
		}

		selector := metav1.FormatLabelSelector(deploy.Spec.Selector)
		pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
			LabelSelector: selector,
		})

		if err != nil || len(pods.Items) == 0 {
			results = append(results, PodInfo{
				DeploymentName: deployName,
				PodName:        "-",
				ContainerName:  "-",
				Image:          "No Pods Found",
				Tag:            "-",
				Status:         "Not Found",
				Phase:          "Not Found",
			})
			continue
		}

		for _, pod := range pods.Items {
			for _, container := range pod.Spec.Containers {
				image, tag := helper.ParseImageTag(container.Image)

				results = append(results, PodInfo{
					DeploymentName: deployName,
					PodName:        pod.Name,
					ContainerName:  container.Name,
					Image:          image,
					Tag:            tag,
					Status:         determinePodStatus(&pod),
					Phase:          helper.MapPodPhaseToConstant(pod.Status.Phase),
				})
			}
		}
	}

	return results, nil
}

func determinePodStatus(pod *corev1.Pod) string {
	switch pod.Status.Phase {
	// 🟡 Pending
	case corev1.PodPending:
		for _, cond := range pod.Status.Conditions {
			if cond.Type == corev1.PodScheduled && cond.Status == corev1.ConditionFalse {
				return fmt.Sprintf("%s:%s", constants.PodStatusPending, cond.Reason)
			}
		}
		return constants.PodStatusPending

	// 🟢 Running
	case corev1.PodRunning:
		containersReady := false
		podReady := false

		for _, cond := range pod.Status.Conditions {
			if cond.Type == corev1.ContainersReady && cond.Status == corev1.ConditionTrue {
				containersReady = true
			}
			if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
				podReady = true
			}
		}

		// Container is running but not all ready
		if containersReady && podReady {
			return constants.PodStatusRunning
		}
		// 🟠 Running but not fully ready yet
		return constants.PodStatusStarting

	// ✅ Succeeded
	case corev1.PodSucceeded:
		return constants.PodStatusSucceeded

	// 🔴 Failed
	case corev1.PodFailed:
		return constants.PodStatusFailed

	// ⚪️ Unknown
	case corev1.PodUnknown:
		return constants.PodStatusUnknown
	}

	// 🔍 Container-level state checks
	for _, cs := range pod.Status.ContainerStatuses {
		if !cs.Ready {
			// 🟡 Waiting
			if cs.State.Waiting != nil {
				return fmt.Sprintf("%s:%s", constants.PodStatusWaiting, cs.State.Waiting.Reason)
			}
			// 🔴 Terminated with error
			if cs.State.Terminated != nil {
				return fmt.Sprintf("%s:%s", constants.PodStatusTerminated, cs.State.Terminated.Reason)
			}
		}
	}

	// Default case if we can't determine status
	return constants.PodStatusUnknown
}
