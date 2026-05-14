# saga - File Operations

A comprehensive file I/O library providing utilities for reading, writing, archiving, and managing files with support for multiple archive formats.

## Features

- **File Handles**: Simple abstraction over file paths with convenient methods
- **Read/Write Operations**: Support for plain text, JSON, and binary data
- **File Metadata**: Check existence, size, and modification time
- **Copy & Move**: File duplication and relocation
- **Append Operations**: Add content to existing files
- **Archive Support**:
  - ZIP compression/decompression
  - GZIP compression/decompression
  - TAR archiving/extraction
  - TAR.GZ compression/decompression
- **Hashing**: Generate SHA256 hashes of file contents
- **Security**: Protection against ZIP slip and TAR slip attacks

## File Handle Operations

**Getting a File Handle:**

```go
file := saga.File("/path/to/file.txt")
```

**Reading Files:**

```go
// Read as bytes
data, err := file.Read()

// Read as string
text, err := file.Text()

// Read as JSON
var data MyStruct
err := file.JSON(&data)

// Get file path
path := file.Path()
```

**Checking File Properties:**

```go
// Check if file exists
exists, err := file.Exists()

// Get file size
size, err := file.Size()

// Get modification time
modTime, err := file.ModTime()

// Get SHA256 hash
hash, err := file.Hash()
```

**Writing Files:**

```go
// Write text
err := file.WriteText("Hello, World!")

// Write JSON
err := file.WriteJSON(map[string]string{"key": "value"})

// Write bytes
err := file.Write([]byte("data"))

// Write from reader
err := file.Write(io.Reader)
```

**File Operations:**

```go
// Append data
err := file.Append([]byte("\nNew line"))

// Copy file
err := file.Copy("/path/to/dest")

// Move file
err := file.Move("/path/to/new/location")

// Delete file
err := file.Delete()

// Get reader
reader, err := file.Reader()

// Get bytes buffer
buffer, err := file.Bytes()
```

## Archive Operations

**ZIP Operations:**

```go
// Create a ZIP archive
file1 := saga.File("file1.txt")
file2 := saga.File("directory/")
err := saga.Zip("archive.zip", file1, file2)

// Extract a ZIP archive
archive := saga.File("archive.zip")
err := saga.Unzip(archive, "destination/")
```

**GZIP Operations:**

```go
// Compress a single file
file := saga.File("data.txt")
err := saga.Gz("data.txt.gz", file)

// Decompress
archive := saga.File("data.txt.gz")
err := saga.Ungz(archive, "data.txt")
```

**TAR Operations:**

```go
// Create a TAR archive
file1 := saga.File("file1.txt")
file2 := saga.File("directory/")
err := saga.Tar("archive.tar", file1, file2)

// Extract a TAR archive
archive := saga.File("archive.tar")
err := saga.Untar(archive, "destination/")
```

**TAR.GZ Operations:**

```go
// Create compressed TAR archive
file1 := saga.File("file1.txt")
file2 := saga.File("directory/")
err := saga.TarGz("archive.tar.gz", file1, file2)

// Extract compressed TAR archive
archive := saga.File("archive.tar.gz")
err := saga.UntarGz(archive, "destination/")
```

## Utility Functions

**Path Resolution:**

```go
// Resolve paths with home directory expansion
resolvedPath, err := saga.ResolvePath("~/documents/file.txt")
```

## Security Features

- **ZIP Slip Protection**: Validates extracted file paths to prevent directory traversal
- **TAR Slip Protection**: Validates extracted file paths in TAR archives
