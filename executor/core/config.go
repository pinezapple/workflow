package core

import (
	"context"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	mainConfig     *MainConfig
	mainConfigLock sync.RWMutex

	JobCount     int = 0
	JobCountLock sync.Mutex

	OldTaskCounter int32 = 0
	TaskHandleLock sync.RWMutex

	K8SClientSet     *kubernetes.Clientset
	K8SClientSetLock sync.RWMutex

	syncAndReadyFlag bool = false
	syncAndReadyLock sync.RWMutex
)

func IsGoodToGo(threshold int) bool {
	JobCountLock.Lock()
	defer JobCountLock.Unlock()
	if JobCount < threshold {
		JobCount++
		return true
	} else {
		return false
	}
}

func IncreaseJobCount() {
	JobCountLock.Lock()
	defer JobCountLock.Unlock()
	JobCount++
}

func DecreaseJobCount() {
	JobCountLock.Lock()
	defer JobCountLock.Unlock()
	JobCount--
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

// HTTPServerConf binding configuration for webserver
type HTTPServerConf struct {
	PublicIP  string
	Port      int
	Cert      string
	Key       string
	ClientCAs string // csv list of trusted CAs
}

type MainConfig struct {
	ServiceName          string
	MaximumConcurrentJob int
	FUSEMountpoint       string
	MinioEndpoint        string
	FailOverTime         int
	HttpServerConfig     *HTTPServerConf
	K8SConfig            *K8SConf
	Environment          string
	LogLevel             string
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
	Initk8sClientset()
}
