package saga

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

// TODO: refactor to follow patterns below, smaller units.
func TestFileHandle_WriteAndRead(t *testing.T) {
	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, "test.txt")
	f := File(filePath)

	// Test WriteText
	err := f.WriteText("hello world")
	if err != nil {
		t.Fatalf("WriteText failed: %v\n", err)
	}

	// Test Exists
	exists, err := f.Exists()
	if err != nil || !exists {
		t.Fatalf("Expected file to exist\n")
	}

	// Test Text (Read)
	content, err := f.Text()
	if err != nil || content != "hello world" {
		t.Fatalf("Expected 'hello world', got %q\n", content)
	}

	// Test Size
	size, err := f.Size()
	if err != nil || size != 11 {
		t.Fatalf("Expected size 11, got %d\n", size)
	}
}

func TestFileHandle_JSON(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "data.json")
	f := File(filePath)

	data := map[string]string{"key": "value"}
	if err := f.WriteJSON(data); err != nil {
		t.Fatalf("WriteJSON failed: %v\n", err)
	}

	var result map[string]string
	if err := f.JSON(&result); err != nil || result["key"] != "value" {
		t.Fatalf("JSON read failed or mismatched data\n")
	}
}

// === FileHandle Tests ===

// Verify File() returns handle with correct path
func TestFile_ReturnsHandle(t *testing.T) {
	f := File("test.txt")
	if f.Path() != "test.txt" {
		t.Fatalf("Expected File path to be 'test.txt', got '%s'\n", f.Path())
	}
}

// Read file that exists, verify content matches via Text()
func TestFile_Read_ExistingFile(t *testing.T) {
	f := createTempFile(t, "test content")

	content, err := f.Text()
	if err != nil {
		t.Fatalf("Failed to read file: %v\n", err)
	}
	if content != "test content" {
		t.Fatalf("Expected 'test content', got '%s'\n", content)
	}
}

// Read file that doesn't exist, expect error
func TestFile_Read_NonexistentFile(t *testing.T) {
	f_name := strconv.FormatInt(time.Now().UnixNano(), 10) + "_nonexistent.txt"
	f := File(path.Join(t.TempDir(), f_name))

	_, err := f.Text()
	if err == nil {
		t.Fatalf("Expected error when reading nonexistent file\n")
	}
}

// JSON() unmarshals valid JSON into struct
func TestFile_JSON_UnmarshalsValid(t *testing.T) {
	f := createTempFile(t, `{"name": "saga", "version": 1}`)

	var result struct {
		Name    string `json:"name"`
		Version int    `json:"version"`
	}
	if err := f.JSON(&result); err != nil {
		t.Fatalf("JSON() returned error: %v\n", err)
	}
	if result.Name != "saga" || result.Version != 1 {
		t.Fatalf("JSON() unmarshaled incorrect data: %+v\n", result)
	}
}

// JSON() returns error for malformed JSON
func TestFile_JSON_InvalidJSON(t *testing.T) {
	f := createTempFile(t, `{"name": "saga", "version": 1`) // Missing closing brace

	var result struct {
		Name    string `json:"name"`
		Version int    `json:"version"`
	}
	if err := f.JSON(&result); err == nil {
		t.Fatalf("Expected error for invalid JSON, got none\n")
	}
}

// Bytes() returns populated buffer
func TestFile_Bytes_ReturnsBuffer(t *testing.T) {
	f := createTempFile(t, "buffer content")

	buf, err := f.Bytes()
	if err != nil {
		t.Fatalf("Bytes() returned error: %v\n", err)
	}
	if buf.String() != "buffer content" {
		t.Fatalf("Expected 'buffer content', got '%s'\n", buf.String())
	}
}

// Hash() returns consistent SHA256 for same content
func TestFile_Hash_Consistent(t *testing.T) {
	content := "hash me"
	f1 := createTempFile(t, content)
	f2 := createTempFile(t, content)

	hash1, err1 := f1.Hash()
	hash2, err2 := f2.Hash()

	t.Logf("\nHash1: %s\nHash2: %s\n", hash1, hash2)

	if err1 != nil || err2 != nil {
		t.Fatalf("Hash() returned error: %v, %v\n", err1, err2)
	}
	if hash1 != hash2 {
		t.Fatalf("Expected identical hashes for same content, got %s and %s\n", hash1, hash2)
	}
}

