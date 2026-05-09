package saga

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Zip compresses the provided files or directories into a .zip file
func Zip(dest string, files ...*FileHandle) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	w := zip.NewWriter(out)
	defer w.Close()

	for _, f := range files {
		err := filepath.Walk(f.path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			relPath := strings.TrimPrefix(strings.TrimPrefix(path, f.path), string(filepath.Separator))
			if relPath == "" {
				header.Name = info.Name()
			} else {
				header.Name = filepath.Join(info.Name(), relPath)
			}

			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := w.CreateHeader(header)
			if err != nil || info.IsDir() {
				return err
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Unzip extracts a .zip file to the destination directory
func Unzip(src *FileHandle, dest string) error {
	r, err := zip.OpenReader(src.path)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			continue // Prevent ZipSlip
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// Gz compresses a single file into a .gz file
func Gz(dest string, src *FileHandle) error {
	in, err := os.Open(src.path)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	w := gzip.NewWriter(out)
	defer w.Close()

	_, err = io.Copy(w, in)
	return err
}

// Ungz extracts a .gz file
func Ungz(src *FileHandle, dest string) error {
	in, err := os.Open(src.path)
	if err != nil {
		return err
	}
	defer in.Close()

	gr, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer gr.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, gr)
	return err
}

// Tar archives the provided files or directories into a .tar file
func Tar(dest string, files ...*FileHandle) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	return writeTar(out, files...)
}

// Untar extracts a .tar file to the destination directory
func Untar(src *FileHandle, dest string) error {
	in, err := os.Open(src.path)
	if err != nil {
		return err
	}
	defer in.Close()
	return readTar(in, dest)
}

// TarGz archives and compresses the provided files into a .tar.gz file
func TarGz(dest string, files ...*FileHandle) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	gw := gzip.NewWriter(out)
	defer gw.Close()

	return writeTar(gw, files...)
}

// UntarGz extracts a .tar.gz file to the destination directory
func UntarGz(src *FileHandle, dest string) error {
	in, err := os.Open(src.path)
	if err != nil {
		return err
	}
	defer in.Close()

	gr, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer gr.Close()

	return readTar(gr, dest)
}

func writeTar(w io.Writer, files ...*FileHandle) error {
	tw := tar.NewWriter(w)
	defer tw.Close()

	for _, f := range files {
		err := filepath.Walk(f.path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			relPath := strings.TrimPrefix(strings.TrimPrefix(path, f.path), string(filepath.Separator))
			if relPath == "" {
				header.Name = info.Name()
			} else {
				header.Name = filepath.Join(info.Name(), relPath)
			}

			if err := tw.WriteHeader(header); err != nil || info.IsDir() {
				return err
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tw, file)
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func readTar(r io.Reader, dest string) error {
	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name)
		if !strings.HasPrefix(target, filepath.Clean(dest)+string(os.PathSeparator)) {
			continue // Prevent TarSlip
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
	return nil
}
