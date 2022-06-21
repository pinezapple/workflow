package business

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/vfluxus/workflow-utils/model"
	"github.com/vfluxus/workflow/executor/core"
	executorModel "github.com/vfluxus/workflow/executor/model"
)

type consumerHandler struct {
	lg *core.LogFormat
}

func (c *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

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
			*/
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

func ReceiveMessageFromKafka(parentCtx context.Context) (fn model.Daemon, err error) {
	fn = func() {
		lg := core.GetLogger()

		lg.Info("Start listening to kafka")

		consumer := core.GetKafkaConsumer()
		consumerHandling := &consumerHandler{
			lg: lg,
		}

		for {
			if err := consumer.Consume(parentCtx, core.GetMainConfig().KafkaConfig.ConsumerTopics, consumerHandling); err != nil {
				lg.Errorf(err.Error())
			}

			if parentCtx.Err() != nil {
				return
			}
		}
	}
	return fn, err
}

func PushAckToKafka(jobs *executorModel.TaskAck) (err error) {
	mainConf := core.GetMainConfig()
	producer := core.GetKafkaProducer()
	if producer == nil {
		return fmt.Errorf("no kafka producer found")
	}

	var value []byte
	if value, err = json.Marshal(jobs); err != nil {
		return err
	}

	partition, offset, e := producer.SendMessage(&sarama.ProducerMessage{
		Topic: mainConf.KafkaConfig.ProducerTopics[0],
		Value: sarama.ByteEncoder(value),
	})

	if e != nil {
		mar, _ := json.Marshal(map[string]interface{}{
			"uuid":      jobs.TaskID,
			"partition": partition,
			"offset":    offset,
			"error":     e,
		})
		return fmt.Errorf("%s", string(mar))
	}

	return
}

func PushUpdateStatusToKafka(req *executorModel.UpdateStatusCheck) (err error) {
	mainConf := core.GetMainConfig()
	producer := core.GetKafkaProducer()
	if producer == nil {
		return fmt.Errorf("no kafka producer found")
	}

	var value []byte
	if value, err = json.Marshal(req); err != nil {
		return err
	}

	partition, offset, e := producer.SendMessage(&sarama.ProducerMessage{
		Topic: mainConf.KafkaConfig.ProducerTopics[1],
		Value: sarama.ByteEncoder(value),
	})

	if e != nil {
		mar, _ := json.Marshal(map[string]interface{}{
			"uuid":      req.TaskID,
			"partition": partition,
			"offset":    offset,
			"error":     e,
		})
		return fmt.Errorf("%s", string(mar))
	}

	return
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
