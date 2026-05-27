package saga

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Zip compresses the provided files or directories into a .zip file.
// Entries are stored under the base name of each provided path.
func Zip(dest string, files ...*FileHandle) error {
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create zip destination %s: %w", dest, err)
	}
	defer out.Close()

	w := zip.NewWriter(out)
	defer w.Close()

	for _, f := range files {
		baseName := filepath.Base(f.path)
		err := filepath.Walk(f.path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("walk %s: %w", path, err)
			}
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return fmt.Errorf("create zip header for %s: %w", path, err)
			}

			// Use filepath.Rel to normalize and avoid fragile prefix trimming.
			relPath, err := filepath.Rel(f.path, path)
			if err != nil {
				return fmt.Errorf("resolve relative path for %s: %w", path, err)
			}
			if relPath == "." {
				header.Name = baseName
			} else {
				// Keep archive entries rooted under the base directory name.
				header.Name = filepath.Join(baseName, relPath)
			}

			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := w.CreateHeader(header)
			if err != nil {
				return fmt.Errorf("create zip entry for %s: %w", path, err)
			}
			if info.IsDir() {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("open %s: %w", path, err)
			}
			defer file.Close()
			if _, err = io.Copy(writer, file); err != nil {
				return fmt.Errorf("write %s to zip: %w", path, err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Unzip extracts a .zip file to the destination directory.
// Returns an error if any entry would escape the destination.
func Unzip(src *FileHandle, dest string) error {
	r, err := zip.OpenReader(src.path)
	if err != nil {
		return fmt.Errorf("open zip %s: %w", src.path, err)
	}
	defer r.Close()

	cleanDest := filepath.Clean(dest)
	// Reject any entry that would escape the destination.
	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, cleanDest+string(os.PathSeparator)) {
			// Explicitly fail on traversal attempts instead of silently skipping.
			return fmt.Errorf("path escapes destination: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return fmt.Errorf("create directory %s: %w", fpath, err)
			}
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return fmt.Errorf("create directory %s: %w", filepath.Dir(fpath), err)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("create file %s: %w", fpath, err)
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("open zip entry %s: %w", f.Name, err)
		}

		if _, err = io.Copy(outFile, rc); err != nil {
			outFile.Close()
			rc.Close()
			return fmt.Errorf("extract zip entry %s: %w", f.Name, err)
		}
		outFile.Close()
		rc.Close()
	}
	return nil
}

// Gz compresses a single file into a .gz file.
// The gzip header preserves the original filename.
func Gz(dest string, src *FileHandle) error {
	in, err := os.Open(src.path)
	if err != nil {
		return fmt.Errorf("open source %s: %w", src.path, err)
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create gzip destination %s: %w", dest, err)
	}
	defer out.Close()

	w := gzip.NewWriter(out)
	w.Name = filepath.Base(src.path)
	defer w.Close()

	if _, err = io.Copy(w, in); err != nil {
		return fmt.Errorf("write gzip %s: %w", dest, err)
	}
	return nil
}

// Ungz extracts a .gz file.
func Ungz(src *FileHandle, dest string) error {
	in, err := os.Open(src.path)
	if err != nil {
		return fmt.Errorf("open gzip source %s: %w", src.path, err)
	}
	defer in.Close()

	gr, err := gzip.NewReader(in)
	if err != nil {
		return fmt.Errorf("create gzip reader for %s: %w", src.path, err)
	}
	defer gr.Close()

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create gzip destination %s: %w", dest, err)
	}
	defer out.Close()

	if _, err = io.Copy(out, gr); err != nil {
		return fmt.Errorf("extract gzip to %s: %w", dest, err)
	}
	return nil
}

// Tar archives the provided files or directories into a .tar file.
// Entries are stored under the base name of each provided path.
func Tar(dest string, files ...*FileHandle) error {
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create tar destination %s: %w", dest, err)
	}
	defer out.Close()
	return writeTar(out, files...)
}

// Untar extracts a .tar file to the destination directory.
// Returns an error if any entry would escape the destination.
func Untar(src *FileHandle, dest string) error {
	in, err := os.Open(src.path)
	if err != nil {
		return fmt.Errorf("open tar source %s: %w", src.path, err)
	}
	defer in.Close()
	return readTar(in, dest)
}

// TarGz archives and compresses the provided files into a .tar.gz file.
// Entries are stored under the base name of each provided path.
func TarGz(dest string, files ...*FileHandle) error {
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create tar.gz destination %s: %w", dest, err)
	}
	defer out.Close()

	gw := gzip.NewWriter(out)
	defer gw.Close()

	return writeTar(gw, files...)
}

// UntarGz extracts a .tar.gz file to the destination directory.
// Returns an error if any entry would escape the destination.
func UntarGz(src *FileHandle, dest string) error {
	in, err := os.Open(src.path)
	if err != nil {
		return fmt.Errorf("open tar.gz source %s: %w", src.path, err)
	}
	defer in.Close()

	gr, err := gzip.NewReader(in)
	if err != nil {
		return fmt.Errorf("create gzip reader for %s: %w", src.path, err)
	}
	defer gr.Close()

	return readTar(gr, dest)
}

func writeTar(w io.Writer, files ...*FileHandle) error {
	tw := tar.NewWriter(w)
	defer tw.Close()

	for _, f := range files {
		baseName := filepath.Base(f.path)
		err := filepath.Walk(f.path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("walk %s: %w", path, err)
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return fmt.Errorf("create tar header for %s: %w", path, err)
			}

			// Use filepath.Rel to normalize and avoid fragile prefix trimming.
			relPath, err := filepath.Rel(f.path, path)
			if err != nil {
				return fmt.Errorf("resolve relative path for %s: %w", path, err)
			}
			if relPath == "." {
				header.Name = baseName
			} else {
				// Keep archive entries rooted under the base directory name.
				header.Name = filepath.Join(baseName, relPath)
			}

			if err := tw.WriteHeader(header); err != nil {
				return fmt.Errorf("write tar header for %s: %w", path, err)
			}
			if info.IsDir() {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("open %s: %w", path, err)
			}
			defer file.Close()
			if _, err = io.Copy(tw, file); err != nil {
				return fmt.Errorf("write %s to tar: %w", path, err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func readTar(r io.Reader, dest string) error {
	tr := tar.NewReader(r)
	cleanDest := filepath.Clean(dest)
	// Reject any entry that would escape the destination.

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read tar header: %w", err)
		}

		target := filepath.Join(dest, header.Name)
		if !strings.HasPrefix(target, cleanDest+string(os.PathSeparator)) {
			return fmt.Errorf("path escapes destination: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("create directory %s: %w", target, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("create directory %s: %w", filepath.Dir(target), err)
			}
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("create file %s: %w", target, err)
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("extract tar entry %s: %w", header.Name, err)
			}
			outFile.Close()
		default:
			// Ignore non-regular files (e.g. symlinks, devices) for safety.
		}
	}
	return nil
}
