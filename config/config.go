package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func Setup() {
	fileName := os.Getenv("HOST")
	if fileName == "" {
		fileName = "local"
	}

	viper.AddConfigPath("./src/config/")
	viper.SetConfigType("yaml")
	viper.SetConfigName(fileName)

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	// Here, we want viper to update the key values it holds
	// for every value that start with $, such as $HELLO
	// update the value by fetching the value from the environment of key HELLO
	for key := range viper.AllSettings() {
		replaceEnvKeys(key)
	}
}

func replaceEnvKeys(key string) {
	value := viper.Get(key)
	switch value := value.(type) {
	case string:
		if strings.HasPrefix(value, "$") {
			viper.Set(key, os.Getenv(strings.TrimPrefix(value, "$")))
		}
	case interface{}:
		for k := range value.(map[string]interface{}) {
			newKey := fmt.Sprintf("%s.%s", key, k)
			replaceEnvKeys(newKey)
		}
	}
}
