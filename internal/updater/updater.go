package updater

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	repoOwner  = "ChristianLapinig"
	repoName   = "aem-local-cli"
	binaryName = "aemlocal"
	CacheTTL   = 24 * time.Hour
)

type CacheEntry struct {
	LatestVersion string    `json:"latest_version"`
	CheckedAt     time.Time `json:"checked_at"`
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func cacheFilePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "aemlocal", "version_check.json"), nil
}

func ReadCache() (*CacheEntry, error) {
	path, err := cacheFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

func writeCache(latest string) error {
	path, err := cacheFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	entry := CacheEntry{
		LatestVersion: latest,
		CheckedAt:     time.Now(),
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// FetchLatestVersion queries the GitHub releases API for the latest version tag.
func FetchLatestVersion() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}
	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	return release.TagName, nil
}

// RefreshCache fetches the latest version and updates the local cache.
// Safe to call from a goroutine; errors are silently ignored.
func RefreshCache() {
	latest, err := FetchLatestVersion()
	if err != nil {
		return
	}
	_ = writeCache(latest)
}

// CheckForUpdate returns the latest version string if a newer version is available,
// reading only from the local cache. Returns empty string if up-to-date, on dev builds,
// or if the cache is missing/stale.
func CheckForUpdate(currentVersion string) string {
	if currentVersion == "dev" {
		return ""
	}
	entry, err := ReadCache()
	if err != nil || time.Since(entry.CheckedAt) > CacheTTL {
		return ""
	}
	if IsNewer(entry.LatestVersion, currentVersion) {
		return entry.LatestVersion
	}
	return ""
}

// IsNewer returns true if latest is strictly greater than current.
// Both must be semver strings (e.g. "v1.2.3"). Returns false if either is unparseable.
func IsNewer(latest, current string) bool {
	l := parseSemver(latest)
	c := parseSemver(current)
	if l == nil || c == nil {
		return false
	}
	for i := 0; i < 3; i++ {
		if l[i] != c[i] {
			return l[i] > c[i]
		}
	}
	return false
}

func parseSemver(v string) []int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 {
		return nil
	}
	// Strip pre-release suffix from patch component (e.g. "1-beta" → "1")
	parts[2] = strings.SplitN(parts[2], "-", 2)[0]
	nums := make([]int, 3)
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil
		}
		nums[i] = n
	}
	return nums
}

// SelfUpdate downloads the release for the given version and replaces the running binary.
func SelfUpdate(version string) error {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	var archiveName string
	if goos == "windows" {
		archiveName = fmt.Sprintf("%s_%s_%s.zip", binaryName, goos, goarch)
	} else {
		archiveName = fmt.Sprintf("%s_%s_%s.tar.gz", binaryName, goos, goarch)
	}

	downloadURL := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
		repoOwner, repoName, version, archiveName)
	checksumURL := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/checksums.txt",
		repoOwner, repoName, version)

	tmpDir, err := os.MkdirTemp("", "aemlocal-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, archiveName)
	checksumPath := filepath.Join(tmpDir, "checksums.txt")

	fmt.Print("Downloading... ")
	if err := downloadFile(downloadURL, archivePath); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	fmt.Println("done")

	if err := downloadFile(checksumURL, checksumPath); err != nil {
		return fmt.Errorf("checksum download failed: %w", err)
	}
	if err := verifyChecksum(archivePath, archiveName, checksumPath); err != nil {
		return fmt.Errorf("checksum verification failed: %w", err)
	}

	var newBinaryPath string
	if goos == "windows" {
		newBinaryPath, err = extractZip(archivePath, tmpDir)
	} else {
		newBinaryPath, err = extractTarGz(archivePath, tmpDir)
	}
	if err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not determine executable path: %w", err)
	}

	return replaceExecutable(newBinaryPath, execPath)
}

func downloadFile(url, dest string) error {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func verifyChecksum(archivePath, archiveName, checksumPath string) error {
	data, err := os.ReadFile(checksumPath)
	if err != nil {
		return err
	}
	var expected string
	for _, line := range strings.Split(string(data), "\n") {
		// GoReleaser format: "<sha256>  <filename>"
		if strings.HasSuffix(strings.TrimSpace(line), archiveName) {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				expected = parts[0]
				break
			}
		}
	}
	if expected == "" {
		return fmt.Errorf("checksum for %s not found in checksums.txt", archiveName)
	}

	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	actual := fmt.Sprintf("%x", h.Sum(nil))
	if actual != expected {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expected, actual)
	}
	return nil
}

func extractTarGz(archivePath, destDir string) (string, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if header.Typeflag == tar.TypeReg && filepath.Base(header.Name) == binaryName {
			destPath := filepath.Join(destDir, binaryName)
			out, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				return "", err
			}
			_, copyErr := io.Copy(out, tr)
			out.Close()
			if copyErr != nil {
				return "", copyErr
			}
			return destPath, nil
		}
	}
	return "", fmt.Errorf("binary %q not found in archive", binaryName)
}

func extractZip(archivePath, destDir string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	target := binaryName + ".exe"
	for _, f := range r.File {
		if filepath.Base(f.Name) == target {
			destPath := filepath.Join(destDir, target)
			out, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				return "", err
			}
			rc, err := f.Open()
			if err != nil {
				out.Close()
				return "", err
			}
			_, copyErr := io.Copy(out, rc)
			rc.Close()
			out.Close()
			if copyErr != nil {
				return "", copyErr
			}
			return destPath, nil
		}
	}
	return "", fmt.Errorf("binary %q not found in archive", target)
}

func replaceExecutable(newPath, execPath string) error {
	execDir := filepath.Dir(execPath)

	// Temp file in the same directory ensures the rename is on the same filesystem (atomic).
	tmp, err := os.CreateTemp(execDir, ".aemlocal-update-*")
	if err != nil {
		return fmt.Errorf("cannot write to %s — try: sudo aemlocal update", execDir)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	if err := copyFile(newPath, tmp.Name()); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), 0755); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), execPath)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
