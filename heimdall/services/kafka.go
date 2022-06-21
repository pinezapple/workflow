package services

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"sync"
// 	"time"

// 	"github.com/Shopify/sarama"
// 	"workflow/heimdall/core"
// )

// // ----------------------------------------------------------------------------------------------------------------------------
// // -------------------------------------------------------- CONSUMER ----------------------------------------------------------

// var (
// 	KafkaConsumer     sarama.ConsumerGroup
// 	KafkaConsumerLock sync.RWMutex // prevent reading while setup
// )

// // consumerHandler implement the sarama.ConsumerGroupHandler
// type consumerHandler struct{}

// // Setup is run at the beginning of a new session, before ConsumeClaim.
// func (cHdl *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
// 	return nil
// }

// // Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited but before the offsets are committed for the very last time.
// func (cHdl *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
// 	return nil
// }

// // ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages(). Once the Messages() channel is closed, the Handler must finish its processing loop and exit.
// func (cHdl *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {
// 	// range channel to listening to message
// 	for message := range claim.Messages() {
// 		var resp struct{} // unmarshal to this
// 		if err := json.Unmarshal(message.Value, &resp); err != nil {
// 			// TODO: do sth here and log
// 			log.Println("Consume message error: ", err)
// 		}
// 		// TODO: process the response

// 		session.MarkMessage(message, "") // mark message is processed
// 	}

// 	return nil // why not return err ??? break the function ?
// }

// // SetKafkaConsumerWithConfig ...
// func SetKafkaConsumerWithConfig(kcf *core.KafkaConfig) {
// 	// create config
// 	config := sarama.NewConfig()
// 	version, err := sarama.ParseKafkaVersion(kcf.Version)
// 	if err != nil {
// 		log.Fatalln("Can not create parse kafka version, error: ", err)
// 	}

// 	config.Version = version
// 	config.Consumer.Return.Errors = true
// 	config.Consumer.Offsets.Initial = sarama.OffsetOldest
// 	config.Consumer.Offsets.Retention = 36 * time.Hour

// 	// create consumer
// 	newConsumer, err := sarama.NewConsumerGroup(kcf.ConsumerBrokers, kcf.ConsumerGroup, config)
// 	if err != nil {
// 		log.Fatalln("Can not create kafka consumer, error:", err)
// 	}

// 	// close existed kafka consumer
// 	if old := GetKafkaConsumer(); old != nil {
// 		old.Close()
// 	}

// 	// set up new kafka consumer to global var
// 	KafkaConsumerLock.RLock()
// 	KafkaConsumer = newConsumer
// 	KafkaConsumerLock.RUnlock()
// }

// // GetKafkaConsumer create a new kafka consumer
// func GetKafkaConsumer() sarama.ConsumerGroup {
// 	KafkaConsumerLock.RLock()
// 	consumer := KafkaConsumer
// 	KafkaConsumerLock.RUnlock()
// 	return consumer
// }

// // ReceiveMessageFromKafka consume kafka messages
// func ReceiveMessageFromKafka(parentCtx context.Context) (fn func(), err error) {
// 	// TODO: Add log

// 	// func to continuously listening kafka
// 	fn = func() {
// 		log.Println("Listening to kafka ... ")
// 		consumer := GetKafkaConsumer()
// 		consumerHandling := &consumerHandler{}

// 		for {
// 			if err = consumer.Consume(parentCtx, []string{"some topics here"}, consumerHandling); err != nil {
// 				log.Println("Kafka consume error: ", err)
// 			}

// 			if err = parentCtx.Err(); err != nil {
// 				log.Println("Kafka context error: ", err)
// 			}
// 		}
// 	}

// 	return fn, err
// }

// // ----------------------------------------------------------------------------------------------------------------------------
// // -------------------------------------------------------- PRODUCER ----------------------------------------------------------

// var (
// 	KafkaProducer     sarama.SyncProducer
// 	KafkaProducerLock sync.RWMutex // prevent reading while setup
// )

// // SetKafkaProducerWithConfig ...
// func SetKafkaProducerWithConfig(kcf *core.KafkaConfig) {
// 	// create config
// 	config := sarama.NewConfig()
// 	version, err := sarama.ParseKafkaVersion(kcf.Version)
// 	if err != nil {
// 		log.Panicln("Can not parse kafka version, error: ", err)
// 	}

// 	config.Version = version
// 	config.Producer.RequiredAcks = sarama.WaitForAll
// 	config.Producer.Compression = sarama.CompressionSnappy
// 	config.Producer.Retry.Max = 5
// 	config.Producer.Return.Successes = true

// 	// create producer
// 	newProducer, err := sarama.NewSyncProducer(kcf.ProducerBrokers, config)
// 	if err != nil {
// 		log.Fatalln("Can not create kafka producer, error: ", err)
// 	}

// 	// close any producer if exist
// 	if old := GetKafkaConsumer(); old != nil {
// 		old.Close()
// 	}

// 	// set kafka producer to global var
// 	KafkaProducerLock.RLock()
// 	KafkaProducer = newProducer
// 	KafkaProducerLock.RUnlock()
// }

// // GetKafkaProducer ...
// func GetKafkaProducer() sarama.SyncProducer {
// 	KafkaProducerLock.RLock()
// 	producer := KafkaProducer
// 	KafkaProducerLock.RUnlock()
// 	return producer
// }

// // PushSthToKafka push to kafka template
// func PushSthToKafka(sth struct{}) (err error) {
// 	producer := GetKafkaProducer()
// 	if producer != nil {
// 		return fmt.Errorf("No kafka producer found")
// 	}
// 	// marshal data to send
// 	var value []byte
// 	if value, err = json.Marshal(sth); err != nil {
// 		return err
// 	}
// 	// send message
// 	partition, offset, e := producer.SendMessage(&sarama.ProducerMessage{
// 		Topic: "some topics",
// 		Value: sarama.ByteEncoder(value),
// 	})
// 	if err != nil {
// 		return fmt.Errorf("Can not send message value: %v, partition: %v, offset: %v, error: %v", sth, partition, offset, e)
// 	}

// 	return nil
// }
