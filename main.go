package main

import (
	"encoding/base64"
	"fmt"

	"github.com/itstheanurag/in-memory/store"
	"github.com/itstheanurag/in-memory/store/kv"
)

func main() {

	keyValueStore := kv.NewStore(20)
	loggingStore := NewLoggingMiddleware(keyValueStore)

	encrypted, err := SetWithEncryption(loggingStore, "Hello", "World!")

	fmt.Print("encrypted: %v", encrypted)
	value, err := loggingStore.Get("Hello")

	if err == nil {
		fmt.Println("Value:", value)
	} else {
		fmt.Println("Key not found")
	}

	keyValueStore.Delete("hero")
	_, err = keyValueStore.Get("hero")

	fmt.Println(err)

	fmt.Println("In memory datastore, fully written in Golang")
}

func SetWithEncryption(store store.Storer, key, value string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(value))

	if err := store.Set(key, encoded); err != nil {
		return "", err
	}

	return encoded, nil
}
