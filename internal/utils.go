package internal

import (
	"encoding/json"

	"github.com/ghodss/yaml"
)

// ToYaml converts and interface to yaml
func ToYaml(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

// FromYaml converts a YAML document into a map[string]interface{}.
func FromYaml(str string) map[string]interface{} {
	m := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

// ToJson takes an interface, marshals it to json, and returns a string.
func ToJson(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

// FromJson converts a JSON document into a map[string]interface{}.
func FromJson(str string) map[string]interface{} {
	m := map[string]interface{}{}

	if err := json.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}
