package k8s

import (
	"context"
	"fmt"
	"io"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func GetLogs(clientset *kubernetes.Clientset, namespace, pod, path string) error {
	podLogOpts := &v1.PodLogOptions{
		Follow:     true,
		Timestamps: false,
	}

	req := clientset.CoreV1().Pods(namespace).GetLogs(pod, podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return fmt.Errorf("error getting log stream: %w", err)
	}
	defer podLogs.Close()

	logFilePath := fmt.Sprintf("%s/%s.log", path, pod)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		return fmt.Errorf("error opening log file: %w", err)
	}
	defer logFile.Close()

	fmt.Printf("logs for pod %s → %s\n", pod, logFilePath)

	_, err = io.Copy(logFile, podLogs)
	if err != nil {
		return fmt.Errorf("error streaming logs: %w", err)
	}

	return nil
}
