package config

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

var (
	errEnvNotFound           = errors.New("env file is not found")
	errYamlNotFound          = errors.New("yaml file is not found")
	paths           []string = []string{
		"config.json",
		"internal/config/proxy.json",
	}
	once sync.Once
	cfg  *viper.Viper
)

func Load() *viper.Viper {
	once.Do(func() {
		err := gotenv.Load(".env")
		if err != nil {
			log.Fatalf("%s", errEnvNotFound)
		}

		envMap := getEnvMap()
		v := viper.New()
		v.SetConfigType("json")

		for _, path := range paths {
			jsonBytes, err := os.ReadFile(path)
			if err != nil {
				log.Fatalf("%s. path: %s", errYamlNotFound.Error(), path)
			}

			tmpl, err := template.New(path).Parse(string(jsonBytes))
			if err != nil {
				log.Fatalf("%s", err.Error())
			}
			var rendered bytes.Buffer
			if err := tmpl.Execute(&rendered, envMap); err != nil {
				log.Fatalf("%s", err.Error())
			}

			if err := v.MergeConfig(bytes.NewBuffer(rendered.Bytes())); err != nil {
				log.Fatalf("%s", err.Error())
			}
		}

		cfg = v
	})

	// pp.Print(v.AllSettings())
	return cfg
}

func getEnvMap() map[string]string {
	envMap := make(map[string]string)
	for _, item := range os.Environ() {
		if split := strings.SplitN(item, "=", 2); len(split) == 2 {
			envMap[split[0]] = split[1]
		}
	}
	return envMap
}
