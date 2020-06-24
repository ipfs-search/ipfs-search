package elasticsearch

import (
	"context"
	"fmt"
)

// Exists returns true if the index exists, false otherwise.
func (i *Index) Exists(ctx context.Context) (bool, error) {
	return i.es.IndexExists(i.cfg.Name).Do(ctx)

}

// Create creates an index with given settings.
func (i *Index) Create(ctx context.Context) error {
	_, err := i.es.CreateIndex(i.cfg.Name).BodyJson(map[string]interface{}{
		"settings": i.cfg.Settings,
		"mappings": map[string]interface{}{
			"_doc": i.cfg.Mapping,
		},
	}).Do(ctx)
	return err
}

// GetSettings returns the mapping for an index.
func (i *Index) GetSettings(ctx context.Context) (interface{}, error) {
	responseMap, err := i.es.IndexGetSettings(i.cfg.Name).Do(ctx)
	if err != nil {
		return false, err
	}

	response, found := responseMap[i.cfg.Name]
	if !found {
		return false, fmt.Errorf("index %s not found in result", i.cfg.Name)
	}

	return response.Settings, nil
}

// SetSettings updates the settings of the index.
func (i *Index) SetSettings(ctx context.Context, settings interface{}) error {
	_, err := i.es.IndexPutSettings(i.cfg.Name).BodyJson(settings).Do(ctx)
	return err
}

// GetMapping returns the mapping for an index.
func (i *Index) GetMapping(ctx context.Context) (interface{}, error) {
	responseMap, err := i.es.GetMapping().Index(i.cfg.Name).Type("_doc").Do(ctx)
	if err != nil {
		return false, err
	}

	response, found := responseMap[i.cfg.Name]
	if !found {
		return false, fmt.Errorf("index %s not found in result", i.cfg.Name)
	}

	mappings, ok := response.(map[string]interface{})["mappings"]
	if !ok {
		return false, fmt.Errorf("\"mappings\" not found in result")
	}

	return mappings, nil
}

// SetMapping updates the settings of the index.
func (i *Index) SetMapping(ctx context.Context, mapping interface{}) error {
	_, err := i.es.PutMapping().Index(i.cfg.Name).Type("_doc").BodyJson(mapping.(map[string]interface{})).Do(ctx)
	return err
}
