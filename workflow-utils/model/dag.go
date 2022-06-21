package model

import (
	"log"
	"sync/atomic"
)

func (n *JobLogic) AddChild(childNode *JobLogic) (ok bool) {
	n.ChildrenJobs = append(n.ChildrenJobs, childNode)
	childNode.ParentJobs = append(childNode.ParentJobs, n)
	return true
}

func (n *JobLogic) PushFirstJobToQueue(gbm *GlobalMap) (ok bool) {
	q := gbm.GetQueueOfLevel(n.QueueLevel)
	job := &PushChildrenJob{
		loggin: make(chan int, 1),
		node:   n,
	}
	q.SyncChan <- job
	<-job.loggin
	return true
}

func (n *JobLogic) PushChildrenToQueue(gbm *GlobalMap) (ok bool) {
	q := gbm.GetQueueOfLevel(n.QueueLevel)
	if !q.Queue.Find(n) {
		log.Fatalf("Fucked up")
	}
	/*
		fmt.Println(n.GMeta.gID)
		fmt.Println(n)
		fmt.Println(q.Queue.Find(n))
	*/
	q.Queue.Remove(n)
	for i := 0; i < len(n.ChildrenJobs); i++ {
		atomic.AddUint32(&n.ChildrenJobs[i].ParentsDoneCount, 1) // atomic value later on
		if n.ChildrenJobs[i].ParentsDoneCount == uint32(len(n.ChildrenJobs[i].ParentJobs)) {
			//			fmt.Println(n.Children[i].GMeta.gID)
			job := &PushChildrenJob{
				loggin: make(chan int, 1),
				node:   n.ChildrenJobs[i],
			}
			q.SyncChan <- job
			<-job.loggin
		}
	}
	return true
}

func (n *JobLogic) PushFirstChildrenToQueue(gbm *GlobalMap) (ok bool) {
	q := gbm.GetQueueOfLevel(n.QueueLevel)
	for i := 0; i < len(n.ChildrenJobs); i++ {
		job := &PushChildrenJob{
			loggin: make(chan int, 1),
			node:   n.ChildrenJobs[i],
		}
		q.SyncChan <- job
		<-job.loggin
	}
	return true
}