// Write string content, verify file created
func TestFile_Write_String(t *testing.T) {
	f_content := "hello saga"
	f := File(path.Join(t.TempDir(), "write_test.txt"))
	err := f.WriteText(f_content)
	if err != nil {
		t.Fatalf("WriteText failed: %v\n", err)
	}

	exists, err := f.Exists()
	if err != nil || !exists {
		t.Fatalf("Expected file to exist after WriteText\n")
	}

	content, err := f.Text()
	if err != nil || content != f_content {
		t.Fatalf("Expected 'hello saga', got '%s'\n", content)
	}
}

func TestFile_Write_ByteSlice(t *testing.T) {
	// TODO: Write []byte, verify file created
}

func TestFile_Write_IOReader(t *testing.T) {
	// TODO: Stream from io.Reader to file
}

func TestFile_Write_CreatesDirectories(t *testing.T) {
	// TODO: Write to nested path creates parent dirs
}

func TestFile_WriteJSON_FormatsOutput(t *testing.T) {
	// TODO: WriteJSON() produces indented JSON
}

func TestFile_Append_AddsToFile(t *testing.T) {
	// TODO: Append() adds data to existing file
}

func TestFile_Copy_CreatesIdentical(t *testing.T) {
	// TODO: Copy produces identical file at destination
}

func TestFile_Copy_CreatesDirectories(t *testing.T) {
	// TODO: Copy creates parent dirs if needed
}

func TestFile_Move_UpdatesPath(t *testing.T) {
	// TODO: Move relocates file and updates handle path
}

func TestFile_Move_FallbackCopyDelete(t *testing.T) {
	// TODO: Move handles cross-device via copy+delete
}

func TestFile_Delete_RemovesFile(t *testing.T) {
	// TODO: Delete removes file from filesystem
}

// Exists() returns true for existing file
func TestFile_Exists_TrueForExisting(t *testing.T) {
	f := File(path.Join(t.TempDir(), "exists.txt"))
	if err := f.WriteText("test"); err != nil {
		t.Fatalf("Failed to write file for Exists test: %v", err)
	}
	t.Logf("Testing Exists() for existing file: %s\n", f.Path())

	exists, err := f.Exists()
	if err != nil {
		t.Fatalf("Exists() returned error: %v", err)
	}
	if !exists {
		t.Fatalf("Expected file to exist, but Exists() returned false")
	}
}

// Exists() returns false for missing file
func TestFile_Exists_FalseForMissing(t *testing.T) {
	f_name := strconv.FormatInt(time.Now().UnixNano(), 10) + "_missing.txt"
	f := File(path.Join(t.TempDir(), f_name))
	t.Logf("Testing Exists() for missing file: %s\n", f.Path())

	exists, err := f.Exists()
	if err != nil {
		t.Fatalf("Exists() returned error: %v", err)
	}
	if exists {
		t.Fatalf("Expected file to not exist, but Exists() returned true")
	}
}

func TestFile_Stat_ReturnsInfo(t *testing.T) {
	// TODO: Stat() returns valid FileInfo
}

func TestFile_Size_MatchesContent(t *testing.T) {
	// TODO: Size() matches actual content length
}

func TestFile_ModTime_ReturnsTime(t *testing.T) {
	// TODO: ModTime() returns valid timestamp
}

func TestFile_Extension_Parses(t *testing.T) {
	// TODO: Extension() returns correct extension
}

func TestFile_Name_ReturnsBase(t *testing.T) {
	// TODO: Name() returns filename without path
}

// === Edge Cases & Error Handling ===

func TestFile_PermissionDenied(t *testing.T) {
	// TODO: Read/write on read-only file returns error
}

func TestFile_EmptyPath(t *testing.T) {
	// TODO: Operations on empty path return meaningful error
}

func TestFile_Reader_ClosedProperly(t *testing.T) {
	// TODO: Verify file handles closed after Reader() use
}

func TestFile_ConcurrentAccess(t *testing.T) {
	// TODO: Multiple goroutines reading same file handle
}

func TestFile_EmptyFile(t *testing.T) {
	// TODO: Hash, read, copy empty file succeeds
}

func TestFile_Write_LargeFile(t *testing.T) {
	// TODO: Write(io.Reader) streams without loading all in memory
}

// === Test Helpers ===

func createTempFile(t *testing.T, content string) *FileHandle {
	t.Helper()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "testfile")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	return File(path)
}
