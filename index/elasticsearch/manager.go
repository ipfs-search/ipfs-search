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
	log.Printf("Ignoring update request for settings on index %s; feature unimplemented", i)
	return nil
	// Note; this needs to discriminate between settings which can and settings which cannot be modified.
	// Ref; https://www.elastic.co/guide/en/elasticsearch/reference/current/index-modules.html#index-modules-settings
	// var err error

	// // Close index
	// _, err = i.es.CloseIndex(i.cfg.Name).Do(ctx)
	// if err != nil {
	// 	return fmt.Errorf("update settings index %s, close, %w", i, err)
	// }

	// // Update settings
	// _, err = i.es.IndexPutSettings(i.cfg.Name).
	// 	BodyJson(i.cfg.Settings).
	// 	Do(ctx)
	// if err != nil {
	// 	return fmt.Errorf("update settings index %s, %w", i, err)
	// }

	// // Reopen index
	// _, err = i.es.OpenIndex(i.cfg.Name).Do(ctx)
	// if err != nil {
	// 	return fmt.Errorf("update settings index %s, reopen, %w", i, err)
	// }

	// return err
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

	got := &Config{
		Name:     i.cfg.Name,
		Settings: settings.(map[string]interface{}),
		Mapping:  mapping.(map[string]interface{}),
	}

	equal := configEqual(i.cfg, got)

	// Below is debug only
	diff := cmp.Diff(i.cfg, got)
	log.Printf("Settings do not match (-want +got):\n%s", diff)

	return equal, nil
}

// ConfigUpdate updates the Elasticsearch settings from the configuration.
func (i *Index) ConfigUpdate(ctx context.Context) error {
	if err := i.setSettings(ctx); err != nil {
		return fmt.Errorf("index %v, updating settings: %w", i, err)
	}
	if err := i.setMapping(ctx); err != nil {
		return fmt.Errorf("index %v, updating mapping: %w", i, err)
	}

	log.Printf("index %v configuration updated", i)
	return nil
}
