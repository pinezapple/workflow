package core

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteK8SJob(ctx context.Context, jobName string, doDelete bool) (err error) {
	// manual delete
	if doDelete {
		err = ActualDeleteK8SJob(ctx, jobName)
		if err != nil {
			return err
		}
	}

	// log here
	return nil
}

func ActualDeleteK8SJob(ctx context.Context, jobName string) (err error) {
	mainConf := GetMainConfig()
	if mainConf.K8SConfig.DeleteJob {
		clientset := GetK8SClientSet()
		jobsClient := clientset.BatchV1().Jobs(mainConf.K8SConfig.K8SNameSpace)

		fg := metav1.DeletePropagationBackground
		deleteOptions := metav1.DeleteOptions{PropagationPolicy: &fg}
		//deleteOptions := metav1.DeleteOptions{}

		err = jobsClient.Delete(ctx, jobName, deleteOptions)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateK8SJob(ctx context.Context, k8sjob *batchv1.Job) (err error) {
	clientset := GetK8SClientSet()
	mainConf := GetMainConfig()

	jobsClient := clientset.BatchV1().Jobs(mainConf.K8SConfig.K8SNameSpace)

	// turn on result when add logger
	_, err = jobsClient.Create(ctx, k8sjob, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	//if !IfJobInPro(k8sjob.Name) {
	//	AddJobInPro(k8sjob.Name, ram, cpu, taskUUID)
	//}
	// log here
	return
}
