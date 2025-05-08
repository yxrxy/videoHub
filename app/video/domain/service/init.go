package service

import (
	"context"

	"github.com/yxrxy/videoHub/app/user/domain/repository"
	videorepo "github.com/yxrxy/videoHub/app/video/domain/repository"
)

type VideoService struct {
	db        videorepo.VideoDB          // 视频数据库接口
	cache     videorepo.VideoCache       // 视频缓存接口
	mq        videorepo.VideoMQ          // 视频消息队列接口
	es        videorepo.VideoElastic     // 视频搜索引擎接口
	userDB    repository.UserDB          // 用户数据库接口
	embedding videorepo.EmbeddingService // 向量嵌入服务接口
	vectorDB  videorepo.VectorDB         // 向量数据库接口
	llm       videorepo.LLMService       // 大语言模型服务接口
}

func NewVideoService(
	db videorepo.VideoDB,
	cache videorepo.VideoCache,
	mq videorepo.VideoMQ,
	es videorepo.VideoElastic,
	userDB repository.UserDB,
	embedding videorepo.EmbeddingService,
	vectorDB videorepo.VectorDB,
	llm videorepo.LLMService) *VideoService {
	if db == nil || cache == nil || mq == nil || es == nil || userDB == nil {
		panic("videoService`s db or cache or mq or es or userDB should not be nil")
	}
	svc := &VideoService{
		db:        db,
		cache:     cache,
		mq:        mq,
		es:        es,
		userDB:    userDB,
		embedding: embedding,
		vectorDB:  vectorDB,
		llm:       llm,
	}
	svc.init()
	return svc
}

func (s *VideoService) init() {
	s.initConsumer()
}

func (s *VideoService) initConsumer() {
	go s.ConsumeProcessVideo(context.Background())
}
