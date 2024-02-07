package plugin

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
)

func installPackage(ctx context.Context, client *http.Client, version string, maxSize int64) error {
	// Sanitize user input
	if _, err := semver.NewVersion(version); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidTofuVersion, version)
	}

	err := downloadPackage(
		ctx,
		client,
		"/tmp/tofu.zip",
		fmt.Sprintf(
			"https://github.com/opentofu/opentofu/releases/download/%s/tofu_%s_linux_amd64.zip",
			version,
			strings.TrimPrefix(version, ""),
		),
	)
	if err != nil {
		return err
	}

	return unzip("/tmp/tofu.zip", "/bin", maxSize)
}

func downloadPackage(ctx context.Context, client *http.Client, filepath, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func unzip(src, dest string, maxSize int64) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	_ = os.MkdirAll(dest, defaultDirPerm)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path, err := sanitizeArchivePath(dest, f.Name)
		if err != nil {
			return err
		}

		if f.FileInfo().IsDir() { //nolint: nestif
			_ = os.MkdirAll(path, f.Mode())
		} else {
			_ = os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			written, err := io.CopyN(f, rc, maxSize)
			if err != nil {
				return err
			} else if written == maxSize {
				return fmt.Errorf("%w: %d", ErrMaxSizeSizeLimit, maxSize)
			}
		}

		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func sanitizeArchivePath(d, t string) (string, error) {
	value := filepath.Join(d, t)
	if strings.HasPrefix(value, filepath.Clean(d)) {
		return value, nil
	}

	return "", fmt.Errorf("%w: %v", ErrTaintedPath, t)
}

func deleteDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(path)
}
