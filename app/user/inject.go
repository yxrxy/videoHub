package user

import (
	"github.com/yxrxy/videoHub/app/user/controllers/rpc"
	"github.com/yxrxy/videoHub/app/user/domain/service"
	"github.com/yxrxy/videoHub/app/user/infrastructure/cache"
	"github.com/yxrxy/videoHub/app/user/infrastructure/mysql"
	"github.com/yxrxy/videoHub/app/user/usecase"
	"github.com/yxrxy/videoHub/config"
	"github.com/yxrxy/videoHub/kitex_gen/user"
	"github.com/yxrxy/videoHub/pkg/base/client"
)

func InjectUserHandler() user.UserService {
	gormDB, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}

	re, err := client.NewRedisClient(config.Redis.DB.User)
	if err != nil {
		panic(err)
	}

	db := mysql.NewUserDB(gormDB)
	redisCache := cache.NewUserCache(re)
	svc := service.NewUserService(db, redisCache)
	uc := usecase.NewUserCase(db, svc)

	return rpc.NewUserHandler(uc)
}
