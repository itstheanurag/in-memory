package main

import (
	"fmt"
)

func main() {

	keyValueStore := NewStore(2)

	keyValueStore.Set("hero", "sucker")
	value, err := keyValueStore.Get("hero")

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
