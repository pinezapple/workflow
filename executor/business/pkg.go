package business

import (
	"context"
	"strings"
	"time"

	"workflow/executor/core"
)

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
