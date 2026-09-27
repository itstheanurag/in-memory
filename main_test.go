package main

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"

	"github.com/itstheanurag/in-memory/store/kv"
)

// mockStore helps test error pathways and custom store behavior
type mockStore struct {
	data      map[string]string
	setErr    error
	getErr    error
	deleteErr error
}

func newMockStore() *mockStore {
	return &mockStore{
		data: make(map[string]string),
	}
}

func (m *mockStore) Get(key string) (string, error) {
	if m.getErr != nil {
		return "", m.getErr
	}
	v, ok := m.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return v, nil
}

func (m *mockStore) Set(key, value string) error {
	if m.setErr != nil {
		return m.setErr
	}
	m.data[key] = value
	return nil
}

func (m *mockStore) Delete(key string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.data, key)
	return nil
}

func (m *mockStore) Len() int {
	return len(m.data)
}

func (m *mockStore) Keys() []string {
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

// ==========================================
// Tests for SetWithEncryption
// ==========================================

func TestSetWithEncryption_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		wantEnc   string
		expectErr bool
	}{
		{
			name:      "simple ascii string",
			key:       "greeting",
			value:     "Hello, World!",
			wantEnc:   base64.StdEncoding.EncodeToString([]byte("Hello, World!")),
			expectErr: false,
		},
		{
			name:      "empty value",
			key:       "empty_val",
			value:     "",
			wantEnc:   "",
			expectErr: false,
		},
		{
			name:      "unicode and emojis",
			key:       "emoji_key",
			value:     "🚀 Golang In-Memory Store ⚡️ 世界",
			wantEnc:   base64.StdEncoding.EncodeToString([]byte("🚀 Golang In-Memory Store ⚡️ 世界")),
			expectErr: false,
		},
		{
			name:      "special symbols and newlines",
			key:       "special_chars",
			value:     "line1\nline2\tline3\r\n!@#$%^&*()_+-=[]{}|;':,.<>/?",
			wantEnc:   base64.StdEncoding.EncodeToString([]byte("line1\nline2\tline3\r\n!@#$%^&*()_+-=[]{}|;':,.<>/?")),
			expectErr: false,
		},
		{
			name:      "long string",
			key:       "long_payload",
			value:     strings.Repeat("abcdef123456", 100),
			wantEnc:   base64.StdEncoding.EncodeToString([]byte(strings.Repeat("abcdef123456", 100))),
			expectErr: false,
		},
		{
			name:      "empty key fails with kv.ErrEmptyKey",
			key:       "",
			value:     "some value",
			wantEnc:   "",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := kv.NewStore(10)
			encoded, err := SetWithEncryption(s, tc.key, tc.value)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if encoded != tc.wantEnc {
				t.Errorf("got encoded %q, want %q", encoded, tc.wantEnc)
			}

			// Verify the underlying store contains the base64 value
			rawVal, err := s.Get(tc.key)
			if err != nil {
				t.Fatalf("failed to retrieve key directly from store: %v", err)
			}
			if rawVal != tc.wantEnc {
				t.Errorf("store value = %q, want %q", rawVal, tc.wantEnc)
			}
		})
	}
}

func TestSetWithEncryption_StoreErrors(t *testing.T) {
	t.Run("fails when store returns error on Set", func(t *testing.T) {
		mock := newMockStore()
		mock.setErr = errors.New("disk full or network failure")

		_, err := SetWithEncryption(mock, "testKey", "testVal")
		if err == nil {
			t.Fatal("expected error from store.Set, got nil")
		}
		if !strings.Contains(err.Error(), "disk full") {
			t.Errorf("expected disk full error, got: %v", err)
		}
	})

	t.Run("fails when store reaches capacity", func(t *testing.T) {
		s := kv.NewStore(1) // capacity = 1
		_, err := SetWithEncryption(s, "key1", "val1")
		if err != nil {
			t.Fatalf("first set failed: %v", err)
		}

		// Second set should fail with ErrStoreFull
		_, err = SetWithEncryption(s, "key2", "val2")
		if err == nil {
			t.Fatal("expected ErrStoreFull, got nil")
		}
	})

	t.Run("updating existing key succeeds when store is full", func(t *testing.T) {
		s := kv.NewStore(1)
		_, _ = SetWithEncryption(s, "key1", "val1")

		encoded, err := SetWithEncryption(s, "key1", "new_val")
		if err != nil {
			t.Fatalf("updating existing key on full store failed: %v", err)
		}
		expectedEnc := base64.StdEncoding.EncodeToString([]byte("new_val"))
		if encoded != expectedEnc {
			t.Errorf("got %q, want %q", encoded, expectedEnc)
		}
	})
}

// ==========================================
// Tests for GetWithDecryption
// ==========================================

