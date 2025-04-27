package video

import (
	"github.com/yxrxy/videoHub/app/video/controllers/rpc"
	"github.com/yxrxy/videoHub/app/video/domain/service"
	"github.com/yxrxy/videoHub/app/video/infrastructure/cache"
	"github.com/yxrxy/videoHub/app/video/infrastructure/mq"
	"github.com/yxrxy/videoHub/app/video/infrastructure/mysql"
	"github.com/yxrxy/videoHub/app/video/usecase"
	"github.com/yxrxy/videoHub/config"
	"github.com/yxrxy/videoHub/kitex_gen/video"
	"github.com/yxrxy/videoHub/pkg/base/client"
	"github.com/yxrxy/videoHub/pkg/kafka"
)

func InjectVideoHandler() video.VideoService {
	gormDB, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}
	re, err := client.NewRedisClient(config.Redis.DB)
	if err != nil {
		panic(err)
	}
	redisCache := cache.NewVideoCache(re)
	kafMQ := kafka.NewKafkaInstance()
	kaf := mq.NewVideoMQ(kafMQ)
	db := mysql.NewVideoDB(gormDB, redisCache)
	svc := service.NewVideoService(db, redisCache, kaf)
	uc := usecase.NewVideoCase(db, redisCache, svc)

	return rpc.NewVideoHandler(uc)
}
