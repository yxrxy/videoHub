package es

import (
	"github.com/olivere/elastic/v7"

	"github.com/yxrxy/videoHub/app/video/domain/repository"
)

type VideoElastic struct {
	client *elastic.Client
}

func NewVideoElastic(client *elastic.Client) repository.VideoElastic {
	return &VideoElastic{client: client}
}
