package business

import (
	"context"
	"sync/atomic"
	"time"

	"workflow/workflow-utils/model"
	"workflow/executor/api/scheduler"
	"workflow/executor/core"
)

func sleepContext(ctx context.Context, delay time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(delay):
	}
}

func SelectTaskDaemons(parentCtx context.Context) (fn model.Daemon, err error) {
	fn = func() {
		lg := core.GetLogger()
		lg.Info("Starting selecting task daemon")

		mainConf := core.GetMainConfig()

		sleepContext(parentCtx, time.Duration(mainConf.FailOverTime)*time.Second)
		if !core.GetSyncFlag() {
			core.SetSyncFlag()
		}
		if core.GetSyncFlag() && atomic.LoadInt32(&core.OldTaskCounter) == 0 {
			//			fmt.Println("in executor resume without task")
			_, err := scheduler.ExecutorResumeWithoutTask()
			if err != nil {
				lg.Errorf(err.Error())
			}
		}

		for {
			if parentCtx.Err() != nil {
				return
			}

			sleepContext(parentCtx, time.Duration(mainConf.SelectTaskIntervalTime)*time.Second)
			//fmt.Println(core.CPULeft)
			//fmt.Println(core.RAMLeft)
			if core.GetLengthOfJobInPro() <= mainConf.MaximumConcurrentJob {
				_, er := scheduler.SelectTask(atomic.LoadInt64(&core.CPULeft), atomic.LoadInt64(&core.RAMLeft))
				if er != nil {
					lg.Errorf(er.Error())
				}
			}
		}
	}
	return fn, nil
}
