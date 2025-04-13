package k8s

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func UpdateDeploymentTag(clientset *kubernetes.Clientset, namespace, deploymentName, tag string) error {
	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	deployment, err := deploymentsClient.Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment %s: %w", deploymentName, err)
	}

	for i, container := range deployment.Spec.Template.Spec.Containers {
		imageParts := strings.Split(container.Image, ":")
		imageName := imageParts[0]
		deployment.Spec.Template.Spec.Containers[i].Image = fmt.Sprintf("%s:%s", imageName, tag)
	}

	_, err = deploymentsClient.Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment %s: %w", deploymentName, err)
	}

	return nil
}
