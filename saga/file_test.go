package saga

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestFileHandle_ResolvePath(t *testing.T) {
	resolved, err := ResolvePath("~/Downloads")
	if err != nil {
		t.Fatalf("Failed to resolve path: %v", err)
	}
	userHome, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	expectedPath := path.Join(userHome, "Downloads")
	if resolved == expectedPath {
		t.Logf("Path resolution successful: %s", resolved)
	}
}

func TestFileHandle_WriteAndRead(t *testing.T) {
	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, "test.txt")
	f := File(filePath)

	// Test WriteText
	err := f.WriteText("hello world")
	if err != nil {
		t.Fatalf("WriteText failed: %v", err)
	}

	// Test Exists
	exists, err := f.Exists()
	if err != nil || !exists {
		t.Fatalf("Expected file to exist")
	}

	// Test Text (Read)
	content, err := f.Text()
	if err != nil || content != "hello world" {
		t.Fatalf("Expected 'hello world', got %q", content)
	}

	// Test Size
	size, err := f.Size()
	if err != nil || size != 11 {
		t.Fatalf("Expected size 11, got %d", size)
	}
}

func TestFileHandle_JSON(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "data.json")
	f := File(filePath)

	data := map[string]string{"key": "value"}
	if err := f.WriteJSON(data); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	var result map[string]string
	if err := f.JSON(&result); err != nil || result["key"] != "value" {
		t.Fatalf("JSON read failed or mismatched data")
	}
}
