package utils

import (
	"fmt"
	"path/filepath"
	"strings"

	"workflow/workflow-utils/model"
)

var (
	ErrTaskFailed = fmt.Errorf("Task failed")
	ErrTaskDone   = fmt.Errorf("Task done")
)

func GetFileName(path []string) (files []string) {
	for i := 0; i < len(path); i++ {
		tg := strings.Split(path[i], "/")
		files = append(files, tg[len(tg)-1])
	}
	return
}

func IsRegex(a string) (ok bool) {
	if strings.Contains(a, "?") {
		return true
	} else if strings.Contains(a, "*") {
		return true
	} else if strings.Contains(a, "[") {
		return true
	} else if strings.Contains(a, "\\") {
		return true
	} else {
		return false
	}
}

func FillInput(regexes []*model.ParamWithRegex, fileName []string, filePath []string, filesize []int64) {
	for i := 0; i < len(regexes); i++ {
		fmt.Println(regexes[i].Regex)
		fmt.Println(regexes[i].SecondaryFiles)
		for j := 0; j < len(regexes[i].Regex); j++ {
			for k := 0; k < len(filePath); k++ {
				ok, _ := filepath.Match(regexes[i].Regex[j], fileName[k])
				if ok {
					var ok1 bool = false
					for k1 := 0; k1 < len(regexes[i].From); k1++ {
						if strings.Contains(filePath[k], regexes[i].From[k1]) {
							ok1 = true
							break
						}
					}

					if ok1 {
						newFile := &model.FilteredFiles{
							Filepath: filePath[k],
							Filesize: filesize[k],
						}
						regexes[i].Files = append(regexes[i].Files, newFile)
					}
				}
			}
		}
		for i1 := 0; i1 < len(regexes[i].Files); i1++ {
			fmt.Println(regexes[i].Files[i1].Filepath)
			fmt.Println(regexes[i].Files[i1].Filesize)
		}
	}
}
