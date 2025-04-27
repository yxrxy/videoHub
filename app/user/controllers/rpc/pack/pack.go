package pack

import (
	"github.com/yxrxy/videoHub/app/user/domain/model"
	rpcmodel "github.com/yxrxy/videoHub/kitex_gen/model"
)

// User 将数据库实体转换为RPC响应实体
func User(u *model.User) *rpcmodel.User {
	if u == nil {
		return nil
	}
	return &rpcmodel.User{
		Id:       u.ID,
		Username: u.Username,
		Avatar:   &u.AvatarURL,
	}
}
