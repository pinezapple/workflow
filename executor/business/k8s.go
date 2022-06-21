package business

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vfluxus/workflow-utils/model"
	"github.com/vfluxus/workflow/executor/core"
	executorModel "github.com/vfluxus/workflow/executor/model"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultBashCommand = "/bin/bash"
	defaultBashParam   = "-c"
)

func CheckDeleteTask(lg *core.LogFormat, resp *executorModel.DeleteTaskNoti) (err error) {
	core.TaskContextMapLock.RLock()
	defer core.TaskContextMapLock.RUnlock()

	if core.IfTaskContext(resp.TaskID) {
		return nil
	} else {
		err, _ = core.DeleteK8SJob(context.Background(), resp.TaskID, true)
		return
	}
}

/*
func ProcessSelectTaskResp(ctx context.Context, lg *core.LogFormat, resp *executorModel.SelectTaskResp) (err error) {
	// create context first for safety reason
	core.TaskHandleLock.Lock()
	var childCtxs []context.Context
	for i := 0; i < len(resp.Tasks); i++ {
		job := resp.Tasks[i]
		childCtx, cancel := context.WithCancel(ctx)
		childCtxs = append(childCtxs, childCtx)
		core.AddTaskContext(job.TaskID, cancel)
	}
	core.TaskHandleLock.Unlock()

	for i := 0; i < len(resp.Tasks); i++ {
		job := resp.Tasks[i]
		if core.IfJobInPro(job.TaskID) {
			continue
		}

		core.AddJobInPro(job.TaskID, job.ERam, job.ECPU, job.TaskUUID)
		err = CreateK8SJob(childCtxs[i], &job, lg)
		if err != nil {
			//if err == fmt.Errorf("Missing secondary file") {
			// handler fail
			req := &executorModel.UpdateStatusCheck{
				Success: false,
				TaskID:  job.TaskID,
			}

			PushUpdateStatusToKafka(req) // push to kafka
			/*
				} else {
					PushAckToKafka(&executorModel.TaskAck{
						TaskID: job.TaskID,
						Status: model.StatusInqueue,
					})

				}
			lg.Errorf(err.Error())
			core.DeleteTaskContext(job.TaskID)
			core.DeleteJobInPro(job.TaskID)
		} else {
			atomic.AddInt64(&core.CPULeft, -job.ECPU)
			atomic.AddInt64(&core.RAMLeft, -job.ERam)
			PushAckToKafka(&executorModel.TaskAck{
				TaskID: job.TaskID,
				Status: model.StatusRunning,
			})
			core.DeleteTaskContext(job.TaskID)
		}
	}
	return
}
*/

func CreateK8SJob(ctx context.Context, job *model.TaskDTO, lg *core.LogFormat, temporalWfID, temporalRunID string) (err error) {
	// TODO: fix this to fail by secondary files
	k8sjob, err := ConstructK8SJob(job, temporalWfID, temporalRunID)
	if err != nil {
		return err
	}
	err = core.CreateK8SJob(context.Background(), k8sjob)

	// turn on this when log is added
	//lg.Action = fmt.Sprintf("Created job %q.\n", result.GetObjectMeta().GetName())
	return
}

