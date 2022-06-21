package infra

type MailVariables struct {
	Name  string
	Value interface{}
}

type CreateMailReq struct {
	FromID                int              `json:"from_id"`
	Subject               string           `json:"subject"`
	TemplateNameAsSubject bool             `json:"template_name_as_subject"`
	TemplateID            int              `json:"template_id"`
	To                    []string         `json:"to"`
	Vars                  []*MailVariables `json:"variables"`
}

type CreateMailResp struct {
	errorResp
	MailID int `json:"mail_id"`
}

type SendMailReq struct {
	MailID int `json:"mailID"`
}

type SendMailResp struct {
	errorResp
}

type errorResp struct {
	Error struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
	} `json:"error"`
}
