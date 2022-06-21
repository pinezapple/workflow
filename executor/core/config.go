package core

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	mainConfig     *MainConfig
	mainConfigLock sync.RWMutex

	CPULeft int64
	RAMLeft int64

	TaskContextMap     map[string]context.CancelFunc
	TaskContextMapLock sync.RWMutex

	JobInPro     map[string]*JobInProgressObj
	JobInProLock sync.RWMutex

	OldTaskCounter int32 = 0
	TaskHandleLock sync.RWMutex

	KafkaProducer     sarama.SyncProducer
	KafkaConsumer     sarama.ConsumerGroup
	KafkaProducerLock sync.RWMutex
	KafkaConsumerLock sync.RWMutex

	K8SClientSet     *kubernetes.Clientset
	K8SClientSetLock sync.RWMutex

	syncAndReadyFlag bool = false
	syncAndReadyLock sync.RWMutex

	K8SDeleteJobMap     map[string]time.Time
	K8SDeleteJobMapLock sync.Mutex
)

type JobInProgressObj struct {
	Ram      int64
	Cpu      int64
	TaskUUID string
}

func GetK8SDeleteJobMap() (mapper map[string]time.Time) {
	K8SDeleteJobMapLock.Lock()
	defer K8SDeleteJobMapLock.Unlock()
	return K8SDeleteJobMap
}

func AddToK8SDeleteJobMap(key string, value time.Time) {
	K8SDeleteJobMapLock.Lock()
	defer K8SDeleteJobMapLock.Unlock()
	K8SDeleteJobMap[key] = value
}

func RemoveFromK8SDeleteJobMap(key string) {
	K8SDeleteJobMapLock.Lock()
	defer K8SDeleteJobMapLock.Unlock()
	delete(K8SDeleteJobMap, key)
}

func GetSyncFlag() bool {
	syncAndReadyLock.RLock()
	defer syncAndReadyLock.RUnlock()
	return syncAndReadyFlag
}

func SetSyncFlag() {
	syncAndReadyLock.Lock()
	defer syncAndReadyLock.Unlock()
	syncAndReadyFlag = true
	return
}

func AddJobInPro(name string, ram, cpu int64, taskUUID string) {
	JobInProLock.Lock()
	defer JobInProLock.Unlock()
	if _, ok := JobInPro[name]; ok {
		return
	} else {
		tmp := &JobInProgressObj{
			Ram:      ram,
			Cpu:      cpu,
			TaskUUID: taskUUID,
		}
		JobInPro[name] = tmp
	}
}

func GetTaskUUIDFromJobInPro(name string) (taskUUID string, ok bool) {
	JobInProLock.RLock()
	defer JobInProLock.RUnlock()
	j, ok := JobInPro[name]
	if !ok {
		return "", false
	} else {
		return j.TaskUUID, true
	}
}

func IfJobInPro(name string) bool {
	JobInProLock.RLock()
	defer JobInProLock.RUnlock()
	_, ok := JobInPro[name]
	return ok
}

func GetLengthOfJobInPro() int {
	JobInProLock.RLock()
	defer JobInProLock.RUnlock()
	return len(JobInPro)
}

func DeleteJobInPro(name string) (ram, cpu int64, taskUUID string, ok bool) {
	JobInProLock.Lock()
	defer JobInProLock.Unlock()
	if _, ok := JobInPro[name]; ok {
		ram = JobInPro[name].Ram
		cpu = JobInPro[name].Cpu
		taskUUID = JobInPro[name].TaskUUID
		delete(JobInPro, name)
		return ram, cpu, taskUUID, true
	} else {
		return 0, 0, "", false
	}
}

func AddTaskContext(id string, cancl context.CancelFunc) {
	TaskContextMapLock.Lock()
	TaskContextMap[id] = cancl
	TaskContextMapLock.Unlock()
}

func DeleteTaskContext(id string) (ok bool) {
	TaskContextMapLock.Lock()
	defer TaskContextMapLock.Unlock()
	if _, ok := TaskContextMap[id]; ok {
		delete(TaskContextMap, id)
		return true
	} else {
		return false
	}
}

func IfTaskContext(id string) (ok bool) {
	TaskContextMapLock.Lock()
	defer TaskContextMapLock.Unlock()
	cancel, ok := TaskContextMap[id]
	if ok {
		cancel()
		delete(TaskContextMap, id)
		return true
	} else {
		return false
	}
}

// HTTPServerConf binding configuration for webserver
type HTTPServerConf struct {
	PublicIP  string
	Port      int
	Cert      string
	Key       string
	ClientCAs string // csv list of trusted CAs
}

