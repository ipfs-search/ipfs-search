package elasticsearch

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"log"
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

// getSettings returns the mapping for an index.
func (i *Index) getSettings(ctx context.Context) (interface{}, error) {
	responseMap, err := i.es.IndexGetSettings(i.cfg.Name).Do(ctx)
	if err != nil {
		return false, err
	}

	response, found := responseMap[i.cfg.Name]
	if !found {
		return false, fmt.Errorf("index %s not found in result", i)
	}

	return response.Settings, nil
}

// setSettings updates the settings of the index.
func (i *Index) setSettings(ctx context.Context) error {
	var err error

	// Remove number_of_shards, which cannot be updated
	// Ref: https://www.elastic.co/guide/en/elasticsearch/reference/current/index-modules.html#index-modules-settings
	newSettings := i.cfg.Settings
	indexSettings := newSettings["index"].(map[string]interface{})
	delete(indexSettings, "number_of_shards")

	// Close index, necessary for some settings
	_, err = i.es.CloseIndex(i.cfg.Name).Do(ctx)
	if err != nil {
		return fmt.Errorf("update settings index %s, close, %w", i, err)
	}

	// Update settings
	_, err = i.es.IndexPutSettings(i.cfg.Name).
		BodyJson(i.cfg.Settings).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("update settings index %s, %w", i, err)
	}

	// Reopen index, necessary for some settings
	_, err = i.es.OpenIndex(i.cfg.Name).Do(ctx)
	if err != nil {
		return fmt.Errorf("update settings index %s, reopen, %w", i, err)
	}

	return nil
}

// getMapping returns the mapping for an index.
func (i *Index) getMapping(ctx context.Context) (interface{}, error) {
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

// setMapping updates the settings of the index.
func (i *Index) setMapping(ctx context.Context) error {
	_, err := i.es.PutMapping().
		Index(i.cfg.Name).
		Type("_doc").
		BodyJson(i.cfg.Mapping).
		Do(ctx)
	return err
}

// ConfigUpToDate checks whether the settings in Elasticsearch matches the settings in the configuration.
func (i *Index) ConfigUpToDate(ctx context.Context) (bool, error) {
	settings, err := i.getSettings(ctx)
	if err != nil {
		return false, fmt.Errorf("index %v, getting settings: %w", i, err)
	}

	mapping, err := i.getMapping(ctx)
	if err != nil {
		return false, fmt.Errorf("index %v, getting mapping: %w", i, err)
	}

	got := Config{
		Name:     i.cfg.Name,
		Settings: settings.(map[string]interface{}),
		Mapping:  mapping.(map[string]interface{}),
	}

	settingsEqual := configEqual(i.cfg.Settings, got.Settings)
	mappingEqual := configEqual(i.cfg.Mapping, got.Mapping)

	if settingsEqual && mappingEqual {
		return true, nil
	}

	// Below is debug only
	diff := cmp.Diff(i.cfg, &got)
	log.Printf("Settings do not match (-want +got):\n%s", diff)

	return false, nil
}

// ConfigUpdate updates the Elasticsearch settings from the configuration.
func (i *Index) ConfigUpdate(ctx context.Context) error {
	if err := i.setSettings(ctx); err != nil {
		return fmt.Errorf("index %v, updating settings: %w", i, err)
	}
	if err := i.setMapping(ctx); err != nil {
		return fmt.Errorf("index %v, updating mapping: %w", i, err)
	}

	log.Printf("index %v configuration update requested", i)
	return nil
}
