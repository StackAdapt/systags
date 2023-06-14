package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// TODO: Documentation
func FormatJson(tags Tags) (string, error) {

	// Try and convert specified tags data to JSON
	out, err := json.MarshalIndent(tags, "", "  ")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// TODO: Documentation
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

// TODO: Documentation
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

// TODO: Documentation
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

// TODO: Documentation
type Format func(Tags) (string, error)

// TODO: Documentation
var Formats = map[string]Format{
	"json":     FormatJson,
	"consul":   FormatConsul,
	"env":      FormatEnv,
	"telegraf": FormatTelegraf,
}
