package saga

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestZipAndUnzip(t *testing.T) {
	tempDir := t.TempDir()

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

// === Archive: Zip Tests ===

func TestZip_SingleFile(t *testing.T) {
	// TODO: Archive single file, verify contents
}

func TestZip_MultipleFiles(t *testing.T) {
	// TODO: Archive multiple files, verify all present
}

func TestZip_Directory(t *testing.T) {
	// TODO: Archive directory with nested structure
}

func TestZip_PreservesStructure(t *testing.T) {
	// TODO: Relative paths maintained in archive
}

// === Archive: Unzip Tests ===

func TestUnzip_FlatArchive(t *testing.T) {
	// TODO: Extract flat zip to destination
}

func TestUnzip_NestedStructure(t *testing.T) {
	// TODO: Extract archive with directories
}

func TestUnzip_PathTraversalBlocked(t *testing.T) {
	// TODO: Reject entries with ../ that escape dest
}

func TestUnzip_PreservesPermissions(t *testing.T) {
	// TODO: File modes preserved from archive
}

// === Archive: Gz/Ungz Tests ===

func TestGz_CompressesFile(t *testing.T) {
	// TODO: Compress file, output smaller than input
}

func TestUngz_RestoresContent(t *testing.T) {
	// TODO: Decompress matches original content
}

// Gzip header contains original filename
func TestGz_PreservesFilename(t *testing.T) {
	temp := t.TempDir()
	srcPath := filepath.Join(temp, "sample.txt")
	if err := os.WriteFile(srcPath, []byte("hello"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	gzPath := filepath.Join(temp, "sample.txt.gz")
	if err := Gz(gzPath, File(srcPath)); err != nil {
		t.Fatalf("Gz failed: %v", err)
	}

	in, err := os.Open(gzPath)
	if err != nil {
		t.Fatalf("open gz failed: %v", err)
	}
	defer in.Close()

	gr, err := gzip.NewReader(in)
	if err != nil {
		t.Fatalf("gzip reader failed: %v", err)
	}
	defer gr.Close()

	if gr.Name != "sample.txt" {
		t.Fatalf("Expected gzip header name sample.txt, got %s", gr.Name)
	}
	data, err := io.ReadAll(gr)
	if err != nil || string(data) != "hello" {
		t.Fatalf("gzip content mismatch: %v", err)
	}
}

// === Archive: Tar Tests ===

func TestTar_SingleFile(t *testing.T) {
	// TODO: Archive single file to tar
}

func TestTar_Directory(t *testing.T) {
	// TODO: Archive directory with nested structure
}

func TestUntar_ExtractsCorrectly(t *testing.T) {
	// TODO: Extract tar restores file structure
}

func TestUntar_PathTraversalBlocked(t *testing.T) {
	// TODO: Reject malicious tar entries
}

func TestUntar_IgnoresNonRegular(t *testing.T) {
	// TODO: Symlinks/devices skipped safely
}

// === Archive: TarGz Tests ===

func TestTarGz_RoundTrip(t *testing.T) {
	// TODO: Create and extract tar.gz, verify content
}

func TestUntarGz_HandlesNested(t *testing.T) {
	// TODO: Extract nested directories from compressed tar
}
