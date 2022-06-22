package business

import (
	"context"
	"strings"
	"time"

	"workflow/executor/core"
)

/*
func (c *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {

		//c.lg.LogInfo(fmt.Sprintf("Message claimed: %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic))

		var resp executorModel.SelectTaskResp
		var resp1 executorModel.DeleteTaskNoti

		if err := json.Unmarshal(message.Value, &resp); err == nil {
			if len(resp.Tasks) != 0 {
				c.lg.Info(fmt.Sprintf("Message claimed: %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic))
			}

			// if wanna blazing fast system, can do async here: Ex: go ProducerSelectTaskResp(...)
			//_ = ProcessSelectTaskResp(context.Background(), c.lg, &resp)
			/*
				if er != nil {
					c.lg.LogErr(er)
				}
		} else {
			c.lg.Info(fmt.Sprintf("Message claimed: %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic))
			if err1 := json.Unmarshal(message.Value, &resp1); err1 == nil {
				c.lg.Dataf(resp1.TaskID)

				// process delete job
				e := CheckDeleteTask(c.lg, &resp1)
				if err != nil {
					c.lg.Errorf(e.Error())
				}
			} else {
				c.lg.Info(fmt.Sprintf("Message claimed: %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic))
				c.lg.Errorf(err.Error())
				c.lg.Errorf(err1.Error())
			}
		}

		session.MarkMessage(message, "")
	}
	return nil
}

*/

func sleepContext(ctx context.Context, delay time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(delay):
	}
}

func GetParentDirectory(path string) (dir string) {
	var j int
	for i := len(path) - 1; i >= 0; i-- {
		if string(path[i]) == "/" {
			j = i
			break
		}
	}
	if j == 0 {
		return ""
	}
	var k = 0

	for {
		if k > j {
			break
		}
		dir = dir + string(path[k])
		k++
	}
	return
}

func GetFileFUSEPath(path string) (localfilepath string) {
	mainConf := core.GetMainConfig()
	return strings.ReplaceAll(path, mainConf.MinioEndpoint, mainConf.FUSEMountpoint)
}
