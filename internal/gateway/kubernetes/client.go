package kubernetes

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientSet *kubernetes.Clientset
}

func NewClient(cfg ClientConfig) (*Client, error) {
	var (
		restCfg *rest.Config
		err     error
	)
	if cfg.UseInClusterConfig {
		restCfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("loading in-cluster config: %w", err)
		}
	} else {
		restCfg, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{
				ExplicitPath: cfg.KubeConfigPathOrDefault(),
			},
			&clientcmd.ConfigOverrides{
				CurrentContext: cfg.KubeConfigContext,
			}).ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("loading kubeconfig: %w", err)
		}
	}

	clientSet, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes client: %w", err)
	}

	return &Client{
		clientSet: clientSet,
	}, nil
}
