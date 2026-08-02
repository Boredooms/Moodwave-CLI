package playback

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// FFplayStatus describes the current state of ffplay availability.
type FFplayStatus struct {
	Available bool
	Path      string
	Source    string // "system", "cache", "none"
}

// CheckFFplay returns the status of ffplay on the current system.
func CheckFFplay() FFplayStatus {
	execName := "ffplay"
	if runtime.GOOS == "windows" {
		execName = "ffplay.exe"
	}

	// Check system PATH
	if p, err := exec.LookPath(execName); err == nil {
		return FFplayStatus{Available: true, Path: p, Source: "system"}
	}

	// Check moodwave cache
	if cacheDir, err := os.UserCacheDir(); err == nil {
		p := filepath.Join(cacheDir, "moodwave", execName)
		if _, err := os.Stat(p); err == nil {
			return FFplayStatus{Available: true, Path: p, Source: "cache"}
		}
	}

	return FFplayStatus{Available: false, Source: "none"}
}

// InstallFFplay downloads and installs ffplay to the moodwave cache.
// It supports Windows (zip), macOS (zip), and Linux (tar.gz).
// Returns the installed path or an error.
func InstallFFplay(ctx context.Context, progressFn func(msg string)) (string, error) {
	if progressFn == nil {
		progressFn = func(string) {}
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine cache dir: %w", err)
	}
	mwCache := filepath.Join(cacheDir, "moodwave")
	_ = os.MkdirAll(mwCache, 0755)

	execName := "ffplay"
	if runtime.GOOS == "windows" {
		execName = "ffplay.exe"
	}
	targetPath := filepath.Join(mwCache, execName)

	// Determine download URL based on OS/Arch
	dlURL, archiveType, err := ffplayDownloadURL()
	if err != nil {
		return "", err
	}

	progressFn(fmt.Sprintf("Downloading ffplay for %s/%s...", runtime.GOOS, runtime.GOARCH))

	// Download
	client := &http.Client{Timeout: 5 * time.Minute}
	req, err := http.NewRequestWithContext(ctx, "GET", dlURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "moodwave-cli/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	// Save to temp file
	tmpFile := filepath.Join(mwCache, "ffmpeg-download.tmp")
	out, err := os.Create(tmpFile)
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}

	progressFn("Downloading... (this may take a minute)")
	written, err := io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		os.Remove(tmpFile)
		return "", fmt.Errorf("download write failed: %w", err)
	}
	progressFn(fmt.Sprintf("Downloaded %.1f MB, extracting...", float64(written)/1024/1024))

	// Extract ffplay from archive
	extractDir := filepath.Join(mwCache, "ffmpeg-extract")
	os.RemoveAll(extractDir)
	_ = os.MkdirAll(extractDir, 0755)

	switch archiveType {
	case "zip":
		err = extractFromZip(tmpFile, extractDir, execName)
	case "tar.gz":
		err = extractFromTarGz(tmpFile, extractDir, execName)
	default:
		err = fmt.Errorf("unsupported archive type: %s", archiveType)
	}

	os.Remove(tmpFile)
	if err != nil {
		os.RemoveAll(extractDir)
		return "", fmt.Errorf("extraction failed: %w", err)
	}

	// Find the extracted ffplay binary
	var foundPath string
	_ = filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.EqualFold(info.Name(), execName) {
			foundPath = path
			return filepath.SkipAll
		}
		return nil
	})

	if foundPath == "" {
		os.RemoveAll(extractDir)
		return "", fmt.Errorf("ffplay not found in archive")
	}

	// Move to final location
	if err := copyFile(foundPath, targetPath); err != nil {
		os.RemoveAll(extractDir)
		return "", fmt.Errorf("install failed: %w", err)
	}

	// Make executable (unix)
	if runtime.GOOS != "windows" {
		os.Chmod(targetPath, 0755)
	}

	// Also extract ffmpeg and ffprobe if present
	for _, name := range []string{"ffmpeg", "ffprobe"} {
		n := name
		if runtime.GOOS == "windows" {
			n += ".exe"
		}
		_ = filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() && strings.EqualFold(info.Name(), n) {
				dest := filepath.Join(mwCache, n)
				_ = copyFile(path, dest)
				if runtime.GOOS != "windows" {
					os.Chmod(dest, 0755)
				}
				return filepath.SkipAll
			}
			return nil
		})
	}

	// Cleanup
	os.RemoveAll(extractDir)
	progressFn("ffplay installed successfully")

	return targetPath, nil
}

