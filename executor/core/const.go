package core

import (
	"fmt"
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
