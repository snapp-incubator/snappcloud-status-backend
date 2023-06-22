package config

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	delimeter = "."
	seperator = "__"

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

	// load default configuration from file
	if err := loadConfigmap(k); err != nil {
		log.Fatalf("error loading from configmap: \n%v", err)
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

func loadConfigmap(k *koanf.Koanf) error {
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		if err == rest.ErrNotInCluster {
			return nil
		}
		panic(fmt.Errorf("error creating Kubernetes config: \n%v", err))
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		panic(fmt.Errorf("error creating Kubernetes client: %v", err))
	}

	// Retrieve the ConfigMap data
	namespace, cmName := os.Getenv("POD_NAMESPACE"), "my-config"
	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), cmName, metav1.GetOptions{})
	if err != nil {
		panic(fmt.Errorf("error retrieving ConfigMap data: %v", err))
	}

	if err := k.Load(rawbytes.Provider([]byte(cm.Data["values.yml"])), nil); err != nil {
		return fmt.Errorf("error loading values: %s", err)
	}

	return nil
}
