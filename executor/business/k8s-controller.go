package business

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"workflow/executor/core"
	"workflow/workflow-utils/model"

	"go.temporal.io/sdk/client"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const maxRetries = 1

var podMap = make(map[string]bool)
var podMapLock sync.RWMutex

func SetPod(name string) {
	podMapLock.Lock()
	defer podMapLock.Unlock()
	podMap[name] = true
}

func GetPod(name string) (ok bool) {
	podMapLock.Lock()
	defer podMapLock.Unlock()
	_, ok = podMap[name]
	return
}

// Event indicate the informerEvent
type ControllerEvent struct {
	key          string
	eventType    string
	namespace    string
	resourceType string
}

type Controller struct {
	logger    *core.LogFormat
	clientset kubernetes.Interface
	queue     workqueue.RateLimitingInterface
	informer  cache.SharedIndexInformer
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	c.logger.Info("Starting k8s watcher")

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	c.logger.Info("Sync and ready")

	wait.Until(c.runWorker, time.Second, stopCh)
}

func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}

func (c *Controller) LastSyncResourceVersion() string {
	return c.informer.LastSyncResourceVersion()
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}
func (c *Controller) processNextItem() bool {
	newEvent, quit := c.queue.Get()

	if quit {
		return false
	}
	defer c.queue.Done(newEvent)
	err := c.processItem(newEvent.(ControllerEvent))

	if err == nil {
		// No error, reset the ratelimit counters
		c.queue.Forget(newEvent)
	} else if c.queue.NumRequeues(newEvent) < maxRetries {
		c.logger.Info(fmt.Sprintf("Error processing %s (will retry): %v", newEvent.(ControllerEvent).key, err))
		c.queue.AddRateLimited(newEvent)
	} else {
		// err != nil and too many retries
		c.logger.Errorf(fmt.Sprintf("Error processing %s (giving up): %v", newEvent.(ControllerEvent).key, err))
		c.queue.Forget(newEvent)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *Controller) OnHandleK8SDeleteNoti(newEvent ControllerEvent) error {
	return nil
}

func (c *Controller) ExecuteFailTaskWorkflow(taskID string) {
	e := GetExecutorTemporal()
	param := model.UpdateTaskFailParam{
		TaskID: taskID,
	}
	wo := client.StartWorkflowOptions{
		TaskQueue: "task-queue-name",
	}

	res, err := e.tempCli.ExecuteWorkflow(context.Background(), wo, model.FailTasktWfName, param)
	if err != nil {
		c.logger.Error(err.Error())
	}
	c.logger.Info("Start Fail Task Workflow with run ID " + res.GetRunID())
}

func (c *Controller) ExecuteDoneTaskWorkflow(taskID string, files []string, size []int64) {
	e := GetExecutorTemporal()
	param := model.UpdateTaskSuccessParam{
		TaskID:   taskID,
		Filename: files,
		Filesize: size,
	}
	wo := client.StartWorkflowOptions{
		TaskQueue: "task-queue-name",
	}

	res, err := e.tempCli.ExecuteWorkflow(context.Background(), wo, model.DoneTasktWfName, param)
	if err != nil {
		c.logger.Error(err.Error())
	}
	c.logger.Info("Start Done Task Workflow with run ID " + res.GetRunID())
}

func (c *Controller) HandleOldK8SJob(objMeta meta_v1.ObjectMeta, podPhase string) error {
	mainConf := core.GetMainConfig()
	atomic.AddInt32(&core.OldTaskCounter, 1)

	//taskUUID := objMeta.Labels["task-uuid"]
	if podPhase == "Succeeded" {
		successTask := objMeta.OwnerReferences[0].Name
		// get meta data like output file name and output file size
		var names []string
		var size []int64
		names, size, err := core.GetAllFileSizeInDirectory(mainConf.K8SConfig.OutputDirPrefix + "/" + successTask)
		if err != nil {
			c.logger.Error(err.Error())
		}
		core.DeleteK8SJob(context.Background(), successTask, true)

		// push signal to Temporal Workflow
		c.ExecuteDoneTaskWorkflow(successTask, names, size)
	}
	if podPhase == "Failed" {
		failTask := objMeta.OwnerReferences[0].Name
		core.DeleteK8SJob(context.Background(), failTask, true)

		// TODO: push signal to Temporal Workflow
		c.ExecuteFailTaskWorkflow(failTask)
	}

	if podPhase == "Running" {
		// TODO: add to running task, calculate to start listen client
		core.IncreaseJobCount()
	}
	return nil
}

func (c *Controller) OnHandleK8SUpdateNoti(newEvent ControllerEvent) error {
	obj, _, err := c.informer.GetIndexer().GetByKey(newEvent.key)
	if err != nil {
		return fmt.Errorf("Error fetching object with key %s from store: %v", newEvent.key, err)
	}
	podPhase := core.GetPodPhase(obj)
	objMeta := core.GetObjectMetaData(obj)
	if objMeta.OwnerReferences == nil {
		return nil
	}

	// If executor is in synchronize time
	if !core.GetSyncFlag() {
		return c.HandleOldK8SJob(objMeta, podPhase)
	}

	mainConf := core.GetMainConfig()
	//podStatus := core.GetPodStatus(obj)

	if podPhase == "Succeeded" {
		c.logger.Info("Job " + objMeta.Name + " is SUCCEEDED")

		// handler Success
		for i := 0; i < len(objMeta.OwnerReferences); i++ {
			if (objMeta.OwnerReferences[i].Kind == "Job") && !GetPod(objMeta.OwnerReferences[i].Name) {
				SetPod(objMeta.OwnerReferences[i].Name)

				var names []string
				var size []int64
				// get meta data like output file name and output file size
				names, size, err := core.GetAllFileSizeInDirectory(mainConf.K8SConfig.OutputDirPrefix + "/" + objMeta.OwnerReferences[i].Name)
				if err != nil {
					c.logger.Error(err.Error())
					return nil
				}
				core.DeleteK8SJob(context.Background(), objMeta.OwnerReferences[i].Name, false)

				// push signal to Temporal Workflow
				c.ExecuteDoneTaskWorkflow(objMeta.OwnerReferences[i].Name, names, size)
			}
		}
	} else if podPhase == "Failed" {
		c.logger.Info("Job " + objMeta.Name + " is FAILED")

		for i := 0; i < len(objMeta.OwnerReferences); i++ {
			if (objMeta.OwnerReferences[i].Kind == "Job") && !GetPod(objMeta.OwnerReferences[i].Name) {
				SetPod(objMeta.OwnerReferences[i].Name)

				_ = core.DeleteK8SJob(context.Background(), objMeta.OwnerReferences[i].Name, false)

				// push signal to Temporal Workflow
				c.ExecuteFailTaskWorkflow(objMeta.OwnerReferences[i].Name)
			}
		}
	}

	return nil
}

func (c *Controller) processItem(newEvent ControllerEvent) error {
	if newEvent.eventType == "update" || newEvent.eventType == "create" {
		return c.OnHandleK8SUpdateNoti(newEvent)
	} else if newEvent.eventType == "delete" {
		return c.OnHandleK8SDeleteNoti(newEvent)
	}
	return nil
}

func NewControllerObj(client kubernetes.Interface, queue workqueue.RateLimitingInterface, informer cache.SharedIndexInformer, resourceType string) (c *Controller) {
	var newEvent ControllerEvent
	var err error

	lg := core.GetLogger()
	config := core.GetMainConfig()

	// TODO: Change logger
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			newEvent.key, err = cache.MetaNamespaceKeyFunc(obj)
			newEvent.eventType = "create"
			newEvent.resourceType = resourceType
			newEvent.namespace = core.GetObjectMetaData(obj).Namespace

			if newEvent.namespace == config.K8SConfig.K8SNameSpace {
				lg.Dataf(fmt.Sprintf("Processing add to %v: %s", resourceType, newEvent.key))

				if err == nil {
					queue.Add(newEvent)
				}
			}
		},
		UpdateFunc: func(old, new interface{}) {
			newEvent.key, err = cache.MetaNamespaceKeyFunc(old)
			newEvent.eventType = "update"
			newEvent.resourceType = resourceType
			newEvent.namespace = core.GetObjectMetaData(old).Namespace

			if newEvent.namespace == config.K8SConfig.K8SNameSpace {
				lg.Dataf(fmt.Sprintf("Processing update to %v: %s", resourceType, newEvent.key))

				if err == nil {
					queue.Add(newEvent)
				}
			}

		},
		DeleteFunc: func(obj interface{}) {
			newEvent.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			//newEvent.key, err = cache.MetaNamespaceKeyFunc(obj)
			newEvent.eventType = "delete"
			newEvent.resourceType = resourceType
			newEvent.namespace = core.GetObjectMetaData(obj).Namespace

			if newEvent.namespace == config.K8SConfig.K8SNameSpace {
				lg.Dataf(fmt.Sprintf("Processing delete to %v: %s", resourceType, newEvent.key))

				if err == nil {
					queue.Add(newEvent)
				}
			}
		},
	})

	return &Controller{
		logger:    lg,
		clientset: client,
		queue:     queue,
		informer:  informer,
	}
}
