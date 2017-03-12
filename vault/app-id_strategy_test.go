package vault

import (
	"os"
	"path"
	"testing"

	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
	"github.com/stretchr/testify/assert"
)

func TestNewAppIDAuthStrategy(t *testing.T) {
	defer os.Unsetenv("VAULT_APP_ID")
	defer os.Unsetenv("VAULT_USER_ID")
	defer os.Unsetenv("VAULT_AUTH_APP_ID_MOUNT")
	defer os.Unsetenv("VAULT_APP_ID_FILE")
	defer os.Unsetenv("VAULT_USER_ID_FILE")

	assert.Nil(t, NewAppIDAuthStrategy())

	os.Setenv("VAULT_APP_ID", "foo")
	assert.Nil(t, NewAppIDAuthStrategy())

	os.Unsetenv("VAULT_APP_ID")
	os.Setenv("VAULT_USER_ID", "bar")
	assert.Nil(t, NewAppIDAuthStrategy())

	os.Setenv("VAULT_APP_ID", "foo")
	os.Setenv("VAULT_USER_ID", "bar")
	auth := NewAppIDAuthStrategy()
	assert.Equal(t, "foo", auth.AppID)
	assert.Equal(t, "bar", auth.UserID)
	assert.Equal(t, "app-id", auth.Mount)

	os.Setenv("VAULT_APP_ID", "baz")
	os.Setenv("VAULT_USER_ID", "qux")
	os.Setenv("VAULT_AUTH_APP_ID_MOUNT", "quux")
	auth = NewAppIDAuthStrategy()
	assert.Equal(t, "baz", auth.AppID)
	assert.Equal(t, "qux", auth.UserID)
	assert.Equal(t, "quux", auth.Mount)
}

func TestNewAppIDAuthStrategy_Files(t *testing.T) {
	defer os.Unsetenv("VAULT_APP_ID_FILE")
	defer os.Unsetenv("VAULT_USER_ID_FILE")

	os.Setenv("VAULT_APP_ID_FILE", "foo")
	assert.Nil(t, NewAppIDAuthStrategy())

	os.Unsetenv("VAULT_APP_ID_FILE")
	os.Setenv("VAULT_USER_ID_FILE", "bar")
	assert.Nil(t, NewAppIDAuthStrategy())

	secretsDir := "/run/secrets"
	appIDPath := path.Join(secretsDir, "appId")
	userIDPath := path.Join(secretsDir, "userId")
	fs := memfs.Create()
	err := vfs.MkdirAll(fs, secretsDir, 0700)
	assert.NoError(t, err)
	f, err := vfs.Create(fs, appIDPath)
	assert.NoError(t, err)
	f.Write([]byte("foo"))
	f, err = vfs.Create(fs, userIDPath)
	assert.NoError(t, err)
	f.Write([]byte("bar"))

	os.Setenv("VAULT_APP_ID_FILE", appIDPath)
	os.Setenv("VAULT_USER_ID_FILE", userIDPath)
	auth := NewAppIDAuthStrategy(fs)
	assert.Equal(t, "foo", auth.AppID)
	assert.Equal(t, "bar", auth.UserID)
}
