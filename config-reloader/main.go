package main

import (
	"config-reloader/config"
	"config-reloader/reloader"
	"k8s.io/client-go/kubernetes"
	"log"
)

func main(){
	clientConfig, err := config.LoadK8SConfig()
	if err != nil{
		log.Fatalf("couldnt load config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		log.Fatalf("couldnt create client: %v", err)
	}

	deploymentReloader := reloader.NewDeploymentReloader(clientset)
	deploymentReloader.RunInformer()

}