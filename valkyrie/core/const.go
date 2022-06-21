package core

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
)

const (
	testDSN = "root:123@/test?charset=utf8&collation=utf8_general_ci&parseTime=true&loc=Asia%2FHo_Chi_Minh"
)

var (
	ErrBadRequest = fmt.Errorf("Bad request")

	ErrNoAckFound = fmt.Errorf("No ack found")

	ErrExtTermChanCapInvalid = fmt.Errorf("Term chan capacity is invalid")

	// ErrDBObjNull indicate DB Object is nil
	ErrDBObjNull = fmt.Errorf("DB Object is nil")
)

// TODO: form minio key name
func FormMinioKeyName(filename string) *string {
	return aws.String("")
}

func GetFileName(path []string) (files []string) {
	for i := 0; i < len(path); i++ {
		tg := strings.Split(path[i], "/")
		files = append(files, tg[len(tg)-1])
	}
	return
}
