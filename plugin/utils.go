package plugin

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"
)

func installPackage(ctx context.Context, client *http.Client, version string, maxSize int64) error {
	// Sanitize user input
	semverVersion, err := semver.NewVersion(version)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidTofuVersion, version)
	}

	packageURL := fmt.Sprintf(
		"https://github.com/opentofu/opentofu/releases/download/v%s/tofu_%s_linux_amd64.zip",
		semverVersion.String(),
		semverVersion.String(),
	)

	tmpdir, err := os.MkdirTemp("/tmp", "tofu_dl_")
	if err != nil {
		return fmt.Errorf("failed to create tmp dir: %w", err)
	}

	defer os.RemoveAll(tmpdir)

	log.Debug().
		Str("tmpdir", tmpdir).
		Msgf("Download OpenTofu '%s' from URL '%s'", version, packageURL)

	tmpfile := filepath.Join(tmpdir, "tofu.zip")

	if err := downloadPackage(ctx, client, tmpfile, packageURL); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	if err := unzip(tmpfile, tmpdir, maxSize); err != nil {
		return fmt.Errorf("failed to unzip: %w", err)
	}

	if err := os.Rename(filepath.Join(tmpdir, "tofu"), tofuBin); err != nil {
		return fmt.Errorf("failed to rename: %w", err)
	}

	return nil
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

	if resp.StatusCode > http.StatusBadRequest {
		return fmt.Errorf("%w: %v", ErrHTTPError, resp.Status)
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
			if err != nil && !errors.Is(err, io.EOF) {
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
