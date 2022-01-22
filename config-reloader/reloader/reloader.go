package reloader

import (
	"config-reloader/model"
	"context"
	"encoding/json"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"log"
)

type  DeploymentReloader struct {
	clientset *kubernetes.Clientset
}

func NewDeploymentReloader(clientset *kubernetes.Clientset) *DeploymentReloader {
	return &DeploymentReloader{clientset: clientset}
}

func (r *DeploymentReloader) RunInformer(){
	informerFactory := informers.NewSharedInformerFactory(r.clientset, 0)
	cmInformer := informerFactory.Core().V1().ConfigMaps().Informer()
	cmInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: r.onConfigMapUpdate,
	})
	stopChan := make(chan struct{})
	cmInformer.Run(stopChan)
}

func (r *DeploymentReloader) onConfigMapUpdate(old, new interface{}){
	newConfigMap, ok := new.(*v1.ConfigMap)
	if !ok{
		log.Fatalf("unrecognised format of v1ConfigMap")
	}
	if !(newConfigMap.Name =="environment" && newConfigMap.Namespace == "hasura"){
		return
	}

	envVarMap := getDataFromConfigMap(newConfigMap)
	for deploymentName, envs := range envVarMap{
		err := r.patch(deploymentName,envs)
		if err != nil{
			log.Printf("could not patch deployment %s. err: %v",deploymentName,  err)
			continue
		}
		log.Printf("patched deployment: %s", deploymentName)
	}

}

func(r *DeploymentReloader)patch(deploymentName model.DeploymentName, envs []model.EnvVar)error{
	log.Printf("attempting patch of deployment %s", deploymentName)

	result, getErr := r.clientset.AppsV1().Deployments("hasura").Get(context.TODO(), string(deploymentName), v12.GetOptions{})
	if getErr != nil {
		return getErr
	}
	envsForContainer := []v1.EnvVar{}
	for i:= range envs{
		envsForContainer = append(envsForContainer, v1.EnvVar{
			Name:      envs[i].Name,
			Value:     envs[i].Value,
			ValueFrom: nil,
		})
	}
	result.Spec.Template.Spec.Containers[0].Env = envsForContainer
	_, updateErr := r.clientset.AppsV1().Deployments("hasura").Update(context.TODO(), result, v12.UpdateOptions{})

	return updateErr
}




func getDataFromConfigMap(configMap *v1.ConfigMap) model.ConfigMapData{
	data := configMap.Data
	ParsedConfigmapData := model.ConfigMapData{}

	for deployment, envVars := range data{
		envVarList := []model.EnvVar{}

		err := json.Unmarshal( []byte(envVars), &envVarList)
		if err != nil{
			log.Printf("invalid json in deployment %s. error: %v", deployment, err)
			continue
		}
		ParsedConfigmapData[model.DeploymentName(deployment)]= envVarList
	}

	return ParsedConfigmapData
}