package business

import (
	"context"

	"workflow/executor/core"
	"workflow/workflow-utils/model"

	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func RunK8SDaemon(parentCtx context.Context) (fn model.Daemon, err error) {
	lg := core.GetLogger()
	lg.Info("Starting k8s listen daemon")

	clientset := core.GetK8SClientSet()
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var resourceType = "pod"

	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return clientset.CoreV1().Pods(meta_v1.NamespaceAll).List(context.Background(), options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return clientset.CoreV1().Pods(meta_v1.NamespaceAll).Watch(context.Background(), options)
			},
		},
		&api_v1.Pod{},
		0, //Skip resync
		cache.Indexers{},
	)

	controller := NewControllerObj(clientset, queue, informer, resourceType)
	stop := make(chan struct{})
	go controller.Run(stop)

	fn = func() {
		<-parentCtx.Done()

		close(stop)
		lg.Info("Shutting down k8s listen daemon")
	}

	return fn, nil
}
