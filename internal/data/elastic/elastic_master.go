package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	elasticEntity "gold-gym-be/internal/entity/elastic"
)

// IndexDocument indexes a UserDocument into the given ES index
func (r *Repository) IndexDocument(ctx context.Context, index string, doc elasticEntity.UserDocument) (string, error) {
	doc.IndexedAt = time.Now().UTC().Format(time.RFC3339)

	body, err := json.Marshal(doc)
	if err != nil {
		return "", fmt.Errorf("marshal document: %w", err)
	}

	docID := fmt.Sprintf("%d", doc.GoldId)
	res, err := r.client.Index(
		index,
		bytes.NewReader(body),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithDocumentID(docID),
	)
	if err != nil {
		return "", fmt.Errorf("es index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		b, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("es index error [%s]: %s", res.Status(), string(b))
	}

	return docID, nil
}

// SearchDocuments performs a full-text match search in the given index
func (r *Repository) SearchDocuments(ctx context.Context, index string, query string) ([]elasticEntity.UserDocument, error) {
	queryBody := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"gold_nama", "gold_email", "gold_nomorhp"},
			},
		},
	}

	body, err := json.Marshal(queryBody)
	if err != nil {
		return nil, fmt.Errorf("marshal query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(index),
		r.client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("es search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("es search error [%s]: %s", res.Status(), string(b))
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source elasticEntity.UserDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode search result: %w", err)
	}

	docs := make([]elasticEntity.UserDocument, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		docs = append(docs, hit.Source)
	}

	return docs, nil
}

// GetDocumentByID fetches a single document by its ID from the given ES index
func (r *Repository) GetDocumentByID(ctx context.Context, index string, id string) (elasticEntity.UserDocument, error) {
	res, err := r.client.Get(
		index,
		id,
		r.client.Get.WithContext(ctx),
	)
	if err != nil {
		return elasticEntity.UserDocument{}, fmt.Errorf("es get: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		b, _ := io.ReadAll(res.Body)
		return elasticEntity.UserDocument{}, fmt.Errorf("es get error [%s]: %s", res.Status(), string(b))
	}

	var result struct {
		Source elasticEntity.UserDocument `json:"_source"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return elasticEntity.UserDocument{}, fmt.Errorf("decode get result: %w", err)
	}

	return result.Source, nil
}
