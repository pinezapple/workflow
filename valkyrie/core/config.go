package core

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"workflow/valkyrie/model"
)

var (
	minioClient     *s3.S3
	minioClientLock sync.RWMutex

	minioBucketSlice     []*MinioBucketCurrentConf
	minioBucketSliceLock sync.Mutex

	DownloadLock sync.RWMutex

	mainConfig     *MainConfig
	mainConfigLock sync.RWMutex

	DB     *gorm.DB
	DBLock sync.RWMutex
)

func RestartMinioBuckets(size []int64, count []int32) {
	minioBucketSliceLock.Lock()
	defer minioBucketSliceLock.Unlock()

	for i := 0; i < len(size); i++ {
		var newBucket = &MinioBucketCurrentConf{
			FileCount: count[i],
			TotalSize: size[i],
		}
		minioBucketSlice = append(minioBucketSlice, newBucket)
	}
}

func GetMinioBucket(size int64) (name string, iter int, newbucket bool) {
	minioBucketSliceLock.Lock()
	defer minioBucketSliceLock.Unlock()
	for i := 0; i < len(minioBucketSlice); i++ {
		if minioBucketSlice[i].PushToBucket(size) {
			return "bucket" + strconv.Itoa(i), i, false
		}
	}
	var newBucket = &MinioBucketCurrentConf{
		FileCount: 1,
		TotalSize: size,
	}
	minioBucketSlice = append(minioBucketSlice, newBucket)
	return "bucket" + strconv.Itoa(len(minioBucketSlice)), len(minioBucketSlice), true
}

func DeleteFromBucket(i int, size int64) (ok bool) {
	minioBucketSliceLock.Lock()
	defer minioBucketSliceLock.Unlock()

	if len(minioBucketSlice)-1 < i {
		return false
	}

	minioBucketSlice[i].DeleteFromBucket(size)
	return true
}

type MinioBucketCurrentConf struct {
	TotalSize int64
	FileCount int32
}

func (m *MinioBucketCurrentConf) PushToBucket(size int64) (ok bool) {
	mainConf := GetMainConfig()
	if (m.FileCount+1) <= mainConf.MinioConfig.BucketMaxCount && (m.TotalSize+size) <= mainConf.MinioConfig.BucketMaxSize {
		m.TotalSize += size
		m.FileCount++
		return true
	} else {
		return false
	}
}

func (m *MinioBucketCurrentConf) DeleteFromBucket(size int64) {
	m.FileCount--
	m.TotalSize -= size
}

// DB config for future use
type DBConf struct {
	DBName   string
	Username string
	Password string
	Hostname string
	Port     string
}

type ExpectJWTConf struct {
	// Audiences is a list of expected receivers or uses of the token.
	Audiences []string
	// Expiration is the Unix timestamp at which the token becomes expired.
	Expiration int64
	// Issuers is a list of acceptable issuers to expect tokens to contain.
	Issuers []string
	// Purpose is an optional field indicating the type of the token (access,
	// refresh, etc.)
	Purpose string
}

// HTTPServerConf binding configuration for webserver
type HTTPServerConf struct {
	PublicIP  string
	Port      int
	Cert      string
	Key       string
	ClientCAs string // csv list of trusted CAs
}

type MinioConf struct {
	AccessKey      string
	SecretKey      string
	Region         string
	Endpoint       string
	BucketMaxSize  int64
	BucketMaxCount int32
}

type OutterBoundHTTPOption struct {
	Addr    string
	KeyFile string
}

type MainConfig struct {
	HttpServerConfig *HTTPServerConf
	MinioConfig      *MinioConf
	DBConfig         *DBConf
	HeimdallConfig   *OutterBoundHTTPOption

	ServiceName string
	Environment string
	LogLevel    string

	UpdateStorageChangeInterval int
	ExpiredUpdateInterval       int
	APIRetryCount               int

	HardDiskOnly bool
	NoDelete     bool

	UploadPenalty    int
	NormalFileTTL    int
	ImportantFileTTL int

	MinioEndpoint       string
	MinioAuthenFile     string
	FUSEMountpoint      string
	InputDirPrefix      string
	UploadTempDirPrefix string
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

func SetMinIOClient(storageConfig *aws.Config) {
	minioClientLock.Lock()
	defer minioClientLock.Unlock()
	newSession := session.New(storageConfig)
	minioClient = s3.New(newSession)
}

func GetMinIOClient() (client *s3.S3) {
	minioClientLock.RLock()
	defer minioClientLock.RUnlock()
	client = minioClient
	return
}

func SetDBObj(db *gorm.DB) {
	DBLock.Lock()
	DB = db
	DBLock.Unlock()
}

func GetDBObj() *gorm.DB {
	DBLock.RLock()
	db := DB
	DBLock.RUnlock()
	return db
}

func constructPostgresDSN(conf *DBConf) (dsn string) {
	dsn = "user=" + conf.Username + " password=" + conf.Password + " dbname=" + conf.DBName + " host=" + conf.Hostname + " port=" + conf.Port + " sslmode=disable"
	return
}
func InitialConnectToDatabase(ctx context.Context) (db *gorm.DB, er error) {
	mainConf := GetMainConfig()
	dsn := constructPostgresDSN(mainConf.DBConfig)
	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold: time.Second,     // Slow SQL threshold
			LogLevel:      gormLogger.Info, // Log level
			Colorful:      true,            // Disable color
		},
	)

	db, er = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	SetDBObj(db)
	return
}

func InitialTableDatabase(ctx context.Context) (er error) {
	// err := make([]error, 2)
	db := GetDBObj()
	tx := db.WithContext(ctx).Begin() // create transaction
	tx.SavePoint("sp1")

	er = tx.AutoMigrate(&model.Bucket{}, &model.GeneratedFile{}, &model.UploadedFile{}, &model.Sample{}, &model.Dataset{}, &model.SampleContent{})
	if er != nil {
		tx.RollbackTo("sp1")
	}

	return tx.Commit().Error
}

func InitCore(ctx context.Context) {
	InitialConnectToDatabase(ctx)
	InitialTableDatabase(ctx)

	if !mainConfig.HardDiskOnly {
		// Configure to use MinIO Server
		s3Config := &aws.Config{
			Credentials:      credentials.NewStaticCredentials(mainConfig.MinioConfig.AccessKey, mainConfig.MinioConfig.SecretKey, ""),
			Endpoint:         aws.String(mainConfig.MinioConfig.Endpoint),
			Region:           aws.String(mainConfig.MinioConfig.Region),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
		}
		SetMinIOClient(s3Config)
	}
}
