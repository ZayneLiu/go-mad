package saga

import (
	"os"
	"path/filepath"
	"strings"
)

// Resolve expands the ~ to the user's home directory and returns the absolute path.
//
// on Windows, as it doesn't have `~` equivalent, it currently does nothing special, paths are resolved as-is.
//
// TODO:
//   - Support Windows env var expansion like %USERPROFILE% and %HOMEPATH%
func Resolve(path string) (string, error) {
	path = filepath.Clean(path)
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if path == "~" {
			return home, nil
		}
		return filepath.Join(home, path[2:]), nil
	}
	return filepath.Abs(path)
}
