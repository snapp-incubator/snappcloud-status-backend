package config

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

const (
	delimeter = "."
	seperator = "__"

	envPrefix = "SNAPPCLOUD_STATUS_BACKEND__"

	tagName = "koanf"

	upTemplate     = "================ Loaded Configuration ================"
	bottomTemplate = "======================================================"
)

func Load(print bool) *Config {
	k := koanf.New(delimeter)

	// load default configuration from file
	if err := LoadValues(k); err != nil {
		log.Fatalf("error loading default values: \n%v", err)
	}

	// load default configuration from environment variables
	if err := loadEnv(k); err != nil {
		log.Printf("error loading environment variables: \n%v", err)
	}

	config := Config{}
	var tag = koanf.UnmarshalConf{Tag: tagName}
	if err := k.UnmarshalWithConf("", &config, tag); err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	if print {
		// pretty print loaded configuration using provided template
		log.Printf("%s\n%v\n%s\n", upTemplate, spew.Sdump(config), bottomTemplate)
	}

	return &config
}

//go:embed values.yml
var values []byte

func LoadValues(k *koanf.Koanf) error {
	if err := k.Load(rawbytes.Provider(values), yaml.Parser()); err != nil {
		return fmt.Errorf("error loading values: %s", err)
	}

	return nil
}

func loadEnv(k *koanf.Koanf) error {
	callback := func(source string) string {
		base := strings.ToLower(strings.TrimPrefix(source, envPrefix))
		return strings.ReplaceAll(base, seperator, delimeter)
	}

	// load environment variables
	if err := k.Load(env.Provider(envPrefix, delimeter, callback), nil); err != nil {
		return fmt.Errorf("error loading environment variables: %s", err)
	}

	return nil
}
