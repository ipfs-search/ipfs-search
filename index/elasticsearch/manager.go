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

	// These settings can only be chanced on closed indexes, hence we should skip them.
	// Missing: "shard.check_on_startup", requires "flat_settings".
	staticSettings := []string{
		"number_of_shards", "codec", "routing_partition_size",
		"load_fixed_bitset_filters_eagerly", "hidden",
	}

	for _, s := range staticSettings {
		delete(indexSettings, s)
	}

	// Update settings
	_, err = i.es.IndexPutSettings(i.cfg.Name).
		BodyJson(i.cfg.Settings).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("update settings index %s, %w", i, err)
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

func (i *Index) settingsUpToDate(ctx context.Context) (bool, error) {
	settings, err := i.getSettings(ctx)
	if err != nil {
		return false, fmt.Errorf("index %v, getting settings: %w", i, err)
	}

	// Below is debug only
	diff := cmp.Diff(i.cfg.Settings, settings)
	log.Printf("Settings do not match (-want +got):\n%s", diff)

	return configEqual(i.cfg.Settings, settings), nil
}

func (i *Index) mappingUpToDate(ctx context.Context) (bool, error) {
	mapping, err := i.getMapping(ctx)
	if err != nil {
		return false, fmt.Errorf("index %v, getting mapping: %w", i, err)
	}

	// Below is debug only
	diff := cmp.Diff(i.cfg.Mapping, mapping)
	log.Printf("Settings do not match (-want +got):\n%s", diff)

	return configEqual(i.cfg.Mapping, mapping), nil

}

// ConfigUpToDate checks whether the settings in Elasticsearch matches the settings in the configuration.
func (i *Index) ConfigUpToDate(ctx context.Context) (bool, error) {
	settingsEqual, err := i.settingsUpToDate(ctx)
	if err != nil {
		return false, err
	}

	mappingEqual, err := i.mappingUpToDate(ctx)
	if err != nil {
		return false, err
	}

	if settingsEqual && mappingEqual {
		return true, nil
	}

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
