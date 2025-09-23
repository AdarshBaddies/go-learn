package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	elastic "github.com/elastic/go-elasticsearch/v8"
)

var (
	ErrNotFound = errors.New("Entity not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       float64 `json:"price"`
}

type SearchResponse struct{
	Hits struct{
		Total struct{
			Value 	int 			`json:"value"`
		}`json:"total"`
		Hits []struct{
			Id 		string 			`json:"_id"`
			Source 	productDocument `json:"_source"`
		}`json:"hits"`
	}`json:"hits"`
}

type MultiGetResponse struct {
	Docs []struct {
		ID     string `json:"_id"`
		Source *productDocument `json:"_source"`
		Found  bool   `json:"found"`
	} `json:"docs"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(
		elastic.Config{
			Addresses: []string{url},
		},
	)
	if err != nil {
		return nil, err
	}

	res, err := client.Info()
	if err != nil {
		return nil, errors.New("failed to connect to Elasticsearch: " + err.Error())
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New("elasticsearch returned error: " + res.String())
	}

	return &elasticRepository{client}, nil
}

func (r *elasticRepository) Close() {
}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error {
	doc := productDocument{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}

	body, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	res, err := r.client.Index(
		"catalog",
		bytes.NewReader(body),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithDocumentID(p.ID),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New("failed to put product: " + res.Status())
	}

	return nil
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.GetSource(
		"catalog",
		id,
		r.client.GetSource.WithContext(ctx),
	)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	p := productDocument{}
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}

	return &Product{
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error){
	res, err := r.client.Search(
		r.client.Search.WithIndex("catalog"),
		r.client.Search.WithFrom(int(skip)),
		r.client.Search.WithSize(int(take)),
		r.client.Search.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	searchRes := SearchResponse{}

	if err := json.NewDecoder(res.Body).Decode(&searchRes); err != nil {
		return nil, err
	}

	products := make([]Product, 0, len(searchRes.Hits.Hits))
	for _,hit := range searchRes.Hits.Hits {
		products = append(products, Product{
			ID: hit.Id,
			Name: hit.Source.Name,
			Description: hit.Source.Description,
			Price: hit.Source.Price,
		})
	}
	return products, nil
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error){

	var body bytes.Buffer
	bodyItems := make([]map[string]string, len(ids))
	for i, id := range ids {
		bodyItems[i] = map[string]string{"_id": id}
	}
	if err := json.NewEncoder(&body).Encode(map[string]interface{}{
		"docs": bodyItems,
	}); err != nil {
		return nil, err
	}

	res, err := r.client.Mget(
		&body,
		r.client.Mget.WithContext(ctx),
		r.client.Mget.WithIndex("catalog"),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	multiGetRes := MultiGetResponse{}
	if err := json.NewDecoder(res.Body).Decode(&multiGetRes); err != nil {
		return nil, err
	}

	products := make([]Product, 0, len(ids))
	for _, doc := range multiGetRes.Docs {
		if doc.Found && doc.Source != nil {
			products = append(products, Product{
				ID:          doc.ID,
				Name:        doc.Source.Name,
				Description: doc.Source.Description,
				Price:       doc.Source.Price,
			})
		}
	}

	return products, nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error){
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query" : query,
				"fields": []string{"name", "description"},
			}, 
		},
	}

	queryBody, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, err		
	}
	
	res, err := r.client.Search(
		r.client.Search.WithIndex("catalog"),
		r.client.Search.WithBody(bytes.NewReader(queryBody)),
		r.client.Search.WithFrom(int(skip)),
		r.client.Search.WithSize(int(take)),
		r.client.Search.WithContext(ctx),
	)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	searchRes := SearchResponse{}

	if err := json.NewDecoder(res.Body).Decode(&searchRes); err != nil {
		return nil, err
	}

	products := make([]Product, 0, len(searchRes.Hits.Hits))
	for _,hit := range searchRes.Hits.Hits {
		products = append(products, Product{
			ID: hit.Id,
			Name: hit.Source.Name,
			Description: hit.Source.Description,
			Price: hit.Source.Price,
		})
	}
	return products, nil

}
