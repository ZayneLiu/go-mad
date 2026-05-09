package saga

import (
	"os"
	"path/filepath"
	"testing"
)

func TestZipAndUnzip(t *testing.T) {
	tempDir, _ := ResolvePath("~/Downloads")

	// Create source file
	srcFile := filepath.Join(tempDir, "source.txt")
	f := File(srcFile)
	if err := f.WriteText("archive content"); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Zip it
	zipPath := filepath.Join(tempDir, "archive.zip")
	if err := Zip(zipPath, f); err != nil {
		t.Fatalf("Zip failed: %v", err)
	}

	t.Logf("Created zip file at: %s", zipPath)

	zipHandle := File(zipPath)
	exists, _ := zipHandle.Exists()
	if !exists {
		t.Fatalf("Zip file was not created")
	}

	// Unzip it
	destDir := filepath.Join(tempDir, "extracted")
	if err := Unzip(zipHandle, destDir); err != nil {
		t.Fatalf("Unzip failed: %v", err)
	}

	// Verify extraction
	extractedContent, err := os.ReadFile(filepath.Join(destDir, "source.txt"))
	if err != nil || string(extractedContent) != "archive content" {
		t.Fatalf("Unzip content mismatch or file missing")
	}
}
