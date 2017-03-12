package vault

import (
	"fmt"
	"os"

	"github.com/blang/vfs"
)

// AppRoleAuthStrategy - an AuthStrategy that uses Vault's approle authentication backend.
type AppRoleAuthStrategy struct {
	*Strategy
	RoleID   string `json:"role_id"`
	SecretID string `json:"secret_id"`
}

// NewAppRoleAuthStrategy - create an AuthStrategy that uses Vault's approle auth
// backend.
func NewAppRoleAuthStrategy(fsOverrides ...vfs.Filesystem) *AppRoleAuthStrategy {
	var fs vfs.Filesystem
	if len(fsOverrides) == 0 {
		fs = vfs.OS()
	} else {
		fs = fsOverrides[0]
	}

	roleID := GetValue("VAULT_ROLE_ID", fs)
	secretID := GetValue("VAULT_SECRET_ID", fs)
	mount := os.Getenv("VAULT_AUTH_APPROLE_MOUNT")
	if mount == "" {
		mount = "approle"
	}
	if roleID != "" && secretID != "" {
		return &AppRoleAuthStrategy{&Strategy{mount, nil}, roleID, secretID}
	}
	return nil
}

func (a *AppRoleAuthStrategy) String() string {
	return fmt.Sprintf("role_id: %s, secret_id: %s, mount: %s", a.RoleID, a.SecretID, a.Mount)
}
