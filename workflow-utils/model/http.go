package model

// ResponseError response error struct
type ResponseError struct {
	Code    int
	Message string
}

// Response response struct
type Response struct {
	Error ResponseError
	Data  interface{}
}

// SetData set data attached to response
func (c *Response) SetData(_dat interface{}) {
	c.Data = _dat
}

// SetCodeMessage set code message
func (c *Response) SetCodeMessage(code int, message string) {
	c.Error.Code = code
	c.Error.Message = string(message)
}
