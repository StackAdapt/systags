package manager

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Tags describe a map of key/value pairs
type Tags map[string]string

// Manager maintains instance information
// for the type which includes in-memory
// representations of the tags as well as
// configurable options.
type Manager struct {
	// The directory for the config files
	ConfigDir string

	// The directory for the system files
	SystemDir string

	logger *slog.Logger

	config Tags
	remote Tags
	system Tags
}

// NewManager initializes a new Manager with
// default configurations and directories.
func NewManager() *Manager {

	// TODO:
	// At the moment, only Linux has been tested, if Windows
	// needs to be supported, then the default ConfigDir and
	// SystemDir will need to be changed.

	m := Manager{
		ConfigDir: "/etc/systags.d",
		SystemDir: "/var/lib/systags",
	}

	m.SetLogger(nil)
	m.Reset()

	return &m
}

// GetLogger returns the current logger from the Manager.
func (m *Manager) GetLogger() *slog.Logger {
	return m.logger
}

// SetLogger sets a new logger for the Manager. If the
// provided logger is nil, it will reset it to default.
func (m *Manager) SetLogger(l *slog.Logger) {

	if l == nil {
		// Avoid invalid loggers
		m.logger = slog.Default()
	} else {
		m.logger = l
	}
}

// Reset clears the current in-memory representation
// of this Manager's config, remote, and system tags.
func (m *Manager) Reset() {

	m.config = make(Tags)
	m.remote = make(Tags)
	m.system = make(Tags)
}

// LoadFiles reads files from the Manager's ConfigDir
// and SystemDir and stores their contents in the
// Manager's internal config, remote, and system tags.
func (m *Manager) LoadFiles() error {

	logger := m.GetLogger()

	configData := make(Tags)
	remoteData := make(Tags)
	systemData := make(Tags)

	// Try to get all files in config directory
	configFiles, err := os.ReadDir(m.ConfigDir)

	if err == nil {

		logger.Debug("reading config directory: " + m.ConfigDir)

		// Iterate through all config files
		for _, file := range configFiles {

			// Ignore folders
			if file.IsDir() {
				continue
			}

			// Ignore files which aren't JSON
			ext := filepath.Ext(file.Name())
			if ext != ".json" {
				continue
			}

			// Construct the full path to the current config file
			configFile := filepath.Join(m.ConfigDir, file.Name())

			logger.Debug(configFile)

			// Attempt to read the contents of the file
			configBytes, err := os.ReadFile(configFile)
			if err != nil {
				return err
			}

			configJson := make(Tags)
			// Try and parse the file as a Tag JSON object
			err = json.Unmarshal(configBytes, &configJson)
			if err != nil {
				return err
			}

			// Merge latest config into result
			for key, value := range configJson {
				configData[key] = value
			}
		}
	}

	// Construct the full path to the system directory files
	remoteFile := filepath.Join(m.SystemDir, "remote.json")
	systemFile := filepath.Join(m.SystemDir, "system.json")

	// Check if remote file exists and then read it
	if _, err := os.Stat(remoteFile); err == nil {

		logger.Debug("reading remote file: " + remoteFile)

		// Attempt to read the contents of the file
		remoteBytes, err := os.ReadFile(remoteFile)
		if err != nil {
			return err
		}

		// Try and parse the file as a Tag JSON object
		err = json.Unmarshal(remoteBytes, &remoteData)
		if err != nil {
			return err
		}
	}

	// Check if system file exists and then read it
	if _, err := os.Stat(systemFile); err == nil {

		logger.Debug("reading system file: " + systemFile)

		// Attempt to read the contents of the file
		systemBytes, err := os.ReadFile(systemFile)
		if err != nil {
			return err
		}

		// Try and parse the file as a Tag JSON object
		err = json.Unmarshal(systemBytes, &systemData)
		if err != nil {
			return err
		}
	}

	m.config = configData
	m.remote = remoteData
	m.system = systemData

	return nil
}

