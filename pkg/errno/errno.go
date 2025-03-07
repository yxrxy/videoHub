package errno

import "fmt"

type ErrNo struct {
	ErrCode int32
	ErrMsg  string
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

var (
	// 成功
	Success = ErrNo{0, "success"}

	// 用户模块错误: 100xx
	ErrUserNotExist     = ErrNo{10001, "用户不存在"}
	ErrUserAlreadyExist = ErrNo{10002, "用户已存在"}
	ErrPasswordWrong    = ErrNo{10003, "密码错误"}
	ErrInvalidToken     = ErrNo{10004, "无效的令牌"}
	ErrUnauthorized     = ErrNo{10005, "未授权访问"}
)