func TestGetWithDecryption_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		storedVal   string
		wantDecoded string
		expectErr   bool
	}{
		{
			name:        "valid base64 string",
			key:         "key1",
			storedVal:   base64.StdEncoding.EncodeToString([]byte("SecretMessage123")),
			wantDecoded: "SecretMessage123",
			expectErr:   false,
		},
		{
			name:        "empty encoded string",
			key:         "key2",
			storedVal:   "",
			wantDecoded: "",
			expectErr:   false,
		},
		{
			name:        "corrupted non-base64 value",
			key:         "bad_key",
			storedVal:   "!!!ThisIsNotValidBase64@@@",
			wantDecoded: "",
			expectErr:   true,
		},
		{
			name:        "json formatted decoded payload",
			key:         "json_key",
			storedVal:   base64.StdEncoding.EncodeToString([]byte(`{"user":"tony","role":"admin"}`)),
			wantDecoded: `{"user":"tony","role":"admin"}`,
			expectErr:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := kv.NewStore(10)
			_ = s.Set(tc.key, tc.storedVal)

			got, err := GetWithDecryption(s, tc.key)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error for invalid data, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tc.wantDecoded {
				t.Errorf("got decrypted %q, want %q", got, tc.wantDecoded)
			}
		})
	}
}

func TestGetWithDecryption_Errors(t *testing.T) {
	t.Run("key not found error", func(t *testing.T) {
		s := kv.NewStore(10)
		_, err := GetWithDecryption(s, "non_existent_key")
		if err == nil {
			t.Fatal("expected key not found error, got nil")
		}
	})

	t.Run("empty key error", func(t *testing.T) {
		s := kv.NewStore(10)
		_, err := GetWithDecryption(s, "")
		if err == nil {
			t.Fatal("expected empty key error, got nil")
		}
	})

	t.Run("underlying store Get error", func(t *testing.T) {
		mock := newMockStore()
		mock.getErr = errors.New("mock read error")

		_, err := GetWithDecryption(mock, "any_key")
		if err == nil {
			t.Fatal("expected mock read error, got nil")
		}
		if !strings.Contains(err.Error(), "mock read error") {
			t.Errorf("expected 'mock read error', got: %v", err)
		}
	})
}

// ==========================================
// Tests for Encryption & Decryption Roundtrip
// ==========================================

func TestEncryptionDecryption_RoundTrip(t *testing.T) {
	s := kv.NewStore(10)

	testPayloads := map[string]string{
		"alpha":   "value-1",
		"beta":    "hello world across systems",
		"gamma":   "{\"status\": 200, \"active\": true}",
		"unicode": "🎉 Go 1.22+ In-Memory Unit Testing 🚀",
		"binary":  "\x00\x01\x02\x03\x04\xff",
	}

	// 1. Write all encrypted
	for k, v := range testPayloads {
		enc, err := SetWithEncryption(s, k, v)
		if err != nil {
			t.Fatalf("SetWithEncryption failed for key %s: %v", k, err)
		}
		if enc == v && v != "" {
			t.Errorf("encrypted string should not match raw plaintext %q", v)
		}
	}

	// 2. Read and decrypt all
	for k, wantVal := range testPayloads {
		gotVal, err := GetWithDecryption(s, k)
		if err != nil {
			t.Fatalf("GetWithDecryption failed for key %s: %v", k, err)
		}
		if gotVal != wantVal {
			t.Errorf("key %s: got %q, want %q", k, gotVal, wantVal)
		}
	}
}

// ==========================================
// Tests for Factory Functions
// ==========================================

func TestCreateStore(t *testing.T) {
	s := CreateStore()
	if s == nil {
		t.Fatal("CreateStore() returned nil")
	}

	// Test store lifecycle
	err := s.Set("k1", "v1")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := s.Get("k1")
	if err != nil || val != "v1" {
		t.Fatalf("Get failed: got (%q, %v), want (%q, nil)", val, err, "v1")
	}

	if s.Len() != 1 {
		t.Errorf("Len = %d, want 1", s.Len())
	}

	keys := s.Keys()
	if len(keys) != 1 || keys[0] != "k1" {
		t.Errorf("Keys = %v, want ['k1']", keys)
	}

	if err := s.Delete("k1"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if s.Len() != 0 {
		t.Errorf("Len after delete = %d, want 0", s.Len())
	}
}

func TestCreateStoreWithMetrics(t *testing.T) {
	ms := CreateStoreWithMetrics()
	if ms == nil {
		t.Fatal("CreateStoreWithMetrics() returned nil")
	}

	// Exercise all methods on MetricsMiddleware
	_ = ms.Set("mKey1", "mVal1")
	_ = ms.Set("mKey2", "mVal2")

	_, _ = ms.Get("mKey1")
	_, _ = ms.Get("mKeyNonExistent") // Miss

	_ = ms.Len()
	_ = ms.Keys()
	_ = ms.Delete("mKey1")

	// Ensure Report() executes cleanly without panic
	ms.Report()
}

// ==========================================
// Test main() execution
// ==========================================

func TestMainFunction(t *testing.T) {
	// Calling main() should run to completion without panic
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("main() panicked: %v", r)
		}
	}()

	main()
}
