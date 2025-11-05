package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"github.com/elastic/go-elasticsearch/v8"
)

type MovieElasticRepository struct {
	client *elasticsearch.Client
}

func NewMovieElasticRepository(client *elasticsearch.Client) *MovieElasticRepository {
	return &MovieElasticRepository{client: client}
}

const movieIndex = "movie"

func (r *MovieElasticRepository) Create(movie *models.MovieElastic) error {
	ctx := context.Background()

	data, err := json.Marshal(movie)
	if err != nil {
		return err
	}

	res, err := r.client.Index(
		movieIndex,
		bytes.NewReader(data),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithDocumentID(movie.ID),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to index document: %s", res.String())
	}
	return nil
}

func (r *MovieElasticRepository) CreateBulk(movies []*models.MovieElastic) error {
	ctx := context.Background()

	var buf bytes.Buffer
	for _, movie := range movies {
		data, err := json.Marshal(movie)
		if err != nil {
			continue
		}
		buf.WriteString(fmt.Sprintf(`{"index":{"_id":"%s"}}%s`, movie.ID, "\n"))
		buf.Write(data)
		buf.WriteString("\n")
	}

	if buf.Len() == 0 {
		return nil
	}

	res, err := r.client.Bulk(
		bytes.NewReader(buf.Bytes()),
		r.client.Bulk.WithContext(ctx),
		r.client.Bulk.WithIndex(movieIndex),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk indexing failed: %s", res.String())
	}
	return nil
}

func (r *MovieElasticRepository) FindByID(id string) (*models.MovieElastic, error) {
	ctx := context.Background()

	res, err := r.client.Get(
		movieIndex,
		id,
		r.client.Get.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("document not found: %s", res.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	source, ok := result["_source"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid source")
	}

	movie := &models.MovieElastic{}
	if v, ok := source["id"].(string); ok {
		movie.ID = v
	}
	if v, ok := source["name"].(string); ok {
		movie.Name = v
	}
	if v, ok := source["movieType"].(string); ok {
		movie.MovieType = v
	}

	return movie, nil
}

func (r *MovieElasticRepository) FindAll() ([]models.MovieElastic, error) {
	ctx := context.Background()

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(movieIndex),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search failed: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits, ok := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, nil
	}

	var movies []models.MovieElastic
	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})

		movie := models.MovieElastic{}
		if v, ok := source["id"].(string); ok {
			movie.ID = v
		}
		if v, ok := source["name"].(string); ok {
			movie.Name = v
		}
		if v, ok := source["movieType"].(string); ok {
			movie.MovieType = v
		}
		if v, ok := source["persona"].([]interface{}); ok {
			for _, p := range v {
				pMap := p.(map[string]interface{})
				persona := models.Persona{}
				if name, ok := pMap["name"].(string); ok {
					persona.Name = name
				}
				if ptype, ok := pMap["type"].(string); ok {
					persona.Type = ptype
				}
				movie.Persona = append(movie.Persona, persona)
			}
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

func (r *MovieElasticRepository) SearchByText(query string) ([]models.MovieElastic, error) {
	ctx := context.Background()

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{"wildcard": map[string]interface{}{"name": map[string]interface{}{"value": "*" + query + "*"}}},
					{"wildcard": map[string]interface{}{"movieType": map[string]interface{}{"value": "*" + query + "*"}}},
					{"wildcard": map[string]interface{}{"persona.name": map[string]interface{}{"value": "*" + query + "*"}}},
				},
				"minimum_should_match": 0,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(movieIndex),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search failed: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits, ok := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, nil
	}

	var movies []models.MovieElastic
	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})

		movie := models.MovieElastic{}
		if v, ok := source["id"].(string); ok {
			movie.ID = v
		}
		if v, ok := source["name"].(string); ok {
			movie.Name = v
		}
		if v, ok := source["movieType"].(string); ok {
			movie.MovieType = v
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

func (r *MovieElasticRepository) DeleteAll() error {
	ctx := context.Background()

	_, err := r.client.Delete(
		movieIndex,
		"_all",
		r.client.Delete.WithContext(ctx),
	)
	// Ignore error if index doesn't exist
	return err
}

func (r *MovieElasticRepository) DeleteByID(id string) error {
	ctx := context.Background()

	_, err := r.client.Delete(
		movieIndex,
		id,
		r.client.Delete.WithContext(ctx),
	)
	return err
}