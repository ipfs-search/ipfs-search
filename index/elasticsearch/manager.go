package elasticsearch

import (
	"context"
	"fmt"
)

// Exists returns true if the index exists, false otherwise.
func (i *Index) Exists(ctx context.Context) (bool, error) {
	return i.Client.IndexExists(i.Name).Do(ctx)

}

// Create creates an index with given settings.
func (i *Index) Create(ctx context.Context, settings interface{}, mapping interface{}) error {
	_, err := i.Client.CreateIndex(i.Name).BodyJson(map[string]interface{}{
		"settings": settings,
		"mappings": map[string]interface{}{
			"_doc": mapping,
		},
	}).Do(ctx)
	return err
}

// GetSettings returns the mapping for an index.
func (i *Index) GetSettings(ctx context.Context) (interface{}, error) {
	responseMap, err := i.Client.IndexGetSettings(i.Name).Do(ctx)
	if err != nil {
		return false, err
	}

	response, found := responseMap[i.Name]
	if !found {
		return false, fmt.Errorf("index %s not found in result", i.Name)
	}

	return response.Settings, nil
}

// SetSettings updates the settings of the index.
func (i *Index) SetSettings(ctx context.Context, settings interface{}) error {
	_, err := i.Client.IndexPutSettings(i.Name).BodyJson(settings).Do(ctx)
	return err
}

// GetMapping returns the mapping for an index.
func (i *Index) GetMapping(ctx context.Context) (interface{}, error) {
	responseMap, err := i.Client.GetMapping().Index(i.Name).Type("_doc").Do(ctx)
	if err != nil {
		return false, err
	}

	response, found := responseMap[i.Name]
	if !found {
		return false, fmt.Errorf("index %s not found in result", i.Name)
	}

	mappings, ok := response.(map[string]interface{})["mappings"]
	if !ok {
		return false, fmt.Errorf("\"mappings\" not found in result")
	}

	return mappings, nil
}

// SetMapping updates the settings of the index.
func (i *Index) SetMapping(ctx context.Context, mapping interface{}) error {
	_, err := i.Client.PutMapping().Index(i.Name).Type("_doc").BodyJson(mapping.(map[string]interface{})).Do(ctx)
	return err
}
