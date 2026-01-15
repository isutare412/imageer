package kubernetes

import (
	"path/filepath"
	"time"

	"k8s.io/client-go/util/homedir"
)

type ClientConfig struct {
	UseInClusterConfig bool
	KubeConfigPath     string
	KubeConfigContext  string
}

func (c ClientConfig) KubeConfigPathOrDefault() string {
	if c.KubeConfigPath != "" {
		return c.KubeConfigPath
	}

	home := homedir.HomeDir()
	return filepath.Join(home, ".kube", "config")
}

type LeaderElectorConfig struct {
	LeaseName      string
	LeaseNamespace string
	LeaseDuration  time.Duration
	RenewDeadline  time.Duration
	RetryPeriod    time.Duration
}