// ffplayDownloadURL returns the appropriate download URL for the current platform.
func ffplayDownloadURL() (url string, archiveType string, err error) {
	switch runtime.GOOS {
	case "windows":
		// gyan.dev essentials build — static, no DLL dependencies
		return "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip", "zip", nil
	case "darwin":
		// evermeet.cx provides macOS static builds
		if runtime.GOARCH == "arm64" {
			return "https://evermeet.cx/ffmpeg/getrelease/zip/ffplay/arm64", "zip", nil
		}
		return "https://evermeet.cx/ffmpeg/getrelease/zip/ffplay", "zip", nil
	case "linux":
		// johnvansickle.com provides Linux static builds
		arch := "amd64"
		if runtime.GOARCH == "arm64" {
			arch = "arm64"
		} else if runtime.GOARCH == "arm" {
			arch = "armhf"
		}
		return fmt.Sprintf("https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-%s-static.tar.xz", arch),
			"tar.gz", nil
	default:
		return "", "", fmt.Errorf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
	}
}

func extractFromZip(zipPath, destDir, targetName string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		// Extract all bin/ files (ffplay, ffmpeg, ffprobe)
		name := filepath.Base(f.Name)
		if !isBinaryFile(name) {
			continue
		}

		destPath := filepath.Join(destDir, name)
		rc, err := f.Open()
		if err != nil {
			continue
		}
		outFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			rc.Close()
			continue
		}
		_, _ = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
	}
	return nil
}

func extractFromTarGz(archivePath, destDir, targetName string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	var reader io.Reader
	// Try gzip first, then xz via pipe
	gz, gzErr := gzip.NewReader(f)
	if gzErr != nil {
		// May be xz — use xz command if available
		f.Close()
		return extractFromXz(archivePath, destDir, targetName)
	}
	defer gz.Close()
	reader = gz

	tr := tar.NewReader(reader)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		name := filepath.Base(hdr.Name)
		if !isBinaryFile(name) {
			continue
		}

		destPath := filepath.Join(destDir, name)
		outFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			continue
		}
		_, _ = io.Copy(outFile, tr)
		outFile.Close()
	}
	return nil
}

func extractFromXz(archivePath, destDir, targetName string) error {
	// Use system xz/tar if available (Linux typically has these)
	xzPath, err := exec.LookPath("xz")
	if err != nil {
		// Try unxz
		xzPath, err = exec.LookPath("unxz")
		if err != nil {
			return fmt.Errorf("xz not found — install xz-utils to extract")
		}
	}

	// Decompress .xz to .tar
	tarPath := archivePath + ".tar"
	cmd := exec.Command(xzPath, "-dk", archivePath)
	if err := cmd.Run(); err != nil {
		// Try alternative: pipe xz -d to tar
		cmd = exec.Command("tar", "-xJf", archivePath, "-C", destDir)
		return cmd.Run()
	}
	defer os.Remove(tarPath)

	// Extract from tar
	f, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer f.Close()

	tr := tar.NewReader(f)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		name := filepath.Base(hdr.Name)
		if !isBinaryFile(name) {
			continue
		}
		destPath := filepath.Join(destDir, name)
		outFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			continue
		}
		_, _ = io.Copy(outFile, tr)
		outFile.Close()
	}
	return nil
}

func isBinaryFile(name string) bool {
	lower := strings.ToLower(name)
	binaries := []string{"ffplay", "ffplay.exe", "ffmpeg", "ffmpeg.exe", "ffprobe", "ffprobe.exe"}
	for _, b := range binaries {
		if lower == b {
			return true
		}
	}
	return false
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
