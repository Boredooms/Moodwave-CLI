# Moodwave Handoff Guide for Claude 🚀

This repository has been structured and prepared for you to take over and build the final product landing page, deployment pipeline, and release automated distribution.

## 📂 Repository Structure

- **`/cli`**: Contains the complete Go implementation of the Moodwave music player.
  - `/cli/cmd/moodwave`: Main entry point and CLI commands (e.g. `play`, `scan`, `search`, `init`, `welcome`, `update`).
  - `/cli/internal`: The core packages (concurrency-safe playback controllers, scanner engine, recommender, and visualizers).
  - `/cli/scripts`: The standalone shell (`install.sh`) and PowerShell (`install.ps1`) installation scripts.
- **`/website`**: Contains the static files for the product landing page (`index.html` and `style.css`).

---

## ⚡ How the Zero-Dependency Installer Works

To allow users to run Moodwave without having Go installed:
1. When a user runs the one-line install command, it downloads the installer script from the `main` branch.
2. The installer script detects the user's OS and CPU architecture.
3. It makes a request to the GitHub Releases API (`api.github.com/repos/Boredooms/Moodwave-CLI/releases/latest`) to discover the latest version.
4. It downloads the pre-compiled binary matching their platform (e.g., `moodwave-windows-amd64.exe` or `moodwave-darwin-arm64`).
5. It extracts the binary, moves it into a persistent path (e.g., `~/.local/bin` or `$env:LOCALAPPDATA\moodwave\bin`), and appends that directory to the system `PATH`.

---

## 📋 What Needs to Be Done Next

### 1. Set Up GitHub Actions for Automated Releases
Create a `.github/workflows/release.yml` file to compile the Go binary for multiple platforms when a tag is pushed (e.g., `v1.0.1`):
- Targets:
  - Windows (`GOOS=windows GOARCH=amd64` & `GOOS=windows GOARCH=arm64`)
  - Linux (`GOOS=linux GOARCH=amd64` & `GOOS=linux GOARCH=arm64` & `GOOS=linux GOARCH=arm`)
  - macOS (`GOOS=darwin GOARCH=amd64` & `GOOS=darwin GOARCH=arm64`)
- Built binaries should be named matching the patterns:
  - `moodwave-windows-amd64.exe`
  - `moodwave-linux-amd64`
  - `moodwave-darwin-arm64`
- Publish them as assets on a new draft release matching the tag.

### 2. Deploy the Website
Host the files inside the `/website` folder on a static hosting provider.
- Recommended: **GitHub Pages**, **Vercel**, or **Netlify**.
- Ensure the custom domain (if any) is configured.

---

## 🤖 Handoff Prompt for Claude

If you are a Claude assistant taking over this project, copy and paste this prompt to start:

```markdown
Hi Claude! We are building Moodwave, a terminal-based mood-music companion that scans the codebase, infers the developer's working mood, and streams matching music (via YouTube, Radio Browser, etc.).

The codebase is organized into:
1. `/cli` - Go implementation of the CLI.
2. `/website` - Static website files (index.html, style.css) for the product landing page.

Please perform the following tasks:
1. Create a GitHub Actions workflow `.github/workflows/release.yml` that automatically cross-compiles the Go binaries in `cli/cmd/moodwave` for Windows (amd64, arm64), macOS (amd64, arm64), and Linux (amd64, arm64, arm) whenever a new git tag (e.g., v*) is pushed. The built binaries must be uploaded to the release draft as assets with names like `moodwave-[os]-[arch]`.
2. Inspect `cli/scripts/install.sh` and `cli/scripts/install.ps1` to make sure they match the release asset names perfectly and work smoothly.
3. Deploy the `/website` folder to GitHub Pages by setting up a GitHub Actions workflow or configuring the repository settings so the landing page is live.
4. Improve or enhance the landing page design if you can make it even more visually stunning.
```
