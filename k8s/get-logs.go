package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

func GetLogs(clientset *kubernetes.Clientset, namespace string, service string, path string) (string, error) {
	//TODO add options for tail and since time
	podLogOpts := &v1.PodLogOptions{}

	req := clientset.CoreV1().Pods(namespace).GetLogs(service, podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	// Open the file to append logs
	logFile, err := os.OpenFile(fmt.Sprintf("%s/%s.log", path, service), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer logFile.Close()

	buf := new(bytes.Buffer)
	writer := io.MultiWriter(buf, logFile)

	_, err = io.Copy(writer, podLogs)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s.log", path, service), nil
}
