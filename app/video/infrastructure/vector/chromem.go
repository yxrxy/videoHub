package vector

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/philippgille/chromem-go"
	"github.com/yxrxy/videoHub/app/video/domain/model"
	"github.com/yxrxy/videoHub/app/video/domain/repository"
)

type ChromemDB struct {
	db         *chromem.DB
	collection *chromem.Collection
}

func NewChromemDB(collectionName string) (repository.VectorDB, error) {
	db := chromem.NewDB()

	// 创建集合，使用默认的OpenAI嵌入(或者您可以指定其他嵌入函数)
	// 注意：这会需要设置OPENAI_API_KEY环境变量
	collection, err := db.CreateCollection(collectionName, nil, nil)
	if err != nil {
		return nil, err
	}

	return &ChromemDB{
		db:         db,
		collection: collection,
	}, nil
}

func (c *ChromemDB) StoreVector(ctx context.Context, videoID int64, vector []float32, metadata *model.VideoMetadata) error {
	// 准备元数据，string类型
	metadataMap := map[string]string{
		"title":       metadata.Title,
		"description": metadata.Description,
		"category":    metadata.Category,
		"user_id":     strconv.FormatInt(metadata.UserID, 10),
	}

	// 标签处理
	if len(metadata.Tags) > 0 {
		metadataMap["tags"] = strings.Join(metadata.Tags, ",")
	}

	// 内容文本（用于展示）
	content := fmt.Sprintf("%s %s %s",
		metadata.Title,
		metadata.Description,
		strings.Join(metadata.Tags, " "))

	// 创建文档（不需要手动传入向量）
	doc := chromem.Document{
		ID:        strconv.FormatInt(videoID, 10),
		Content:   content,
		Metadata:  metadataMap,
		Embedding: vector, // 直接设置向量
	}

	return c.collection.AddDocument(ctx, doc)
}

func (c *ChromemDB) SearchSimilar(ctx context.Context, queryVector []float32, limit int32, filter *model.VectorSearchFilter) ([]int64, []float32, error) {
	// 准备过滤器
	var metadataFilter map[string]string
	if filter != nil && filter.Category != nil {
		metadataFilter = map[string]string{
			"category": *filter.Category,
		}
	}

	// 创建查询选项
	options := chromem.QueryOptions{
		QueryEmbedding: queryVector,
		NResults:       int(limit),
		Where:          metadataFilter,
	}

	// 执行查询
	results, err := c.collection.QueryWithOptions(ctx, options)
	if err != nil {
		return nil, nil, err
	}

	// 解析结果
	ids := make([]int64, 0, len(results))
	scores := make([]float32, 0, len(results))

	for _, result := range results {
		id, err := strconv.ParseInt(result.ID, 10, 64)
		if err != nil {
			return nil, nil, err
		}
		ids = append(ids, id)
		scores = append(scores, result.Similarity)
	}

	return ids, scores, nil
}

func (c *ChromemDB) DeleteEmbedding(ctx context.Context, videoID int64) error {
	//where: Conditional filtering on metadata. Optional.
	//whereDocument: Conditional filtering on documents. Optional.
	//ids: The ids of the documents to delete. If empty, all documents are deleted.
	return c.collection.Delete(ctx, nil, nil, strconv.FormatInt(videoID, 10))
}
