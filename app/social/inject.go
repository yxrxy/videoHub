package social

import (
	"github.com/yxrxy/videoHub/app/social/controllers/rpc"
	"github.com/yxrxy/videoHub/app/social/domain/service"
	"github.com/yxrxy/videoHub/app/social/infrastructure/cache"
	"github.com/yxrxy/videoHub/app/social/infrastructure/mysql"
	"github.com/yxrxy/videoHub/app/social/usecase"
	"github.com/yxrxy/videoHub/config"
	"github.com/yxrxy/videoHub/kitex_gen/social"
	"github.com/yxrxy/videoHub/pkg/base/client"
)

func InjectSocialHandler() social.SocialService {
	gormDB, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}

	re, err := client.NewRedisClient(config.Redis.DB.Social)
	if err != nil {
		panic(err)
	}
	db := mysql.NewSocialDB(gormDB)
	cache0 := cache.NewSocialCache(re)

	svc := service.NewSocialService(db, cache0)
	uc := usecase.NewSocialCase(db, cache0, svc)

	return rpc.NewSocialHandler(uc)
}
