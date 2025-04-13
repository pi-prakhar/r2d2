package utils

import (
	"github.com/pi-prakhar/r2d2/k8s"
)

type App interface {
	Run() error
	UpdateTable(data []k8s.Info)
	Stop()
}