// Kafka config
type KafkaConf struct {
	ProducerBrokers []string
	ConsumerBrokers []string
	ProducerTopics  []string
	ConsumerTopics  []string
	ConsumerGroup   string
}

type MainConfig struct {
	ServiceName            string
	SelectTaskIntervalTime int64
	MaximumConcurrentJob   int
	FUSEMountpoint         string
	MinioEndpoint          string
	APIRetryCount          int
	FailOverTime           int
	HttpServerConfig       *HTTPServerConf
	KafkaConfig            *KafkaConf
	K8SConfig              *K8SConf
	SchedulerConfig        *OutterBoundHTTPOption
	EddaConfig             *OutterBoundHTTPOption
	Environment            string
	LogLevel               string
}

type OutterBoundHTTPOption struct {
	Addr    string
	KeyFile string
}

type K8SConf struct {
	K8SConfigFile          string
	OutputDirPrefix        string
	K8SNameSpace           string
	DeleteJob              bool
	JobTTLAfterFinished    int
	JobDeleteIntervalCheck int
	NodeLabelKey           string
	NodeLabelValue         string
}

func GetMainConfig() (mcf *MainConfig) {
	mainConfigLock.RLock()
	mcf = mainConfig
	mainConfigLock.RUnlock()
	return mcf
}

func SetMainConfig(mcf *MainConfig) {
	mainConfigLock.Lock()
	mainConfig = mcf
	mainConfigLock.Unlock()
}

func SetKafkaProducerWithConfig(v *KafkaConf) {
	config := sarama.NewConfig()

	version, err := sarama.ParseKafkaVersion("2.6.0")
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}
	config.Version = version

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	newProducer, err := sarama.NewSyncProducer(v.ProducerBrokers, config)
	if err != nil {
		panic(err)
	}

	if old := GetKafkaProducer(); old != nil {
		old.Close()
	}

	KafkaProducerLock.Lock()
	KafkaProducer = newProducer
	KafkaProducerLock.Unlock()
}

func GetKafkaProducer() sarama.SyncProducer {
	KafkaProducerLock.RLock()
	producer := KafkaProducer
	KafkaProducerLock.RUnlock()
	return producer
}

func SetKafkaConsumerClusterWithConfig(v *KafkaConf) {

	/*
		config := cluster.NewConfig()
		config.Consumer.Return.Errors = true
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
		config.Consumer.Offsets.Retention = 36 * time.Hour

		newConsumer, err := cluster.NewConsumer(v.ConsumerBrokers, v.ConsumerGroup, v.ConsumerTopics, config)
		if err != nil {
			panic(err)
		}

	*/

	config := sarama.NewConfig()

	version, err := sarama.ParseKafkaVersion("2.6.0")
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}
	config.Version = version

	//	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.Retention = 36 * time.Hour

	newConsumer, err := sarama.NewConsumerGroup(v.ConsumerBrokers, v.ConsumerGroup, config)
	if err != nil {
		panic(err)
	}

	if old := GetKafkaConsumer(); old != nil {
		_ = old.Close()
	}

	KafkaConsumerLock.RLock()
	KafkaConsumer = newConsumer
	KafkaConsumerLock.RUnlock()
}

func GetKafkaConsumer() sarama.ConsumerGroup {
	KafkaConsumerLock.RLock()
	consumer := KafkaConsumer
	KafkaConsumerLock.RUnlock()
	return consumer
}

func Initk8sClientset() {
	config, err := clientcmd.BuildConfigFromFlags("", mainConfig.K8SConfig.K8SConfigFile)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	SetK8SClientSet(clientset)
}

func GetK8SClientSet() (cl *kubernetes.Clientset) {
	K8SClientSetLock.RLock()
	cl = K8SClientSet
	K8SClientSetLock.RUnlock()
	return
}

func SetK8SClientSet(cl *kubernetes.Clientset) {
	K8SClientSetLock.Lock()
	K8SClientSet = cl
	K8SClientSetLock.Unlock()
}

func InitCore(ctx context.Context) {
	// init kafka

	logger = InitLogger(mainConfig)
	SetKafkaProducerWithConfig(mainConfig.KafkaConfig)
	SetKafkaConsumerClusterWithConfig(mainConfig.KafkaConfig)

	Initk8sClientset()

	TaskContextMap = make(map[string]context.CancelFunc)
	JobInPro = make(map[string]*JobInProgressObj)
	K8SDeleteJobMap = make(map[string]time.Time)
}
