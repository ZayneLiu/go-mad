// saga is a package that handle file I/O, archiving, and other operations related to managing files in a structured way.
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
	"time"
)

type FileHandle struct {
	path string
	info os.FileInfo
}

// File returns a File handle for the given path
func File(path string) *FileHandle {
	return &FileHandle{path: path}
}

// Path returns the file path
func (f *FileHandle) Path() string {
	return f.path
}

// Exists checks if the file exists
func (f *FileHandle) Exists() (bool, error) {
	_, err := os.Stat(f.path)
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
	info, err := os.Stat(f.path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// ModTime returns the file modification time
func (f *FileHandle) ModTime() (time.Time, error) {
	info, err := os.Stat(f.path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// Read reads the entire file into memory
func (f *FileHandle) Read() ([]byte, error) {
	return os.ReadFile(f.path)
}

// Text reads the file as a string
func (f *FileHandle) Text() (string, error) {
	data, err := f.Read()
	return string(data), err
}

// JSON reads and unmarshals the file as JSON
func (f *FileHandle) JSON(v interface{}) error {
	data, err := f.Read()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// JSONString is a convenience for getting JSON as a string
func (f *FileHandle) JSONString() (string, error) {
	return f.Text()
}

func (f *FileHandle) Write(data any) error {
	switch v := data.(type) {
	case string:
		return os.WriteFile(f.path, []byte(v), 0644)
	case []byte:
		return os.WriteFile(f.path, v, 0644)
	case io.Reader:
		file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, v)
		return err
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
		return err
	}
	return f.Write(data)
}

// Append appends data to the file
func (f *FileHandle) Append(data []byte) error {
	file, err := os.OpenFile(f.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}

// Delete removes the file
func (f *FileHandle) Delete() error {
	return os.Remove(f.path)
}

// Copy copies the file to a new location
func (f *FileHandle) Copy(dest string) error {
	src, err := os.Open(f.path)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// Move moves the file to a new location
func (f *FileHandle) Move(dest string) error {
	return os.Rename(f.path, dest)
}

// Reader returns an io.ReadCloser for the file
func (f *FileHandle) Reader() (io.ReadCloser, error) {
	return os.Open(f.path)
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
		return "", err
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
