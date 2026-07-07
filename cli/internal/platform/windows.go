//go:build windows

package platform

import (
	"golang.org/x/sys/windows"
)

// windowsVTEnabled returns true if Virtual Terminal Processing is
// enabled on the Windows console, which is required for ANSI colors.
func windowsVTEnabled() bool {
	handle := windows.Stdout
	var mode uint32
	err := windows.GetConsoleMode(handle, &mode)
	if err != nil {
		return false
	}
	// Try to enable ENABLE_VIRTUAL_TERMINAL_PROCESSING (0x0004)
	if mode&0x0004 == 0 {
		mode |= 0x0004
		_ = windows.SetConsoleMode(handle, mode)
		_ = windows.GetConsoleMode(handle, &mode)
	}
	return mode&0x0004 != 0
}
