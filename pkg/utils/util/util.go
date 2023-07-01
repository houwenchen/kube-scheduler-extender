package util

import (
	"path/filepath"

	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func BuildKubeConfig(configPath string) (*restclient.Config, error) {
	if len(configPath) != 0 {
		clientConfig, err := clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			return BuildkubeConfigFromEnv()
		}
		return clientConfig, nil
	}

	return BuildkubeConfigFromEnv()
}

func BuildkubeConfigFromEnv() (*restclient.Config, error) {
	clientConfig, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err == nil {
		return clientConfig, nil
	}

	return restclient.InClusterConfig()
}
