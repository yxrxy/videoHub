package es

var mapping = `{
  "settings": {
    "number_of_shards": 3,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "text_analyzer": {
          "tokenizer": "ik_max_word",
          "filter": ["lowercase", "asciifolding", "trim"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": { "type": "long" },
      "name": { "type": "keyword" },
      "title": {
        "type": "text",
        "analyzer": "text_analyzer",
        "fields": {
          "keyword": { "type": "keyword" }
        }
      },
      "description": {
        "type": "text",
        "analyzer": "text_analyzer"
      },
      "tags": { "type": "keyword" },
      "category": { "type": "keyword" },
      "author_id": { "type": "long" },
      "created_at": { "type": "date" },
      "view_count": { "type": "long" },
      "is_deleted": { "type": "boolean" },
      "search_text": {
        "type": "text",
        "analyzer": "text_analyzer"
      }
    }
  }
}`
