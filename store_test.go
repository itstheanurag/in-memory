package main

import (
	"reflect"
	"testing"
)

func TestKeys_ReturnsAllKeysSorted(t *testing.T) {
	store := NewStore(5)

	store.Set("ironman", "tony stark")
	store.Set("captainamerica", "steve rogers")
	store.Set("thor", "thor")
	store.Set("spiderman", "peter parker")
	store.Set("hulk", "bruce banner")

	keys_in_store := store.Keys()
	keys_should_be := []string{"captainamerica", "hulk", "ironman", "spiderman", "thor"}

	if len(keys_in_store) != len(keys_should_be) {
		t.Errorf("Keys() = %v, want %v", keys_in_store, keys_should_be)
	}

	if !reflect.DeepEqual(keys_in_store, keys_should_be) {
		t.Errorf("Keys() = %v, want %v", keys_in_store, keys_should_be)
	}

}

func TestKeys_EmptyStore(t *testing.T) {
	store := NewStore(0)
	keys_in_store := store.Keys()

	if len(keys_in_store) != 0 {
		t.Errorf("Keys() = %v, want %v", keys_in_store, []string{})
	}
}

func TestSetGet_RoundTrip(t *testing.T) {
	store := NewStore(1)
	store.Set("ironman", "tony stark")

	if value, err := store.Get("ironman"); err != nil || value != "tony stark" {
		t.Errorf("Get() = %v, want %v", value, "tony stark")
	}

	if value, err := store.Get("missing"); err == nil || value != "" {
		t.Errorf("Get() = %v, want %v", value, "")
	}
}

func TestEmptyKey(t *testing.T) {
	store := NewStore(1)
	if value, err := store.Get(""); err == nil || value != "" {
		t.Errorf("Get() = %v, want %v", value, "")
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name      string
		setup     map[string]string
		deleteKey string
		wantLen   int
	}{
		{
			name:      "deletes existing key",
			setup:     map[string]string{"ironman": "tony stark"},
			deleteKey: "ironman",
			wantLen:   0,
		},
		{
			name:      "does nothing for missing key",
			setup:     map[string]string{"ironman": "tony stark"},
			deleteKey: "missing",
			wantLen:   1,
		},
		{
			name:      "does nothing on empty store",
			setup:     map[string]string{},
			deleteKey: "ironman",
			wantLen:   0,
		},
		{
			name: "deletes one of multiple keys",
			setup: map[string]string{
				"ironman":        "tony stark",
				"captainamerica": "steve rogers",
			},
			deleteKey: "captainamerica",
			wantLen:   1,
		},
		{
			name:      "deletes only key",
			setup:     map[string]string{"ironman": "tony stark"},
			deleteKey: "ironman",
			wantLen:   0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := NewStore(len(tc.setup))

			for key, value := range tc.setup {
				store.Set(key, value)
			}

			store.Delete(tc.deleteKey)

			if got := len(store.Keys()); got != tc.wantLen {
				t.Errorf(
					"after Delete(%q), len(Keys()) = %v, want %v",
					tc.deleteKey,
					got,
					tc.wantLen,
				)
			}

			if _, err := store.Get(tc.deleteKey); err == nil {
				t.Errorf(
					"after Delete(%q), key still present",
					tc.deleteKey,
				)
			}
		})
	}
}
