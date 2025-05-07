package video

import (
	usermysql "github.com/yxrxy/videoHub/app/user/infrastructure/mysql"
	"github.com/yxrxy/videoHub/app/video/controllers/rpc"
	"github.com/yxrxy/videoHub/app/video/domain/service"
	videocache "github.com/yxrxy/videoHub/app/video/infrastructure/cache"
	"github.com/yxrxy/videoHub/app/video/infrastructure/embedding"
	"github.com/yxrxy/videoHub/app/video/infrastructure/es"
	"github.com/yxrxy/videoHub/app/video/infrastructure/llm"
	"github.com/yxrxy/videoHub/app/video/infrastructure/mq"
	videomysql "github.com/yxrxy/videoHub/app/video/infrastructure/mysql"
	"github.com/yxrxy/videoHub/app/video/infrastructure/vector"
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
	re, err := client.NewRedisClient(config.Redis.DB.Video)
	if err != nil {
		panic(err)
	}
	elastic, err := client.NewEsVideoClient()
	if err != nil {
		panic(err)
	}

	redisCache := videocache.NewVideoCache(re)
	kafMQ := kafka.NewKafkaInstance()
	kaf := mq.NewVideoMQ(kafMQ)
	db := videomysql.NewVideoDB(gormDB, redisCache)
	esClient := es.NewVideoElastic(elastic)
	userDB := usermysql.NewUserDB(gormDB)
	emb := embedding.NewOpenAIEmbedding(config.ApiKey.Key, config.ApiKey.BaseURL, config.ApiKey.Proxy)
	vec, _ := vector.NewChromemDB("videos")
	llm := llm.NewOpenAILLM(config.ApiKey.Key, config.ApiKey.BaseURL, config.ApiKey.Proxy)
	svc := service.NewVideoService(db, redisCache, kaf, esClient, userDB, emb, vec, llm)
	uc := usecase.NewVideoCase(db, redisCache, esClient, svc)
	return rpc.NewVideoHandler(uc)
}
