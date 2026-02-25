package elastic

import (
	"context"

	elasticEntity "gold-gym-be/internal/entity/elastic"
)

// IndexUser indexes a user document into Elasticsearch
func (s *Service) IndexUser(ctx context.Context, index string, doc elasticEntity.UserDocument) (string, error) {
	return s.elastic.IndexDocument(ctx, index, doc)
}

// SearchUsers performs a full-text search across user documents
func (s *Service) SearchUsers(ctx context.Context, index string, query string) ([]elasticEntity.UserDocument, error) {
	return s.elastic.SearchDocuments(ctx, index, query)
}

// GetUserByID retrieves a single user document by ID
func (s *Service) GetUserByID(ctx context.Context, index string, id string) (elasticEntity.UserDocument, error) {
	return s.elastic.GetDocumentByID(ctx, index, id)
}
