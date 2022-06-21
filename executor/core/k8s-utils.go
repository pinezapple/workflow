package core

import (
	"strings"

	apps_v1 "k8s.io/api/apps/v1"
	batch_v1 "k8s.io/api/batch/v1"
	api_v1 "k8s.io/api/core/v1"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	rbac_v1beta1 "k8s.io/api/rbac/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetTaskIDFromPodID(podID string) (taskID string) {
	mainConfig := GetMainConfig()

	tmp := podID[len(mainConfig.K8SConfig.K8SNameSpace)+1:]
	tmpSlice := strings.Split(tmp, "-")

	for i := 0; i < len(tmpSlice)-1; i++ {
		taskID = taskID + tmpSlice[i] + "-"
	}
	return taskID[0 : len(taskID)-1]
}

func GetPodStatus(obj interface{}) (objStatus api_v1.PodStatus) {
	switch object := obj.(type) {
	case *api_v1.Pod:
		objStatus = object.Status
	default:
		objStatus = api_v1.PodStatus{}
	}

	return objStatus
}

func GetPodSpec(obj interface{}) (objStatus api_v1.PodSpec) {
	switch object := obj.(type) {
	case *api_v1.Pod:
		objStatus = object.Spec
	default:
		objStatus = api_v1.PodSpec{}
	}

	return objStatus
}

// GetObjectMetaData returns metadata of a given k8s object
func GetObjectMetaData(obj interface{}) (objectMeta meta_v1.ObjectMeta) {

	switch object := obj.(type) {
	case *apps_v1.Deployment:
		objectMeta = object.ObjectMeta
	case *api_v1.ReplicationController:
		objectMeta = object.ObjectMeta
	case *apps_v1.ReplicaSet:
		objectMeta = object.ObjectMeta
	case *apps_v1.DaemonSet:
		objectMeta = object.ObjectMeta
	case *api_v1.Service:
		objectMeta = object.ObjectMeta
	case *api_v1.Pod:
		objectMeta = object.ObjectMeta
	case *batch_v1.Job:
		objectMeta = object.ObjectMeta
	case *api_v1.PersistentVolume:
		objectMeta = object.ObjectMeta
	case *api_v1.Namespace:
		objectMeta = object.ObjectMeta
	case *api_v1.Secret:
		objectMeta = object.ObjectMeta
	case *ext_v1beta1.Ingress:
		objectMeta = object.ObjectMeta
	case *api_v1.Node:
		objectMeta = object.ObjectMeta
	case *rbac_v1beta1.ClusterRole:
		objectMeta = object.ObjectMeta
	case *api_v1.ServiceAccount:
		objectMeta = object.ObjectMeta
	case *api_v1.Event:
		objectMeta = object.ObjectMeta
	}
	return objectMeta
}

// GetObjectMetaData returns metadata of a given k8s object
func GetPodPhase(obj interface{}) (phase string) {

	switch object := obj.(type) {
	case *api_v1.Pod:
		phase = string(object.Status.Phase)
	default:
		phase = ""
	}
	return phase
}
