package infra

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	mailServiceTimeOut  = time.Second * 10
	emailRegex          = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
	notifyEmailTemplate = 12
	notifyEmailSubject  = "Your run has ended"
)

var (
	ErrNotStausOK         = errors.New("not status ok")
	ErrInvalidEmailFormat = errors.New("invalid email format")
)

type MailService struct {
	address          string
	accID            int
	httpClient       *http.Client
	emailRegexp      *regexp.Regexp
	notifyTemplateID int
}

var (
	globalMailService = new(MailService)
)

func ReplaceGlobalMailSrv(srv *MailService) {
	globalMailService = srv
}

func GetMailSrv() *MailService {
	return globalMailService
}

func NewMailService(address string, accID, notifyTemplateID int) (*MailService, error) {
	emailRegexp, err := regexp.Compile(emailRegex)
	if err != nil {
		return nil, err
	}

	return &MailService{
		address: address,
		accID:   accID,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: mailServiceTimeOut,
		},
		emailRegexp:      emailRegexp,
		notifyTemplateID: notifyTemplateID,
	}, nil
}

func (m MailService) formCreateURL() string {
	return m.address + "/internal/mail"
}

func (m MailService) formSendURL() string {
	return m.address + "/internal/mail/send"
}

func (m MailService) post(ctx context.Context, url string, reqBody, respExpected interface{}) error {
	if reqBody == nil {
		return errors.New("reqBody is nil")
	}

	reqByte, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqByte))
	if err != nil {
		return err
	}

	resp, err := m.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if respExpected != nil {
		if err := json.NewDecoder(resp.Body).Decode(respExpected); err != nil {
			return err
		}
	}

	if http.StatusOK != http.StatusOK {
		return ErrNotStausOK
	}

	return nil
}

func (m MailService) newMail(ctx context.Context, to, subject string, templateID int, dataOpt ...MailServiceDataOpt) (int, error) {
	// validate to
	if !m.emailRegexp.MatchString(to) {
		return 0, ErrInvalidEmailFormat
	}

	var (
		mailVars = make([]*MailVariables, 0, len(dataOpt))
	)

	for i := range dataOpt {
		dataOpt[i].apply(&mailVars)
	}

	var (
		reqBody = &CreateMailReq{
			FromID:                m.accID,
			Subject:               subject,
			TemplateNameAsSubject: true,
			TemplateID:            templateID,
			To:                    []string{to},
			Vars:                  mailVars,
		}
		resp = new(CreateMailResp)
	)

	if err := m.post(ctx, m.formCreateURL(), reqBody, resp); err != nil {
		return 0, err
	}

	return resp.MailID, nil
}

func (m MailService) sendMail(ctx context.Context, mailID int) error {
	var (
		reqBody = &SendMailReq{
			MailID: mailID,
		}
		resp = new(SendMailResp)
	)

	if err := m.post(ctx, m.formSendURL(), reqBody, resp); err != nil {
		return fmt.Errorf("send email error: %v, msg: %s", err, resp.Error.Msg)
	}

	return nil
}

func (m MailService) SendNotifyEmail(ctx context.Context, email, runID, workflowName, endStatus, taskUUID, taskName string) error {
	// handle the workflow name
	workflowName = strings.Replace(workflowName, ".cwl", "", -1)
	switch {
	case strings.Contains(strings.ToLower(workflowName), "bionet"):
		workflowName = "WES - (without capture kit)"

	case strings.Contains(strings.ToLower(workflowName), "wes"):
		workflowName = "WES - (with capture kit)"
	}

	// generate the taskName
	taskNameList := strings.Split(taskName, "-")
	taskName = taskNameList[len(taskNameList)-1]

	// generate the status color
	var statusColor string
	switch endStatus {
	case "COMPLETED":
		statusColor = "green"
	default:
		statusColor = "red"
	}

	// create a new email
	newMailID, err := m.newMail(ctx, email, notifyEmailSubject, m.notifyTemplateID,
		AddData("RunID", runID),
		AddData("WorkflowName", workflowName),
		AddData("EndStatus", endStatus),
		AddData("TaskUUID", taskUUID),
		AddData("TaskName", taskName),
		AddData("StatusColor", statusColor),
	)
	if err != nil {
		return err
	}

	// send the email
	if err := m.sendMail(ctx, newMailID); err != nil {
		return err
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- DATA OPTIONS ----------------------------------------------------------
type MailServiceDataOpt interface {
	apply(vars *[]*MailVariables)
}

type mailServiceDataOptFn func(vars *[]*MailVariables)

func (fn mailServiceDataOptFn) apply(vars *[]*MailVariables) {
	fn(vars)
}

func AddData(key string, value interface{}) MailServiceDataOpt {
	return mailServiceDataOptFn(func(vars *[]*MailVariables) {
		if vars == nil {
			return
		}
		*vars = append(*vars, &MailVariables{
			Name:  key,
			Value: value,
		})
	})
}
