package main

import (
	"encoding/base64"
	"fmt"

	"github.com/itstheanurag/in-memory/middleware"
	"github.com/itstheanurag/in-memory/store"
	"github.com/itstheanurag/in-memory/store/kv"
)

func main() {
	store := CreateStoreWithMetrics()
	store.Set("pageviews", "0")

	const hits = 5

	cmds := make([]Command, hits)

	for i := range cmds {
		cmds[i] = Command{Op: "INCR", Key: "pageviews"}
	}

	RestoreOnBoot(store, cmds)
	store.Report()

}

func logConfig() middleware.LogConfig {
	return middleware.LogConfig{
		Get:    true,
		Set:    true,
		Delete: true,
		Len:    true,
		Keys:   true,
		Incr:   false,
		Decr:   false,
	}
}

func CreateStore() store.Storer {
	keyValueStore := kv.NewStore(20)
	loggingStore := middleware.NewLoggingMiddleware(keyValueStore, logConfig())
	return loggingStore
}

func CreateStoreWithMetrics() *middleware.MetricsMiddleware {
	keyValueStore := kv.NewStore(20)
	loggingStore := middleware.NewLoggingMiddleware(
		keyValueStore,
		logConfig(),
	)

	return middleware.NewMetricsMiddleware(loggingStore)
}

func SetWithEncryption(store store.Storer, key, value string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(value))

	if err := store.Set(key, encoded); err != nil {
		return "", err
	}

	return encoded, nil
}

func GetWithDecryption(store store.Storer, key string) (string, error) {
	encoded, err := store.Get(key)

	if err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

func PopulateDefaults(store store.Storer) error {
	defaults := map[string]string{
		"enviornment": "production",
		"version":     "0.0.1",
		"debug":       "true",
		"author":      "gaurav",
		"project":     "in-memory",
	}

	for key, value := range defaults {
		if err := store.Set(key, value); err != nil {
			return fmt.Errorf("PopulateDefaults: %w", err)
		}
	}

	return nil
}

func SetGetCommandsReboot() {
	cmds := []Command{
		{Op: "SET", Key: "env", Value: "production"},
		{Op: "SET", Key: "version", Value: "0.0.1"},
		{Op: "SET", Key: "debug", Value: "true"},
		{Op: "GET", Key: "version"},
		{Op: "SET", Key: "region", Value: "eu-west-1"},
		{Op: "SET", Key: "project", Value: "in-memory"},
		{Op: "GET", Key: "env"},
	}
	store := CreateStoreWithMetrics()
	RestoreOnBoot(store, cmds)
	store.Report()
}
