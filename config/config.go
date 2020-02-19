package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v2"
)

// Provider is the main config provider instance for managing config.
var Provider = &configProvider{}

type configProvider struct {
	mu     sync.Mutex
	config *Config
}

func new() *Config {
	return &Config{
		Presets: make(map[string]Preset),
	}
}

func (p *configProvider) readConfigOrDie() {
	c := new()
	path := getConfigPathOrDie()
	b, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		ioutil.WriteFile(path, b, os.ModePerm)
	} else if err != nil {
		log.Fatalf("Error reading config: %s", err.Error())
	}

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		log.Fatalf("Error parsing config: %s", err.Error())
	}
	p.config = c
}

func (p *configProvider) GetConfig() *Config {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.config == nil {
		p.readConfigOrDie()
	}
	return p.config
}

func (p *configProvider) SetConfig(c *Config) {
	p.mu.Lock()
	defer p.mu.Unlock()

	b, _ := yaml.Marshal(c)
	ioutil.WriteFile(getConfigPathOrDie(), b, os.ModePerm)
	p.config = c
}

// Preset represents a defined tf-provider-swap workflow.
type Preset struct {
	Provider  string // provider to update
	BinPath   string // path to binary
	PreUpdate string // shell script to execute before update
}

// Config represents the CLI configuration.
type Config struct {
	Presets map[string]Preset // presets defined by user
}

// getConfigPathOrDie gets path to the tf-provider-swap config or panics.
func getConfigPathOrDie() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %s\n", err.Error())
	}

	return filepath.Join(home, ".tf-provider-swap")
}
