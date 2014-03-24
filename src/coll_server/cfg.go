package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
)

const ConfigDefaultFileName = AppName + ".cfg.json"

// Server configuration
type Config struct {
	Server struct {
		Address string
		Storage struct {
			Backends map[string]string
		}
	}

	configFileName string
	isCreated      bool
}

// Load configuration from given config file name
//
// If the file does not exist, will attempt to create default config file.
// If the config file name is empty, will use default config file:
//   $HOME/.config/AppName.cfg.json
func LoadConfig(configFileName string) (*Config, error) {
	var cfg *Config
	var err error

	// create empty config
	if configFileName == "" {
		cfg, err = getDefaultConfig()
		if err != nil {
			return nil, fmt.Errorf("cannot initialize default config: %v", err)
		}
	} else {
		cfg = &Config{configFileName: configFileName}
	}

	// check for existence
	_, err = os.Stat(cfg.configFileName)
	if os.IsNotExist(err) {
		err = createDefaultConfig(cfg)
		if err != nil {
			return nil, fmt.Errorf("cannot initialize empty config: %v", err)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("error on stat %s: %v", cfg.configFileName, err)
	}

	err = ReadConfig(cfg.configFileName, cfg)
	if err != nil {
		return nil, fmt.Errorf("error on reading config: %v", err)
	}

	return cfg, nil
}

// General function to read a JSON file into any structure
func ReadConfig(jsonFileName string, into interface{}) error {
	rv := reflect.ValueOf(into)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("cannot read into given value.")
	}

	f, err := os.Open(jsonFileName)
	if err != nil {
		return fmt.Errorf("could not open JSON file %s: %v", jsonFileName, err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(into)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}
	return nil
}

// Populates given cfg with default values
func populateDefaultConfig(cfg *Config) {
	cfg.Server.Address = ":8080"
	cfg.Server.Storage.Backends = make(map[string]string)
	cfg.Server.Storage.Backends["mongodb"] = "mongodb://localhost"
}

// Returns empty config object with default config file name set
func getDefaultConfig() (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("cannot lookup current user: %v", err)
	}
	cfgFileName := filepath.Join(usr.HomeDir, ".config", AppName, ConfigDefaultFileName)
	cfg := &Config{configFileName: cfgFileName}
	return cfg, nil
}

// Create a default config at cfg.configFileName
//
// cfg will be populated with default values
func createDefaultConfig(cfg *Config) error {
	cfgPath := filepath.Dir(cfg.configFileName)
	err := os.MkdirAll(cfgPath, 0755)
	if err != nil {
		return fmt.Errorf("cannot create config dir %s: %v", cfg.configFileName, err)
	}
	cfgFile, err := os.OpenFile(cfg.configFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open config file %s for writing: %v", cfg.configFileName, err)
	}
	defer cfgFile.Close()
	populateDefaultConfig(cfg)
	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("cannot encode JSON: %v", err)
	}
	// JSON beautification
	buf := bytes.NewBuffer(nil)
	err = json.Indent(buf, jsonStr, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot indent JSON: %v", err)
	}
	_, err = io.Copy(cfgFile, buf)
	if err != nil {
		return fmt.Errorf("error writing buf to file: %v", err)
	}
	cfg.isCreated = true
	return nil
}
