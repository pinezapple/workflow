package responsedto

type RespDto struct {
	Error ErrRespDto  `json:"error"`
	Data  interface{} `json:"data"`
}

type ErrRespDto struct {
	Code    int64  `json:"code"`
	Message string `json:"message"` // empty if success 200
}

func (resp *RespDto) SetCodeMessage(code int64, message string) {
	resp.Error.Code = code
	resp.Error.Message = message
}

func (resp *RespDto) SetData(data interface{}) {
	if data == nil {
		return
	}

	resp.Data = data
}
