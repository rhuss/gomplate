package vault

import (
	"os"
	"path"
	"testing"

	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
	"github.com/stretchr/testify/assert"
)

func TestNewAppRoleAuthStrategy(t *testing.T) {
	defer os.Unsetenv("VAULT_ROLE_ID")
	defer os.Unsetenv("VAULT_SECRET_ID")
	defer os.Unsetenv("VAULT_AUTH_APPROLE_MOUNT")
	defer os.Unsetenv("VAULT_ROLE_ID_FILE")
	defer os.Unsetenv("VAULT_SECRET_ID_FILE")

	os.Unsetenv("VAULT_ROLE_ID")
	os.Unsetenv("VAULT_SECRET_ID")
	assert.Nil(t, NewAppRoleAuthStrategy())

	os.Setenv("VAULT_ROLE_ID", "foo")
	assert.Nil(t, NewAppRoleAuthStrategy())

	os.Unsetenv("VAULT_ROLE_ID")
	os.Setenv("VAULT_SECRET_ID", "bar")
	assert.Nil(t, NewAppRoleAuthStrategy())

	os.Setenv("VAULT_ROLE_ID", "foo")
	os.Setenv("VAULT_SECRET_ID", "bar")
	auth := NewAppRoleAuthStrategy()
	assert.Equal(t, "foo", auth.RoleID)
	assert.Equal(t, "bar", auth.SecretID)
	assert.Equal(t, "approle", auth.Mount)

	os.Setenv("VAULT_ROLE_ID", "baz")
	os.Setenv("VAULT_SECRET_ID", "qux")
	os.Setenv("VAULT_AUTH_APPROLE_MOUNT", "quux")
	auth = NewAppRoleAuthStrategy()
	assert.Equal(t, "baz", auth.RoleID)
	assert.Equal(t, "qux", auth.SecretID)
	assert.Equal(t, "quux", auth.Mount)
}

func TestNewAppRoleAuthStrategy_Files(t *testing.T) {
	defer os.Unsetenv("VAULT_ROLE_ID_FILE")
	defer os.Unsetenv("VAULT_SECRET_ID_FILE")

	assert.Nil(t, NewAppRoleAuthStrategy())

	os.Setenv("VAULT_ROLE_ID_FILE", "foo")
	assert.Nil(t, NewAppRoleAuthStrategy())

	os.Unsetenv("VAULT_ROLE_ID_FILE")
	os.Setenv("VAULT_SECRET_ID_FILE", "bar")
	assert.Nil(t, NewAppRoleAuthStrategy())

	secretsDir := "/run/secrets"
	roleIDPath := path.Join(secretsDir, "roleID")
	secretIDPath := path.Join(secretsDir, "secretID")
	fs := memfs.Create()
	err := vfs.MkdirAll(fs, secretsDir, 0700)
	assert.NoError(t, err)
	f, err := vfs.Create(fs, roleIDPath)
	assert.NoError(t, err)
	f.Write([]byte("foo"))
	f, err = vfs.Create(fs, secretIDPath)
	assert.NoError(t, err)
	f.Write([]byte("bar"))

	os.Setenv("VAULT_ROLE_ID_FILE", roleIDPath)
	os.Setenv("VAULT_SECRET_ID_FILE", secretIDPath)
	auth := NewAppRoleAuthStrategy(fs)
	assert.Equal(t, "foo", auth.RoleID)
	assert.Equal(t, "bar", auth.SecretID)
}
