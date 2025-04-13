package k8s

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Info struct {
	DeploymentName string
	ContainerName  string
	Image          string
	Tag            string
}

// FetchDeploymentTags retrieves deployments for given names and extracts container image tags.
func FetchDeploymentInfo(clientset *kubernetes.Clientset, namespace string, services []string) ([]Info, error) {
	var results []Info

	for _, name := range services {
		deploy, err := clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			results = append(results, Info{
				DeploymentName: name,
				ContainerName:  "-",
				Image:          "Not Found",
				Tag:            "-",
			})
			continue
		}

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
			})
		}
	}

	return results, nil
}
