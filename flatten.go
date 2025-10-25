package flattener

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"
)

type Flattener struct {
	Delimeter string
}

// JSON flattens a JSON object seperated by f.Delimeter
func (f Flattener) JSON(input json.RawMessage) ([]byte, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	var data map[string]any
	err := json.Unmarshal(input, &data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal input: %w", err)
	}

	result := make(map[string]any)
	f.Flatten("", data, result)

	sorted := sortMap(result)

	b, err := json.Marshal(sorted)
	if err != nil {
		return nil, fmt.Errorf("could not marshal result: %w", err)
	}

	return b, nil
}

// YAML flattens a YAML document separated by f.Delimeter
func (f Flattener) YAML(input []byte) ([]byte, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	var data map[string]any
	err := yaml.Unmarshal(input, &data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal input: %w", err)
	}

	result := make(map[string]any)
	f.Flatten("", data, result)

	sorted := sortMap(result)

	b, err := yaml.Marshal(sorted)
	if err != nil {
		return nil, fmt.Errorf("could not marshal result: %w", err)
	}

	return b, nil
}

// Flatten recursively flattens nested data into result map
func (f Flattener) Flatten(prefix string, data any, result map[string]any) {
	switch t := data.(type) {
	case map[string]any:
		for k, v := range t {
			var newPrefix string
			if prefix != "" {
				newPrefix = prefix + f.Delimeter
			}

			f.Flatten(newPrefix+k, v, result)
		}
	case []any:
		for i, v := range t {
			var newPrefix string
			if prefix != "" {
				newPrefix = prefix + f.Delimeter
			}

			f.Flatten(newPrefix+fmt.Sprintf("%d", i), v, result)
		}
	default:
		result[prefix] = t
	}
}

// sortMap sorts a map by keys and returns a new map with sorted order
func sortMap(m map[string]any) map[string]any {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sorted := make(map[string]any, len(m))
	for _, k := range keys {
		sorted[k] = m[k]
	}

	return sorted
}
