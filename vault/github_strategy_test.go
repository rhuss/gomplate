package vault

import (
	"os"
	"path"
	"testing"

	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
	"github.com/stretchr/testify/assert"
)

func TestNewGitHubAuthStrategy(t *testing.T) {
	defer os.Unsetenv("VAULT_AUTH_GITHUB_TOKEN")
	defer os.Unsetenv("VAULT_AUTH_GITHUB_MOUNT")

	os.Unsetenv("VAULT_AUTH_GITHUB_TOKEN")
	assert.Nil(t, NewGitHubAuthStrategy())

	os.Setenv("VAULT_AUTH_GITHUB_TOKEN", "foo")
	auth := NewGitHubAuthStrategy()
	assert.Equal(t, "foo", auth.Token)
	assert.Equal(t, "github", auth.Mount)

	os.Setenv("VAULT_AUTH_GITHUB_MOUNT", "bar")
	auth = NewGitHubAuthStrategy()
	assert.Equal(t, "foo", auth.Token)
	assert.Equal(t, "bar", auth.Mount)
}

func TestNewGitHubAuthStrategy_File(t *testing.T) {
	defer os.Unsetenv("VAULT_AUTH_GITHUB_TOKEN_FILE")

	os.Unsetenv("VAULT_AUTH_GITHUB_TOKEN_FILE")
	assert.Nil(t, NewGitHubAuthStrategy())

	secretsDir := "/run/secrets"
	tokenPath := path.Join(secretsDir, "token")
	fs := memfs.Create()
	err := vfs.MkdirAll(fs, secretsDir, 0700)
	assert.NoError(t, err)
	f, err := vfs.Create(fs, tokenPath)
	assert.NoError(t, err)
	f.Write([]byte("foo"))

	os.Setenv("VAULT_AUTH_GITHUB_TOKEN_FILE", tokenPath)
	auth := NewGitHubAuthStrategy(fs)
	assert.Equal(t, "foo", auth.Token)
}