// SaveFiles saves the current state of the Manager's
// remote and system tags to corresponding files in
// the SystemDir. Before writing new data, it attempts
// to create a backup of the existing files.
func (m *Manager) SaveFiles() error {

	logger := m.GetLogger()

	// Construct the full path to the system directory files
	remoteFile := filepath.Join(m.SystemDir, "remote.json")
	systemFile := filepath.Join(m.SystemDir, "system.json")

	// Attempt to convert the remote data to JSON
	remoteJson, err := json.MarshalIndent(m.remote, "", "\t")
	if err != nil {
		return err
	}

	// Attempt to convert the system data to JSON
	systemJson, err := json.MarshalIndent(m.system, "", "\t")
	if err != nil {
		return err
	}

	// Check if remote file exists and then read it
	if _, err := os.Stat(remoteFile); err == nil {

		remoteBackup := remoteFile + ".bak"

		logger.Debug("writing remote backup: " + remoteBackup)

		// Attempt to read the contents of the file
		remoteBytes, err := os.ReadFile(remoteFile)
		if err != nil {
			return err
		}

		// Try and backup the contents of the file
		err = os.WriteFile(remoteBackup, remoteBytes, 0666)
		if err != nil {
			return err
		}
	}

	// Check if system file exists and then read it
	if _, err := os.Stat(systemFile); err == nil {

		systemBackup := systemFile + ".bak"

		logger.Debug("writing system backup: " + systemBackup)

		// Attempt to read the contents of the file
		systemBytes, err := os.ReadFile(systemFile)
		if err != nil {
			return err
		}

		// Try and backup the contents of the file
		err = os.WriteFile(systemBackup, systemBytes, 0666)
		if err != nil {
			return err
		}
	}

	logger.Debug("writing remote file: " + remoteFile)

	// Attempt to write the current tag content
	err = os.WriteFile(remoteFile, remoteJson, 0666)
	if err != nil {
		return err
	}

	logger.Debug("writing system file: " + systemFile)

	// Attempt to write the current tag content
	err = os.WriteFile(systemFile, systemJson, 0666)
	if err != nil {
		return err
	}

	return nil
}

// UpdateRemote fetches remote tags from the instance
// it's currently running on. This operation may take
// some time to complete, controlled by timeout. If a
// retry duration is provided, this function will auto
// retry for the duration if empty tags are returned.
// A bounded exponential backoff strategy is employed.
func (m *Manager) UpdateRemote(timeout time.Duration, retry time.Duration) error {

	// TODO:
	// At the moment, only AWS is supported, but if you want
	// to support GCP or other cloud providers in the future
	// then you will need to implement a function to detect
	// which cloud provider is being used and add the feature
	// to update the tags similar to how it's done now in AWS.

	var err error
	var res Tags

	// Sleep duration to start with
	curInterval := 1 * time.Second

	// Maximum duration to sleep for
	maxInterval := 5 * time.Second

	startTime := time.Now()
	untilTime := startTime.Add(retry)

	for {
		res, err = getAwsTags(m.GetLogger(), timeout)
		if err != nil {
			return err
		}

		// Tags are not empty or we have reached time limit
		if len(res) > 0 || time.Since(startTime) > retry {
			break
		}

		// Avoid exceeding the time limit
		interval := time.Until(untilTime)
		if interval > curInterval {
			interval = curInterval
		}

		time.Sleep(interval)

		// Adjust sleep duration to take longer the next time
		curInterval = time.Duration(float64(curInterval) * 2)
		if curInterval > maxInterval {
			curInterval = maxInterval
		}
	}

	m.remote = res
	return nil
}

// ConfigTags returns the Manager's config tags.
func (m *Manager) ConfigTags() Tags {

	return m.config
}

// RemoteTags returns the Manager's remote tags.
func (m *Manager) RemoteTags() Tags {

	return m.remote
}

