package main

import (
	"sync"

	"github.com/itstheanurag/in-memory/store"
)

type Command struct {
	Op    string
	Key   string
	Value string
}

// lowercase function for using locally
func dispatch(s store.Storer, cmd Command) {
	switch cmd.Op {
	case "SET":
		s.Set(cmd.Key, cmd.Value)
	case "GET":
		s.Get(cmd.Key)
	case "DELETE":
		s.Delete(cmd.Key)
	}
}

// Uppercase function for external use ( exported function )
func RestoreOnBoot(s store.Storer, cmds []Command) {

	var wg sync.WaitGroup

	for _, cmd := range cmds {
		// Old way
		// wg.Add(1)
		// go func(cmd Command) {
		// 	defer wg.Done()
		// 	dispatch(s, cmd)
		// }(cmd)

		// New way
		wg.Go(func() {
			dispatch(s, cmd)
		})
	}

	wg.Wait()
}
