package plugin

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"
	"github.com/thegeeklab/wp-opentofu/tofu"
)

func installPackage(ctx context.Context, client *http.Client, version string) error {
	// Sanitize user input
	semverVersion, err := semver.NewVersion(version)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidTofuVersion, version)
	}

	packageURL := fmt.Sprintf(
		"https://github.com/opentofu/opentofu/releases/download/v%s/tofu_%s_linux_%s.zip",
		semverVersion.String(),
		semverVersion.String(),
		runtime.GOARCH,
	)

	tmpdir, err := os.MkdirTemp("/tmp", "tofu_dl_")
	if err != nil {
		return fmt.Errorf("failed to create tmp dir: %w", err)
	}

	defer func() {
		_ = os.RemoveAll(tmpdir)
	}()

	log.Debug().
		Str("tmpdir", tmpdir).
		Msgf("Download OpenTofu '%s' from URL '%s'", version, packageURL)

	tmpfile := filepath.Join(tmpdir, "tofu.zip")

	if err := downloadPackage(ctx, client, tmpfile, packageURL); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	if err := unzip(tmpfile, tmpdir); err != nil {
		return fmt.Errorf("failed to unzip: %w", err)
	}

	if err := os.Rename(filepath.Join(tmpdir, "tofu"), tofu.TofuBin); err != nil {
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

func unzip(src, dest string) error {
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

			if f.FileInfo().Size() == 0 {
				return nil
			}

			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			defer func() {
				if err := outFile.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.CopyN(outFile, rc, f.FileInfo().Size())
			if err != nil {
				return err
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

// boolPtr returns a pointer to the bool value.
func boolPtr(b bool) *bool {
	return &b
}

// parseStringMapToInitOptions converts a string map to InitOptions struct.
func parseStringMapToInitOptions(optionMap map[string]string) tofu.InitOptions {
	initOptions := tofu.InitOptions{}

	for key, value := range optionMap {
		switch key {
		case "backend":
			if parsedBool, err := strconv.ParseBool(value); err == nil {
				initOptions.Backend = boolPtr(parsedBool)
			}
		case "backend-config":
			if value != "" {
				initOptions.BackendConfig = []string{value}
			}
		case "lock":
			if parsedBool, err := strconv.ParseBool(value); err == nil {
				initOptions.Lock = boolPtr(parsedBool)
			}
		case "lockfile":
			if value != "" {
				initOptions.Lockfile = value
			}
		case "lock-timeout":
			if value != "" {
				initOptions.LockTimeout = value
			}
		}
	}

	return initOptions
}

// parseStringMapToFmtOptions converts a string map to FmtOptions struct.
func parseStringMapToFmtOptions(optionMap map[string]string) tofu.FmtOptions {
	fmtOptions := tofu.FmtOptions{}

	for key, value := range optionMap {
		switch key {
		case "list":
			if parsedBool, err := strconv.ParseBool(value); err == nil {
				fmtOptions.List = boolPtr(parsedBool)
			}
		case "write":
			if parsedBool, err := strconv.ParseBool(value); err == nil {
				fmtOptions.Write = boolPtr(parsedBool)
			}
		case "diff":
			if parsedBool, err := strconv.ParseBool(value); err == nil {
				fmtOptions.Diff = boolPtr(parsedBool)
			}
		case "check":
			if parsedBool, err := strconv.ParseBool(value); err == nil {
				fmtOptions.Check = boolPtr(parsedBool)
			}
		}
	}

	return fmtOptions
}
