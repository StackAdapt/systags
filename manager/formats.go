package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// FormatJson attempts to convert tags into a JSON
// string. Returns error if the conversion fails.
func FormatJson(tags Tags) (string, error) {

	// Try and convert specified tags data to JSON
	out, err := json.MarshalIndent(tags, "", "  ")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// FormatYaml attempts to convert tags into a YAML
// string. Returns error if the conversion fails.
func FormatYaml(tags Tags) (string, error) {

	// Try to convert tags to YAML
	out, err := yaml.Marshal(tags)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// FormatToml attempts to convert tags into a TOML
// string. Returns error if the conversion fails.
func FormatToml(tags Tags) (string, error) {

	// Try to convert tags to TOML
	out, err := toml.Marshal(tags)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// FormatEnv attempts to convert tags into a string
// compatible with shell environment variables. Returns
// error if the conversion fails.
func FormatEnv(tags Tags) (string, error) {

	// Regex to replace invalid Bash characters
	validKey := regexp.MustCompile("[^A-Z0-9_]")

	filtered := make(Tags)
	// Iterate through all the tags
	for key, value := range tags {

		// Skip empty
		if key == "" {
			continue
		}

		// Attempt to normalize the key
		k := validKey.ReplaceAllString(
			strings.ToUpper(key), "_",
		)

		// Skip keys starting with digit
		if '0' <= k[0] && k[0] <= '9' {
			continue
		}

		// Values are more permissions but the single
		// quote needs to be properly escaped in Bash
		v := strings.Replace(value, "'", "'\\''", -1)

		filtered[k] = v // Keys are deduplicated here
	}

	result := ""
	// Iterate through filtered tags
	for key, value := range filtered {

		result += fmt.Sprintf(
			"%s='%s' ", key, value,
		)
	}

	return strings.Trim(result, " "), nil
}

// FormatTelegraf attempts to convert tags into a
// string which could be stored as a config file
// for Telegraf by InfluxData. Returns error if the
// conversion fails.
func FormatTelegraf(tags Tags) (string, error) {

	telegraf := struct {
		GlobalTags Tags `json:"global_tags"`
	}{
		GlobalTags: tags,
	}

	var out bytes.Buffer
	enc := toml.NewEncoder(&out)
	enc.SetIndentTables(true)

	// Try encoding tags as TOML
	err := enc.Encode(telegraf)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

// FormatConsul attempts to convert tags into a
// string which could be stored as a config file
// for Consul by HashiCorp. Returns error if the
// conversion fails.
func FormatConsul(tags Tags) (string, error) {

	consul := struct {
		NodeMeta Tags `json:"node_meta"`
	}{
		NodeMeta: tags,
	}

	// Try and convert specified consul data to JSON
	out, err := json.MarshalIndent(consul, "", "  ")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// Format is a type that defines a function signature
// for formatting tags into a specific string format.
type Format func(Tags) (string, error)

// Formats is a registry of tag formatting functions.
var Formats = map[string]Format{
	"json":     FormatJson,
	"yaml":     FormatYaml,
	"yml":      FormatYaml,
	"toml":     FormatToml,
	"env":      FormatEnv,
	"telegraf": FormatTelegraf,
	"consul":   FormatConsul,
}
