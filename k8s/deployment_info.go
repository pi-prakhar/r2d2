package k8s

import (
	"context"
	"strings"

	"github.com/pi-prakhar/r2d2/constants"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Info struct {
	DeploymentName string
	ContainerName  string
	Image          string
	Tag            string
	Status         string
}

// FetchDeploymentInfo retrieves deployments for given names and extracts container image tags and status.
func FetchDeploymentInfo(clientset *kubernetes.Clientset, namespace string, names []string) ([]Info, error) {
	var results []Info
	ctx := context.Background()

	for _, name := range names {
		deploy, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			results = append(results, Info{
				DeploymentName: name,
				ContainerName:  "-",
				Image:          "Not Found",
				Tag:            "-",
				Status:         "Not Found",
			})
			continue
		}

		// Determine deployment status based on conditions
		deploymentStatus := determineDeploymentStatus(deploy)

		// Process each container
		for _, container := range deploy.Spec.Template.Spec.Containers {
			imageParts := strings.Split(container.Image, ":")
			tag := "latest"
			if len(imageParts) > 1 {
				tag = imageParts[len(imageParts)-1]
			}

			results = append(results, Info{
				DeploymentName: name,
				ContainerName:  container.Name,
				Image:          container.Image,
				Tag:            tag,
				Status:         deploymentStatus,
			})
		}
	}

	return results, nil
}

func determineDeploymentStatus(deploy *appsv1.Deployment) string {
	replicas := int32(0)
	if deploy.Spec.Replicas != nil {
		replicas = *deploy.Spec.Replicas
	}

	// Extract conditions
	available := getDeploymentCondition(deploy, appsv1.DeploymentAvailable)
	progressing := getDeploymentCondition(deploy, appsv1.DeploymentProgressing)
	replicaFailure := getDeploymentCondition(deploy, appsv1.DeploymentReplicaFailure)

	// 🔴 1. Critical failure
	if replicaFailure != nil && replicaFailure.Status == corev1.ConditionTrue {
		return constants.DeploymentStatusReplicaFailure
	}
	if progressing != nil && progressing.Status == corev1.ConditionFalse {
		if progressing.Reason != "" {
			return constants.DeploymentStatusFailed + ":" + progressing.Reason
		}
		return constants.DeploymentStatusFailed
	}

	// 🟢 2. Fully rolled out
	if available != nil && available.Status == corev1.ConditionTrue {
		if progressing != nil &&
			progressing.Status == corev1.ConditionTrue &&
			progressing.Reason == "NewReplicaSetAvailable" {
			return constants.DeploymentStatusComplete
		}
		return constants.DeploymentStatusAvailable
	}

	// 🟠 3. Updating / Progressing
	if progressing != nil && progressing.Status == corev1.ConditionTrue {
		if deploy.Status.UpdatedReplicas < replicas {
			return constants.DeploymentStatusUpdating
		}
		return constants.DeploymentStatusProgressing
	}

	// 🟡 4. Scaling or Starting
	if deploy.Status.AvailableReplicas < replicas {
		if deploy.Status.AvailableReplicas > 0 {
			return constants.DeploymentStatusScaling
		}
		return constants.DeploymentStatusStarting
	}

	// ⚪️ 5. Unknown
	return constants.DeploymentStatusUnknown
}

func getDeploymentCondition(deploy *appsv1.Deployment, condType appsv1.DeploymentConditionType) *appsv1.DeploymentCondition {
	for i := range deploy.Status.Conditions {
		if deploy.Status.Conditions[i].Type == condType {
			return &deploy.Status.Conditions[i]
		}
	}
	return nil
}