func ConstructK8SJob(job *model.TaskDTO, temporalWfID, temporalRunID string) (k8sJob *batchv1.Job, err error) {
	volumes, volumeMounts, err := ConstructMountPoint(job)
	if err != nil {
		return nil, err
	}
	args := Args(job)
	workingDir := WorkingDir(job)
	var dockerImage string
	var backoffLimit = int32(0)
	if len(job.DockerImage) != 0 {
		dockerImage = job.DockerImage[0]
	}
	mainConfig := core.GetMainConfig()
	nodeLabel := map[string]string{mainConfig.K8SConfig.NodeLabelKey: mainConfig.K8SConfig.NodeLabelValue}

	if !mainConfig.K8SConfig.DeleteJob {
		k8sJob = &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name: job.TaskID,
			},
			Spec: batchv1.JobSpec{
				//maybe add selector
				/*Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "demo",},},*/
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app":             job.TaskID,
							"task-uuid":       job.TaskUUID,
							"temporal-wf-id":  temporalWfID,
							"temporal-run-id": temporalRunID,
						},
					},
					Spec: apiv1.PodSpec{
						Volumes:      volumes,
						NodeSelector: nodeLabel,
						Containers: []apiv1.Container{
							{
								Name:            job.TaskID,
								Image:           dockerImage,
								ImagePullPolicy: "IfNotPresent",
								VolumeMounts:    volumeMounts,
								WorkingDir:      workingDir,

								Command: []string{defaultBashCommand},
								Args:    args,
							},
						},
						RestartPolicy: "Never",
					},
				},
				BackoffLimit: &backoffLimit,
			}}
	} else {
		ttl := int32(mainConfig.K8SConfig.JobTTLAfterFinished)
		k8sJob = &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name: job.TaskID,
			},
			Spec: batchv1.JobSpec{
				TTLSecondsAfterFinished: &ttl,
				//maybe add selector
				/*Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "demo",},},*/
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app":             job.TaskID,
							"task-uuid":       job.TaskUUID,
							"temporal-wf-id":  temporalWfID,
							"temporal-run-id": temporalRunID,
						},
					},
					Spec: apiv1.PodSpec{
						Volumes:      volumes,
						NodeSelector: nodeLabel,
						Containers: []apiv1.Container{
							{
								Name:            job.TaskID,
								Image:           dockerImage,
								ImagePullPolicy: "IfNotPresent",
								VolumeMounts:    volumeMounts,

								Command: []string{defaultBashCommand},
								Args:    args,
							},
						},
						RestartPolicy: "Never",
					},
				},
				BackoffLimit: &backoffLimit,
			}}
	}
	return
}

func CreateDirectory(path string) {
	if _, er := os.Stat(path); os.IsNotExist(er) {
		os.MkdirAll(path, 0777)
	}
	return
}

func ConstructMountPoint(job *model.TaskDTO) (volume []apiv1.Volume, volumeMount []apiv1.VolumeMount, err error) {
	dirType := apiv1.HostPathType("DirectoryOrCreate")
	//outdirType := apiv1.HostPathType("Directory")
	//fileType := apiv1.HostPathType("")
	var dirMap = make(map[string]bool)
	for i := 0; i < len(job.ParamsWithRegex); i++ {
		secondaryFile := job.ParamsWithRegex[i].SecondaryFiles

		//if secondaryFile != nil {}
		var counter []int
		var remainder []string
		for j1 := 0; j1 < len(secondaryFile); j1++ {
			count := strings.Count(secondaryFile[j1], "^")
			remain := strings.ReplaceAll(secondaryFile[j1], "^", "")
			counter = append(counter, count)
			remainder = append(remainder, remain)
		}

		/*
			for j := 0; j < len(job.ParamsWithRegex[i].Regex); j++ {
				if _, err := os.Stat(job.ParamsWithRegex[i].Regex[j]); err == nil {
					path := GetParentDirectory(GetFileFUSEPath(job.ParamsWithRegex[i].Regex[j]))
					if _, ok := dirMap[path]; ok {
						continue
					}
					dirMap[path] = true
					v := apiv1.Volume{
						Name: job.TaskID + "-regex-" + strconv.Itoa(i) + strconv.Itoa(j),
						VolumeSource: apiv1.VolumeSource{
							HostPath: &apiv1.HostPathVolumeSource{
								Path: path,
								Type: &dirType,
							},
						},
					}
					volume = append(volume, v)

					vm := apiv1.VolumeMount{
						//MountPath: job.ParamsWithRegex[i].Regex[j],
						MountPath: path,
						Name:      job.TaskID + "-regex-" + strconv.Itoa(i) + strconv.Itoa(j),
					}
					volumeMount = append(volumeMount, vm)
				}
			}
		*/

		for j := 0; j < len(job.ParamsWithRegex[i].Files); j++ {
			path := GetParentDirectory(GetFileFUSEPath(job.ParamsWithRegex[i].Files[j].Filepath))
			if path == "" {
				continue
			}

			filename := filepath.Base(job.ParamsWithRegex[i].Files[j].Filepath)
			for k := 0; k < len(counter); k++ {
				fileElement := strings.Split(filename, ".")
				var secondaryFileFirstPath string
				for k1 := 0; k1 < len(fileElement)-counter[k]; k1++ {
					if k1 != len(fileElement)-counter[k]-1 {
						secondaryFileFirstPath = secondaryFileFirstPath + fileElement[k1] + "."
					} else {
						secondaryFileFirstPath = secondaryFileFirstPath + fileElement[k1]
					}
				}
				secondaryFileName := path + secondaryFileFirstPath + remainder[k]
				if _, err := os.Stat(secondaryFileName); err != nil && os.IsNotExist(err) {
					return nil, nil, fmt.Errorf("Missing secondary file: %s", secondaryFileName)
				}
			}

			if _, ok := dirMap[path]; ok {
				continue
			}
			dirMap[path] = true
			v := apiv1.Volume{
				Name: job.TaskID + "-file-" + strconv.Itoa(i) + strconv.Itoa(j),
				VolumeSource: apiv1.VolumeSource{
					HostPath: &apiv1.HostPathVolumeSource{
						Path: path,
						Type: &dirType,
					},
				},
			}
			volume = append(volume, v)

			vm := apiv1.VolumeMount{
				//MountPath: job.ParamsWithRegex[i].Regex[j],
				MountPath: path,
				Name:      job.TaskID + "-file-" + strconv.Itoa(i) + strconv.Itoa(j),
			}
			volumeMount = append(volumeMount, vm)
		}
	}

	// Create volume for output dir
	outputPrefix := core.GetMainConfig().K8SConfig.OutputDirPrefix
	CreateDirectory(outputPrefix + "/" + job.TaskID)

	vout := apiv1.Volume{
		Name: job.TaskID + "-output",
		VolumeSource: apiv1.VolumeSource{
			HostPath: &apiv1.HostPathVolumeSource{
				Path: outputPrefix + "/" + job.TaskID,
				Type: &dirType,
			},
		},
	}

	vmout := apiv1.VolumeMount{
		MountPath: outputPrefix + "/" + job.TaskID,
		Name:      job.TaskID + "-output",
	}
	volume = append(volume, vout)
	volumeMount = append(volumeMount, vmout)

	return
}

