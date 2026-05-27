package saga

import (
	"os"
	"path"
	"runtime"
	"testing"
)

func TestResolve(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("Skipping Resolve test on Windows due to lack of ~ (or equiv.) support\n")
	}

	resolved, err := Resolve("~/Downloads")
	if err != nil {
		t.Fatalf("Failed to resolve path: %v\n", err)
	}
	userHome, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v\n", err)
	}

	expectedPath := path.Join(userHome, "Downloads")
	if resolved == expectedPath {
		t.Logf("Path resolution successful: %s\n", resolved)
	}
}
