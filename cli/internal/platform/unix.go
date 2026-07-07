//go:build !windows

package platform

// windowsVTEnabled is only meaningful on Windows.
// On other platforms, return true — color is assumed if we reach this.
func windowsVTEnabled() bool {
	return true
}
