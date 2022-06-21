package core

import (
	"os/exec"
)

func MountBucketToDir(bucket, mountpoint, authFile, storagePath string) (err error) {
	_, err = exec.Command("sh", "-c", "./automount.sh"+" "+bucket+" "+mountpoint+" "+authFile+" "+storagePath).Output()
	if err != nil {
		return err
	}
	return
}

// maybe add an unmount script with this
//fusermount -u /home/pine/fuse-mountpoint/bucket1
