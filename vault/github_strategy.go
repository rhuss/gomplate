package vault

import (
	"fmt"
	"os"

	"github.com/blang/vfs"
)

// GitHubAuthStrategy - an AuthStrategy that uses Vault's app-id authentication backend.
type GitHubAuthStrategy struct {
	*Strategy
	Token string `json:"token"`
}

// NewGitHubAuthStrategy - create an AuthStrategy that uses Vault's app-id auth
// backend.
func NewGitHubAuthStrategy(fsOverrides ...vfs.Filesystem) *GitHubAuthStrategy {
	var fs vfs.Filesystem
	if len(fsOverrides) == 0 {
		fs = vfs.OS()
	} else {
		fs = fsOverrides[0]
	}

	mount := os.Getenv("VAULT_AUTH_GITHUB_MOUNT")
	if mount == "" {
		mount = "github"
	}
	token := GetValue("VAULT_AUTH_GITHUB_TOKEN", fs)
	if token != "" {
		return &GitHubAuthStrategy{&Strategy{mount, nil}, token}
	}
	return nil
}

func (a *GitHubAuthStrategy) String() string {
	return fmt.Sprintf("token: %s, mount: %s", a.Token, a.Mount)
}
