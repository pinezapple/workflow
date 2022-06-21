package business

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"workflow/executor/api/eddamanager"
	"workflow/executor/core"
	"workflow/executor/model"

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
	podMap[name] = true
	podMapLock.Unlock()
}

func GetPod(name string) (ok bool) {
	podMapLock.Lock()
	_, ok = podMap[name]
	podMapLock.Unlock()
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
	logger    *core.LogFormat //TODO: logger
	clientset kubernetes.Interface
	queue     workqueue.RateLimitingInterface
	informer  cache.SharedIndexInformer
	//TODO: handler
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
	taskID := core.GetTaskIDFromPodID(newEvent.key)
	ram, cpu, taskUUID, ok := core.DeleteJobInPro(taskID)

	if ok {
		atomic.AddInt64(&core.CPULeft, cpu)
		atomic.AddInt64(&core.RAMLeft, ram)

		// push fail update status
		req := &model.UpdateStatusCheck{
			Success: false,
			TaskID:  taskID,
		}

		err := PushUpdateStatusToKafka(req) // push to kafka
		if err != nil {
			c.logger.Error(err.Error())
			return nil
		}

		_, err = eddamanager.UpdateTaskLogState(taskUUID)
		if err != nil {
			c.logger.Error(err.Error())
			return nil
		}
	}

	return nil
}

func (c *Controller) HandleOldK8SJob(objMeta meta_v1.ObjectMeta, podPhase string) error {
	mainConf := core.GetMainConfig()
	atomic.AddInt32(&core.OldTaskCounter, 1)

	if podPhase == "Succeeded" {
		successTask := objMeta.OwnerReferences[0].Name
		taskUUID := objMeta.Labels["task-uuid"]
		//AddToSuccessTask(objMeta.OwnerReferences[0].Name, objMeta.Labels["task-uuid"])
		var names []string
		var size []int64
		// get meta data like output file name and output file size
		names, size, err := core.GetAllFileSizeInDirectory(mainConf.K8SConfig.OutputDirPrefix + "/" + successTask)
		if err != nil {
			c.logger.Error(err.Error())
		}

		// NOTE: if K8S cluster fail to delete job automatically, change this to false
		//core.DeleteK8SJob(context.Background(), successTask, false)
		core.DeleteK8SJob(context.Background(), successTask, true)

		// TODO: push signal to Temporal Workflow
		req := &model.UpdateStatusCheck{
			Success:  true,
			TaskID:   successTask,
			Filename: names,
			Filesize: size,
		}

		err = PushUpdateStatusToKafka(req) // push to kafka
		if err != nil {
			c.logger.Error(err.Error())
			return nil
		}

		_, err = eddamanager.UpdateTaskLogState(taskUUID)
		if err != nil {
			c.logger.Error(err.Error())
			return nil
		}

	}
	if podPhase == "Failed" {
		//AddToFailTask(objMeta.OwnerReferences[0].Name, objMeta.Labels["task-uuid"])
		failTask := objMeta.OwnerReferences[0].Name
		taskUUID := objMeta.Labels["task-uuid"]

		// NOTE: if K8S cluster fail to delete job automatically, change this to false
		//core.DeleteK8SJob(context.Background(), successTask, false)
		core.DeleteK8SJob(context.Background(), failTask, true)

		// TODO: push signal to Temporal Workflow
		/*
			// handler fail
			req := &model.UpdateStatusCheck{
				Success: false,
				TaskID:  failTask,
			}


			err := PushUpdateStatusToKafka(req) // push to kafka
			if err != nil {
				c.logger.Error(err.Error())
				return nil
			}
		*/

		_, err := eddamanager.UpdateTaskLogState(taskUUID)
		if err != nil {
			c.logger.Error(err.Error())
			return nil
		}
	}
	if podPhase == "Running" {
		// TODO: add to running task, calculate to start listen client
		cpu, err := strconv.ParseInt(objMeta.Labels["cpu"], 10, 64)
		if err != nil {
			c.logger.Error(err.Error())
		}
		ram, err := strconv.ParseInt(objMeta.Labels["ram"], 10, 64)
		if err != nil {
			c.logger.Error(err.Error())
		}
		taskUUID := objMeta.Labels["task-uuid"]

		core.AddJobInPro(objMeta.OwnerReferences[0].Name, ram, cpu, taskUUID)
		atomic.AddInt64(&core.CPULeft, -cpu)
		atomic.AddInt64(&core.RAMLeft, -ram)
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
	podStatus := core.GetPodStatus(obj)
	taskUUID, _ := core.GetTaskUUIDFromJobInPro(objMeta.OwnerReferences[0].Name)

	// Initital new job log reading
	if podPhase != "Succeeded" && podPhase != "Failed" {
		if (podStatus.ContainerStatuses != nil) && (podStatus.ContainerStatuses[0].Ready) {
			podSpec := core.GetPodSpec(obj)
			c.logger.Info("Sending new task log to edda agent with task id " + objMeta.OwnerReferences[0].Name)
			_, err := eddamanager.NewTaskLog(taskUUID, objMeta.Name, mainConf.K8SConfig.K8SNameSpace, podStatus.ContainerStatuses[0].Name, podStatus.ContainerStatuses[0].ContainerID, podSpec.NodeName)
			if err != nil {
				c.logger.Error(err.Error())
				return nil
			}
		}
	}

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

				req := &model.UpdateStatusCheck{
					Success:  true,
					TaskID:   objMeta.OwnerReferences[i].Name,
					Filename: names,
					Filesize: size,
				}

				_, ok := core.DeleteK8SJob(context.Background(), objMeta.OwnerReferences[i].Name, false)
				if ok {
					err = PushUpdateStatusToKafka(req) // push to kafka
					if err != nil {
						c.logger.Error(err.Error())
						return nil

					}
					_, err = eddamanager.UpdateTaskLogState(taskUUID)
					if err != nil {
						c.logger.Error(err.Error())
						return nil
					}
				}
			}
		}
	} else if podPhase == "Failed" {
		c.logger.Info("Job " + objMeta.Name + " is FAILED")

		for i := 0; i < len(objMeta.OwnerReferences); i++ {
			if (objMeta.OwnerReferences[i].Kind == "Job") && !GetPod(objMeta.OwnerReferences[i].Name) {
				SetPod(objMeta.OwnerReferences[i].Name)
				_, ok := core.DeleteK8SJob(context.Background(), objMeta.OwnerReferences[i].Name, false)
				if ok {
					// handler fail
					req := &model.UpdateStatusCheck{
						Success: false,
						TaskID:  objMeta.OwnerReferences[i].Name,
					}

					err = PushUpdateStatusToKafka(req) // push to kafka
					if err != nil {
						c.logger.Error(err.Error())
						return nil
					}

					_, err = eddamanager.UpdateTaskLogState(taskUUID)
					if err != nil {
						c.logger.Error(err.Error())
						return nil
					}
				}
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
