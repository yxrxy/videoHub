package response

type Response struct {
	Base struct {
		Code    int32  `json:"code"`
		Message string `json:"msg"`
	} `json:"base"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *Response {
	return &Response{
		Base: struct {
			Code    int32  `json:"code"`
			Message string `json:"msg"`
		}{
			Code:    0,
			Message: "success",
		},
		Data: data,
	}
}

func Error(code int32, msg string) *Response {
	return &Response{
		Base: struct {
			Code    int32  `json:"code"`
			Message string `json:"msg"`
		}{
			Code:    code,
			Message: msg,
		},
		Data: nil,
	}
}
