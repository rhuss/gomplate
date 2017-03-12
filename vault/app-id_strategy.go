package vault

import (
	"fmt"
	"os"

	"github.com/blang/vfs"
)

// AppIDAuthStrategy - an AuthStrategy that uses Vault's app-id authentication backend.
type AppIDAuthStrategy struct {
	*Strategy
	AppID  string `json:"app_id"`
	UserID string `json:"user_id"`
}

// NewAppIDAuthStrategy - create an AuthStrategy that uses Vault's app-id auth
// backend.
func NewAppIDAuthStrategy(fsOverrides ...vfs.Filesystem) *AppIDAuthStrategy {
	var fs vfs.Filesystem
	if len(fsOverrides) == 0 {
		fs = vfs.OS()
	} else {
		fs = fsOverrides[0]
	}

	appID := GetValue("VAULT_APP_ID", fs)
	userID := GetValue("VAULT_USER_ID", fs)
	mount := os.Getenv("VAULT_AUTH_APP_ID_MOUNT")
	if mount == "" {
		mount = "app-id"
	}
	if appID != "" && userID != "" {
		return &AppIDAuthStrategy{&Strategy{mount, nil}, appID, userID}
	}
	return nil
}

func (a *AppIDAuthStrategy) String() string {
	return fmt.Sprintf("app-id: %s, user-id: %s, mount: %s", a.AppID, a.UserID, a.Mount)
}
