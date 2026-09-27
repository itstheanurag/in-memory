package main

import (
	"encoding/base64"
	"fmt"

	"github.com/itstheanurag/in-memory/middleware"
	"github.com/itstheanurag/in-memory/store"
	"github.com/itstheanurag/in-memory/store/kv"
)

func main() {

	loggingStore := CreateStore()
	encrypted, err := SetWithEncryption(loggingStore, "Hello", "World!")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("encrypted:", encrypted)

	val, err := GetWithDecryption(loggingStore, "Hello")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("decrypted:", val)

	// keyValueStore.Delete("hero")
	// _, err = keyValueStore.Get("hero")

	// fmt.Println(err)

	// fmt.Println("In memory datastore, fully written in Golang")
}

func CreateStore() store.Storer {
	keyValueStore := kv.NewStore(20)
	loggingStore := middleware.NewLoggingMiddleware(keyValueStore)
	return loggingStore
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
