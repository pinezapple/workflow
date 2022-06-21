package logger

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// LogFormat log pattern for whole system
type LogFormat struct {
	ServiceName string      `json:"srv"`
	Source      string      `json:"src,omitempty"`
	Action      string      `json:"act,omitempty"`
	Data        interface{} `json:"dat,omitempty"`   // data
	Err         interface{} `json:"err,omitempty"`   // error
	Success     interface{} `json:"suc,omitempty"`   // success
	Stack       interface{} `json:"stack,omitempty"` // stack trace
	Message     string      `json:"msg,omitempty"`
}

// ToMapStringItf ...
func (lg *LogFormat) ToMapStringItf() map[string]interface{} {
	return map[string]interface{}{
		"srv":   lg.ServiceName,
		"src":   lg.Source,
		"act":   lg.Action,
		"dat":   lg.Data,
		"err":   lg.Err,
		"suc":   lg.Success,
		"stack": lg.Stack,
	}
}

func CreateNewLogObject(name string) (lg *LogFormat) {
	lg = &LogFormat{ServiceName: name}
	return lg
}

func (lg *LogFormat) LogData(data interface{}) {
	if lg == nil {
		return
	}

	lg.Data = data
	js, _ := json.Marshal(lg)
	log.Infof("%s", js)
	lg.Data = ""
}

// LogInfo information logging
func (lg *LogFormat) LogInfo(message string) {
	if lg == nil {
		return
	}

	lg.Action = message
	js, _ := json.Marshal(lg)
	log.Infof("%s", js)
	lg.Action = ""
}

// LogErr error logging
func (lg *LogFormat) LogErr(err error) {
	if err == nil {
		return
	}
	lg.Err = err.Error()
	js, _ := json.Marshal(lg)
	log.Errorf("%s", js)
	lg.Err = ""
}

// LogWarning warning logging
func (lg *LogFormat) LogWarning() {
	if lg == nil {
		return
	}

	js, _ := json.Marshal(lg)
	log.Infof("%s", js)
}
