package business

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"workflow/valkyrie/core"
)

func DeleteFromStorage(bucket string, name string, size int64) (ok bool, err error) {
	iter, err := strconv.ParseInt(bucket[6:], 10, 64)
	if err != nil {
		return false, err
	}

	// Delete from minio
	svc := core.GetMinIOClient()
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	})

	if err != nil {
		return false, err
	}

	// Delele from Ram bucket
	core.DeleteFromBucket(int(iter), size)

	return
}
