package response

import "github.com/yxrrxy/videoHub/pkg/errno"

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
			Code:    errno.Success.ErrCode,
			Message: errno.Success.ErrMsg,
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
