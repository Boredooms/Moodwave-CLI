//go:build windows

package playback

import (
	"os/exec"
	"syscall"
	"unsafe"
)

const (
	TH32CS_SNAPTHREAD     = 0x00000004
	THREAD_SUSPEND_RESUME = 0x0002
)

type THREADENTRY32 struct {
	Size           uint32
	Usage          uint32
	ThreadID       uint32
	OwnerProcessID uint32
	BasePri        int32
	TpBasePri      int32
	Flags          uint32
}

func pauseProcess(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	setProcessSuspended(cmd.Process.Pid, true)
}

func resumeProcess(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	setProcessSuspended(cmd.Process.Pid, false)
}

func setProcessSuspended(pid int, suspend bool) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	pCreateToolhelp32Snapshot := kernel32.NewProc("CreateToolhelp32Snapshot")
	pThread32First := kernel32.NewProc("Thread32First")
	pThread32Next := kernel32.NewProc("Thread32Next")
	pOpenThread := kernel32.NewProc("OpenThread")
	pSuspendThread := kernel32.NewProc("SuspendThread")
	pResumeThread := kernel32.NewProc("ResumeThread")
	pCloseHandle := kernel32.NewProc("CloseHandle")

	snapshot, _, _ := pCreateToolhelp32Snapshot.Call(TH32CS_SNAPTHREAD, 0)
	if snapshot == uintptr(syscall.InvalidHandle) {
		return
	}
	defer pCloseHandle.Call(snapshot)

	var te THREADENTRY32
	te.Size = uint32(unsafe.Sizeof(te))

	ret, _, _ := pThread32First.Call(snapshot, uintptr(unsafe.Pointer(&te)))
	for ret != 0 {
		if te.OwnerProcessID == uint32(pid) {
			threadHandle, _, _ := pOpenThread.Call(THREAD_SUSPEND_RESUME, 0, uintptr(te.ThreadID))
			if threadHandle != 0 {
				if suspend {
					pSuspendThread.Call(threadHandle)
				} else {
					pResumeThread.Call(threadHandle)
				}
				pCloseHandle.Call(threadHandle)
			}
		}
		ret, _, _ = pThread32Next.Call(snapshot, uintptr(unsafe.Pointer(&te)))
	}
}