func ConstructParam(job *executorModel.TaskHTTPResp) (args []string) {
	args = append(args, defaultBashParam)
	var commandArg string
	outputPrefix := core.GetMainConfig().K8SConfig.OutputDirPrefix
	commandArg = "cd " + outputPrefix + "/" + job.TaskID + " && " + job.Command + " "
	for i := 0; i < len(job.ParamsWithRegex); i++ {
		commandArg = commandArg + job.ParamsWithRegex[i].Prefix
		/*
			for j := 0; j < len(job.ParamsWithRegex[i].Regex); j++ {
				if _, err := os.Stat(job.ParamsWithRegex[i].Regex[j]); err == nil {
					commandArg = commandArg + job.ParamsWithRegex[i].Regex[j]
				}
				fmt.Println("regex")
				fmt.Println(job.ParamsWithRegex[i].Regex[j])
			}
		*/
		for j := 0; j < len(job.ParamsWithRegex[i].Files); j++ {
			commandArg = commandArg + job.ParamsWithRegex[i].Files[j].Filepath
		}
	}
	args = append(args, commandArg)
	return
}

func WorkingDir(job *model.TaskDTO) (workingDir string) {
	workingDir = core.GetMainConfig().K8SConfig.OutputDirPrefix + "/" + job.TaskID
	return
}

func Args(job *model.TaskDTO) (args []string) {
	args = append(args, defaultBashParam)
	command := []string{job.Command}

	for i := 0; i < len(job.ParamsWithRegex); i++ {
		commandArg := job.ParamsWithRegex[i].Prefix
		/*
			var isEqualCommand bool = false
			if len(commandArg) > 0 {
				isEqualCommand = commandArg[len(commandArg)-1] == '='
			} else {
				isEqualCommand = true
			}
		*/

		for j := 0; j < len(job.ParamsWithRegex[i].Files); j++ {
			/*
				if j == 0 {
					if isEqualCommand {
						commandArg = commandArg + job.ParamsWithRegex[i].Files[j].Filepath
					} else {
						commandArg = commandArg + " " + job.ParamsWithRegex[i].Files[j].Filepath
					}
				} else {
					commandArg = commandArg + " " + job.ParamsWithRegex[i].Files[j].Filepath
				}
			*/

			commandArg = commandArg + job.ParamsWithRegex[i].Files[j].Filepath
		}
		command = append(command, commandArg)
	}

	args = append(args, strings.Join(command, ""))
	return
}
