package utils

import (
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/vfluxus/valkyrie/core"
)

func GetBucketMetadata(bucket string) (filecount int32, totalfilesize int64, err error) {
	svc := core.GetMinIOClient()
	err = svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String("testbucket1"),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		for _, obj := range p.Contents {
			filecount++
			totalfilesize += *obj.Size
		}
		return true
	})

	return
}

func ConstructUnsafeDownloadURL(userid, runuuid, taskuuid, filename string, bucket string) (durl string) {
	mainConf := core.GetMainConfig()
	if runuuid == "" {
		return mainConf.MinioEndpoint + "/" + bucket + "/" + userid + "/" + filename
	} else {
		return mainConf.MinioEndpoint + "/" + bucket + "/" + userid + "/" + runuuid + "/" + taskuuid + "/" + filename
	}
}

func DeleteFile(path string) error {
	e := os.Remove(path)
	return e
}

func DeleteDir(path string) error {
	err := os.RemoveAll(path)
	return err
}

func GetTaskDirectoryRegexp(taskid string, path string) string {
	newpath := GetRunDir(path)
	spliter := strings.Split(taskid, "-")
	return newpath + spliter[0] + "-" + spliter[1] + "-*"
}

func DeleteDirWithRegex(reg string) error {
	return exec.Command("sh", "-c", "rm -rf "+reg).Run()
}

func GetRunDir(path string) (dir string) {
	spliter := strings.Split(path, "/")
	return path[0:(len(path) - len(spliter[len(spliter)-1]) - len(spliter[len(spliter)-2]) - 1)]
}

func GetTaskDir(path string) (dir string) {
	spliter := strings.Split(path, "/")
	return path[0:(len(path) - len(spliter[len(spliter)-1]))]
}
