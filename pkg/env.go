// all environment variables including default values put here

package pkg

import (
	"os"
	"path/filepath"
)

const (
	Logo = "🦞"
	// AppName is the name of the app
	AppName = "PicoClaw"

	DefaultPicoClawHome = ".picoclaw"
	WorkspaceName       = "workspace"
)

// GetPicoClawHome returns the picoclaw home directory.
// Priority:
//  1. PICOCLAW_HOME environment variable (explicit override)
//  2. Same directory as the current executable (if .picoclaw folder exists)
//  3. Default user home directory (~/.picoclaw)
func GetPicoClawHome(envHomeKey string) string {
	if home := os.Getenv(envHomeKey); home != "" {
		return home
	}

	// Check same directory as the current executable
	if exe, err := os.Executable(); err == nil {
		home := filepath.Join(filepath.Dir(exe), DefaultPicoClawHome)
		if info, err := os.Stat(home); err == nil && info.IsDir() {
			return home
		}
	}

	userHome, _ := os.UserHomeDir()
	return filepath.Join(userHome, DefaultPicoClawHome)
}
