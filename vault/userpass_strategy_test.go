package vault

import (
	"os"
	"path"
	"testing"

	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
	"github.com/stretchr/testify/assert"
)

func TestNewUserPassAuthStrategy(t *testing.T) {
	defer os.Unsetenv("VAULT_AUTH_USERNAME")
	defer os.Unsetenv("VAULT_AUTH_PASSWORD")
	defer os.Unsetenv("VAULT_AUTH_USERPASS_MOUNT")

	assert.Nil(t, NewUserPassAuthStrategy())

	os.Setenv("VAULT_AUTH_USERNAME", "foo")
	assert.Nil(t, NewUserPassAuthStrategy())

	os.Unsetenv("VAULT_AUTH_USERNAME")
	os.Setenv("VAULT_AUTH_PASSWORD", "bar")
	assert.Nil(t, NewUserPassAuthStrategy())

	os.Setenv("VAULT_AUTH_USERNAME", "foo")
	os.Setenv("VAULT_AUTH_PASSWORD", "bar")
	auth := NewUserPassAuthStrategy()
	assert.Equal(t, "foo", auth.Username)
	assert.Equal(t, "bar", auth.Password)
	assert.Equal(t, "userpass", auth.Mount)

	os.Setenv("VAULT_AUTH_USERNAME", "foo")
	os.Setenv("VAULT_AUTH_PASSWORD", "bar")
	os.Setenv("VAULT_AUTH_USERPASS_MOUNT", "baz")
	auth = NewUserPassAuthStrategy()
	assert.Equal(t, "foo", auth.Username)
	assert.Equal(t, "bar", auth.Password)
	assert.Equal(t, "baz", auth.Mount)
}

func TestNewUserPassAuthStrategy_Files(t *testing.T) {
	defer os.Unsetenv("VAULT_AUTH_USERNAME_FILE")
	defer os.Unsetenv("VAULT_AUTH_PASSWORD_FILE")

	os.Setenv("VAULT_AUTH_USERNAME_FILE", "foo")
	assert.Nil(t, NewUserPassAuthStrategy())

	os.Unsetenv("VAULT_AUTH_USERNAME_FILE")
	os.Setenv("VAULT_AUTH_PASSWORD_FILE", "bar")
	assert.Nil(t, NewUserPassAuthStrategy())

	secretsDir := "/run/secrets"
	usernamePath := path.Join(secretsDir, "username")
	passwordPath := path.Join(secretsDir, "password")
	fs := memfs.Create()
	err := vfs.MkdirAll(fs, secretsDir, 0700)
	assert.NoError(t, err)
	f, err := vfs.Create(fs, usernamePath)
	assert.NoError(t, err)
	f.Write([]byte("foo"))
	f, err = vfs.Create(fs, passwordPath)
	assert.NoError(t, err)
	f.Write([]byte("bar"))

	os.Setenv("VAULT_AUTH_USERNAME_FILE", usernamePath)
	os.Setenv("VAULT_AUTH_PASSWORD_FILE", passwordPath)
	auth := NewUserPassAuthStrategy(fs)
	assert.Equal(t, "foo", auth.Username)
	assert.Equal(t, "bar", auth.Password)
}
