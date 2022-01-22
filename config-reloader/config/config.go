package config

import (
	"flag"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func LoadK8SConfig()(*rest.Config, error){
	if os.Getenv("RUNMODE") == "OUTCLUSTER"{
		return LoadOutClusterConfig()
	}
	return LoadInClusterConfig()
}

func LoadOutClusterConfig()(*rest.Config, error){
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	return clientcmd.BuildConfigFromFlags("", *kubeconfig)
}

func LoadInClusterConfig()(*rest.Config, error){
	return rest.InClusterConfig()
}