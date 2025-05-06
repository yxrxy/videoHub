package model

type VideoMetadata struct {
	Title       string
	Description string
	Tags        []string
	Category    string
	UserID      int64
}

type VectorSearchFilter struct {
	Category *string
	FromDate *int64
	ToDate   *int64
	UserID   *int64
}

type RAGResponse struct {
	Summary        string   `json:"summary"`
	Keywords       []string `json:"keywords"`
	RelatedQueries []string `json:"related_queries"`
}