// SystemTags returns the Manager's system tags.
func (m *Manager) SystemTags() Tags {

	return m.system
}

// GetTags returns the combined set of config, system,
// and remote tags, but it filters the tags based on
// regular expressions provided in the "pick" and "omit"
// parameters. If the "regex" parameter is set to false,
// the function treats the "pick" and "omit" parameters
// as comma-separated lists of exact keys to include or
// exclude, respectively.
func (m *Manager) GetTags(
	regex bool,
	pick string,
	omit string,
) Tags {

	if !regex && pick != "" {

		// Escape all regex characters for safety
		pick = regexp.QuoteMeta(pick)

		// Convert all commas into pipe characters
		pick = strings.Replace(pick, ",", "|", -1)

		// Wrap the entire value in regex structure
		pick = fmt.Sprintf("^(%s)$", pick)
	}

	if !regex && omit != "" {

		// Escape all regex characters for safety
		omit = regexp.QuoteMeta(omit)

		// Convert all commas into pipe characters
		omit = strings.Replace(omit, ",", "|", -1)

		// Wrap the entire value in regex structure
		omit = fmt.Sprintf("^(%s)$", omit)
	}

	pickRegex := regexp.MustCompile(pick)
	omitRegex := regexp.MustCompile(omit)

	combined := make(Tags)

	// Merge remote tags into combined
	for key, value := range m.remote {
		combined[key] = value
	}

	// Merge config tags into combined
	for key, value := range m.config {
		combined[key] = value
	}

	// Merge system tags into combined
	for key, value := range m.system {
		combined[key] = value
	}

	picked := make(Tags)
	omited := make(Tags)

	// Filter out all the picked values
	for key, value := range combined {

		if pick == "" {
			// Select every key
			picked[key] = value

		} else {

			// Use regex to select the keys
			if pickRegex.MatchString(key) {
				picked[key] = value
			}
		}
	}

	// Filter out all the omited values
	for key, value := range picked {

		if omit == "" {
			// Select every key
			omited[key] = value

		} else {

			// Use regex to select the keys
			if !omitRegex.MatchString(key) {
				omited[key] = value
			}
		}
	}

	return omited
}

// GetTag returns a tag by its key from system, config,
// or remote tags, in that order. If the key doesn't exist
// in any of the tag sets, it returns the default value.
func (m *Manager) GetTag(key string, def string) string {

	// Whether system has the key
	value, found := m.system[key]
	if found {
		return value
	}

	// Whether config has the key
	value, found = m.config[key]
	if found {
		return value
	}

	// Whether remote has the key
	value, found = m.remote[key]
	if found {
		return value
	}

	return def
}

// SetTag sets a tag with the specified key and value in
// the system tags. It returns the previous value of the
// tag, or an empty string if the tag did not exist.
func (m *Manager) SetTag(key string, val string) string {

	// Retrieve the current value
	existing, _ := m.system[key]

	// Apply new value
	m.system[key] = val
	return existing
}

// RemoveTag removes a tag with the specified key from
// the system tags. It returns the previous value of the
// tag, or an empty string if the tag did not exist.
func (m *Manager) RemoveTag(key string) string {

	// Retrieve the current value
	existing, _ := m.system[key]

	// Delete the value
	delete(m.system, key)
	return existing
}

// PrefixTags returns new tags based on the specified
// tags but with each key prepended with prefix string
func (m *Manager) PrefixTags(tags Tags, prefix string) Tags {

	if prefix == "" {
		return tags
	}

	result := make(Tags)
	// Loop through specified tags
	for key, value := range tags {
		result[prefix+key] = value
	}

	return result
}

// SuffixTags returns new tags based on the specified
// tags but with each key appended with suffix string
func (m *Manager) SuffixTags(tags Tags, suffix string) Tags {

	if suffix == "" {
		return tags
	}

	result := make(Tags)
	// Loop through specified tags
	for key, value := range tags {
		result[key+suffix] = value
	}

	return result
}
