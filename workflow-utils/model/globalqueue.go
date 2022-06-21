package model

import (
	"sync"

	linkedList "github.com/lichti/go-linked-list"
)

func (globalMapper *GlobalMap) GetQueueOfLevel(level int) (q *GlobalQueue) {
	globalMapper.RLock()
	q = globalMapper.Mapper[level]
	globalMapper.RUnlock()
	return q
}

type GlobalMap struct {
	sync.RWMutex
	Mapper []*GlobalQueue
}

//type NodeJobInQueueDB struct{
//	NodeId      string `json:"node_id" gorm:"column:node_id; primary_key"`
//	CurrentJobId  string `json:"current_job_id" gorm:"column:current_job_id"`
//	NextNodeId    string `json:"next_job_id" gorm:"column:next_job_id"`
//}

type NodeInQueue struct {
	JobId string `json:"job_id" gorm:"column:job_in_node; primaryKey"`
	JobLogic JobLogic `gorm:"foreignKey:job_in_node;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

//-----------------GlobleQueue-------------------
type GlobalQueue struct {
	Queue linkedList.List `gorm:"-"`
	//Head     *JobLogic
	SyncChan           chan *PushChildrenJob `gorm:"-"`
	Level              int                   `json:"level" gorm:"column:level, primary_key"`
	Status             string                `json:"status" gorm:"column:status"`
	MaxLeftBehindReq   uint32                `json:"max_left_behind_req" gorm:"column:max_left_behind_req"`
	LeftBehindReqCount uint32                `json:"left_behind_request_count" gorm:"column:left_behind_request_count"`
	sync.RWMutex       `gorm:"-"`
}

type PushChildrenJob struct {
	loggin chan int
	node   *JobLogic
}

func (g *GlobalQueue) GetOneJobFromQueue(CPU int64, RAM int64) (job *JobLogic, RAMleft, CPUleft int64) {
	g.Lock()
	n := g.Queue.First
	for {
		if n == nil {
			break
		}

		j := (n.Value).(*JobLogic)
		if j.ECPU <= CPU && j.ERam <= RAM {
			job = j
			RAM = RAM - j.ERam
			CPU = CPU - j.ECPU
			break
		}
		n = n.Next
	}
	g.Unlock()
	return job, RAM, CPU
}

func (g *GlobalQueue) GetMultipleJobFromQueue(CPU int64, RAM int64) (jobs []*JobLogic, RAMleft, CPUleft int64) {
	var userPenaltyMap = make(map[uint32]bool)
	g.Lock()
	n := g.Queue.First
	for {
		if n == nil {
			break
		}

		j := (n.Value).(*JobLogic)
		_, ok := userPenaltyMap[j.UserID]

		if j.ECPU <= CPU && j.ERam <= RAM && !ok {
			jobs = append(jobs, j)
			CPU = CPU - j.ECPU
			RAM = RAM - j.ERam
		} else {
			userPenaltyMap[j.UserID] = true
		}
		n = n.Next
	}
	g.Unlock()
	return jobs, RAM, CPU
}

func (g *GlobalQueue) ListenFromQueueChannel(wg *sync.WaitGroup) {
	for {
		j, ok := <-g.SyncChan
		if ok {
			g.Queue.Push(j.node)
			j.loggin <- 1
		} else {
			wg.Done()
			return
		}
	}
}
