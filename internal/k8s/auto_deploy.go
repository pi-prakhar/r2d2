package k8s

import (
	"github.com/pi-prakhar/r2d2/config"
	"k8s.io/client-go/kubernetes"
)

func UpdateK8sDeploymentTag(k8sClient *kubernetes.Clientset, config config.AutoDeployConfig) (bool, []error) {
	successCount := 0
	errors := []error{}
	for _, deployment := range config.Deployments {
		err := UpdateDeploymentTag(k8sClient, config.Namespace, deployment, config.Tag)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		successCount++
	}
	allSuccessful := successCount == len(config.Deployments)
	return allSuccessful, errors
}
