package es

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/bytedance/sonic"
	"github.com/olivere/elastic/v7"
	"github.com/yxrxy/videoHub/app/video/domain/model"
	"github.com/yxrxy/videoHub/pkg/constants"
	"github.com/yxrxy/videoHub/pkg/errno"
)

func (es *VideoElastic) IsExist(ctx context.Context, indexName string) bool {
	res, err := es.client.IndexExists(indexName).Do(ctx)
	if err != nil {
		logger.Errorf("VideoElastic.IsExist Error checking if index exists: %v", err)
		return false
	}
	return res
}

func (es *VideoElastic) CreateIndex(ctx context.Context, indexName string) error {
	_, err := es.client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, "VideoElastic.CreateIndex Error creating index: %v", err)
	}
	return nil
}

func (es *VideoElastic) AddItem(ctx context.Context, indexName string, video *model.Video, name string) error {
	VideoES := &model.VideoES{
		ID:          video.ID,
		Name:        name,
		Title:       video.Title,
		Description: video.Description,
		Tags:        strings.Split(video.Tags, ","),
		Category:    video.Category,
		AuthorID:    video.UserID,
		CreatedAt:   time.Now(),
		ViewCount:   video.VisitCount,
		IsDeleted:   false,
		SearchText:  fmt.Sprintf("%s %s %s", video.Title, video.Description, video.Tags),
	}

	_, err := es.client.Index().Index(indexName).
		Id(strconv.FormatInt(VideoES.ID, 10)).
		BodyJson(VideoES).Do(ctx)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, "VideoElastic.AddItem Error adding item: %v", err)
	}
	return nil
}

func (es *VideoElastic) RemoveItem(ctx context.Context, indexName string, id int64) error {
	_, err := es.client.Delete().Index(indexName).Id(fmt.Sprintf("%d", id)).Do(ctx)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, "VideoElastic.RemoveItem failed: %v", err)
	}

	return nil
}

func structToMapUsingJSON(obj interface{}) map[string]interface{} {
	data, _ := sonic.Marshal(obj)
	var result map[string]interface{}
	_ = sonic.Unmarshal(data, &result)
	return result
}

func (es *VideoElastic) UpdateItem(ctx context.Context, indexName string, video *model.VideoES, name string) error {
	videoEs := &model.VideoES{
		ID:          video.ID,
		Name:        name,
		Title:       video.Title,
		Description: video.Description,
		Tags:        video.Tags,
		Category:    video.Category,
		AuthorID:    video.AuthorID,
		ViewCount:   video.ViewCount,
		SearchText:  fmt.Sprintf("%s %s %s", video.Title, video.Description, strings.Join(video.Tags, " ")),
	}
	_, err := es.client.Update().Index(indexName).
		Id(strconv.FormatInt(video.ID, 10)).
		Doc(structToMapUsingJSON(videoEs)).
		Do(ctx)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, "VideoElastic.UpdateItem failed: %v", err)
	}
	return nil
}

func (es *VideoElastic) SearchItems(ctx context.Context, indexName string, query *model.VideoES) ([]int64, int64, error) {
	q := es.BuildQuery(query)

	result, err := es.client.Search().Index(indexName).
		Query(q).
		Do(ctx)
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalESErrorCode, "VideoElastic.SearchItems failed: %v", err)
	}

	rets := make([]int64, 0)
	for _, hit := range result.Hits.Hits {
		var videoEs model.VideoES
		data, err := hit.Source.MarshalJSON()
		if err != nil {
			return nil, 0, errno.Errorf(errno.InternalServiceErrorCode, "VideoElastic.SearchItems failed: %v", err)
		}
		err = sonic.Unmarshal(data, &videoEs)
		if err != nil {
			return nil, 0, errno.Errorf(errno.InternalServiceErrorCode, "VideoElastic.SearchItems failed: %v", err)
		}
		rets = append(rets, videoEs.ID)
	}

	return rets, result.TotalHits(), nil
}

func (es *VideoElastic) BuildQuery(req *model.VideoES) *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	hasCondition := false
	if req.Keywords != "" {
		matchQuery := elastic.NewMatchQuery("search_text", req.Keywords)
		query = query.Must(matchQuery)
		hasCondition = true
	}
	if req.FromDate != nil || req.ToDate != nil {
		dateRange := elastic.NewRangeQuery("created_at")
		if req.FromDate != nil {
			dateRange.Gte(time.Unix(*req.FromDate/constants.MillisecondsPerSecond, 0))
		}
		if req.ToDate != nil {
			dateRange.Lte(time.Unix(*req.ToDate/constants.MillisecondsPerSecond, 0))
		}
		query = query.Must(dateRange)
		hasCondition = true
	}
	if req.Username != nil && *req.Username != "" {
		query = query.Must(elastic.NewTermQuery("name", *req.Username))
		hasCondition = true
	}

	query = query.MustNot(elastic.NewTermQuery("is_deleted", true))

	if !hasCondition {
		query = query.Must(elastic.NewMatchAllQuery())
	}

	return query
}
