package oauth

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTokenStoreSaveLoad(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "tokens", "tokens.json")
	store := NewTokenStore(path)

	expected := &Token{
		AccessToken:  "access",
		RefreshToken: "refresh",
		TokenType:    "Bearer",
		Scope:        "read:profile",
		Expiry:       time.Now().UTC().Truncate(time.Second),
	}

	if err := store.Save(expected); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("token file mode = %o, want 600", got)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.AccessToken != expected.AccessToken {
		t.Fatalf("access token = %q, want %q", loaded.AccessToken, expected.AccessToken)
	}
	if loaded.RefreshToken != expected.RefreshToken {
		t.Fatalf("refresh token = %q, want %q", loaded.RefreshToken, expected.RefreshToken)
	}
	if !loaded.Expiry.Equal(expected.Expiry) {
		t.Fatalf("expiry = %v, want %v", loaded.Expiry, expected.Expiry)
	}
}

func TestTokenStoreDelete(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "tokens.json")
	store := NewTokenStore(path)

	if err := store.Save(&Token{AccessToken: "access", Expiry: time.Now()}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if err := store.Delete(); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("token file still exists, stat err = %v", err)
	}
}
