// saga is a package that handles file I/O, archiving, and other operations related to managing files in a structured way.
// It provides utilities for reading, writing, and organizing files, as well as creating archives and handling file metadata
package saga

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileHandle struct {
	path string
}

// File returns a File handle for the given path
//
// The path should be resolved before creating the handle, saga provides a Resolve utility for this purpose.
//
// The handle does not perform any I/O until methods are called on it.
func File(path string) *FileHandle {
	return &FileHandle{path: path}
}

// Path returns the file path
func (f *FileHandle) Path() string {
	return f.path
}

// File opens the file and returns an os.File handle
//
// The caller is responsible for closing the file when done
func (f *FileHandle) File() (*os.File, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %w", f.path, err)
	}
	return file, nil
}

// Stat returns the FileInfo for the file
func (f *FileHandle) Stat() (os.FileInfo, error) {
	info, err := os.Stat(f.path)
	if err != nil {
		return nil, fmt.Errorf("stat file %s: %w", f.path, err)
	}
	return info, nil
}

// Exists checks if the file exists
func (f *FileHandle) Exists() (bool, error) {
	_, err := f.Stat()
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// Size returns the file size in bytes
func (f *FileHandle) Size() (int64, error) {
	info, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// ModTime returns the file modification time
func (f *FileHandle) ModTime() (time.Time, error) {
	info, err := f.Stat()
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// Extension returns the file extension (e.g., ".txt")
func (f *FileHandle) Extension() string {
	return filepath.Ext(f.path)
}

// Name returns the base name of the file
func (f *FileHandle) Name() string {
	return filepath.Base(f.path)
}

// Read reads the entire file into memory
func (f *FileHandle) Read() ([]byte, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", f.path, err)
	}
	return data, nil
}

// Text reads the file as a string
func (f *FileHandle) Text() (string, error) {
	data, err := f.Read()
	return string(data), err
}

// JSON reads and unmarshals the file as JSON
func (f *FileHandle) JSON(v any) error {
	data, err := f.Read()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshal json from %s: %w", f.path, err)
	}
	return nil
}

// Write writes data to the file, creating directories if they don't exist
func (f *FileHandle) Write(data any) error {
	if err := os.MkdirAll(filepath.Dir(f.path), 0755); err != nil {
		return fmt.Errorf("create directories for %s: %w", f.path, err)
	}

	switch v := data.(type) {
	case string:
		if err := os.WriteFile(f.path, []byte(v), 0644); err != nil {
			return fmt.Errorf("write file %s: %w", f.path, err)
		}
		return nil
	case []byte:
		if err := os.WriteFile(f.path, v, 0644); err != nil {
			return fmt.Errorf("write file %s: %w", f.path, err)
		}
		return nil
	case io.Reader:
		// Stream reader to disk to avoid holding the full payload in memory.
		file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("open file %s for writing: %w", f.path, err)
		}
		defer file.Close()
		if _, err := io.Copy(file, v); err != nil {
			return fmt.Errorf("write file %s from reader: %w", f.path, err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported write type: %T", data)
	}
}

// WriteText writes a string to the file
func (f *FileHandle) WriteText(text string) error {
	return f.Write(text)
}

// WriteJSON marshals and writes data as JSON
func (f *FileHandle) WriteJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json for %s: %w", f.path, err)
	}
	return f.Write(data)
}

// Append appends data to the file
func (f *FileHandle) Append(data []byte) error {
	file, err := os.OpenFile(f.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file %s for append: %w", f.path, err)
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("append to file %s: %w", f.path, err)
	}
	return nil
}

// Delete removes the file
func (f *FileHandle) Delete() error {
	if err := os.Remove(f.path); err != nil {
		return fmt.Errorf("delete file %s: %w", f.path, err)
	}
	return nil
}

// Copy copies the file to a new location
func (f *FileHandle) Copy(dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("create directories for %s: %w", dest, err)
	}

	src, err := os.Open(f.path)
	if err != nil {
		return fmt.Errorf("open source %s: %w", f.path, err)
	}
	defer src.Close()

	dst, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create destination %s: %w", dest, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copy %s to %s: %w", f.path, dest, err)
	}
	return nil
}

// Move moves the file to a new location and updates the handle path
func (f *FileHandle) Move(dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("create directories for %s: %w", dest, err)
	}

	if err := os.Rename(f.path, dest); err == nil {
		f.path = dest
		return nil
	}

	// cross-device link or other error – try copy+delete
	if err := f.Copy(dest); err != nil {
		return err
	}
	if err := f.Delete(); err != nil {
		return err
	}
	f.path = dest
	return nil
}

// Reader returns an io.ReadCloser for the file
func (f *FileHandle) Reader() (io.ReadCloser, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %w", f.path, err)
	}
	return file, nil
}

// Bytes reads the file into a bytes.Buffer
func (f *FileHandle) Bytes() (*bytes.Buffer, error) {
	data, err := f.Read()
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

// Hash returns SHA256 hash of file contents
func (f *FileHandle) Hash() (string, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return "", fmt.Errorf("open file %s: %w", f.path, err)
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("hash file %s: %w", f.path, err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
