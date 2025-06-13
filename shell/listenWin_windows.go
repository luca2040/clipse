//go:build windows
// +build windows

// This file provides stub implementations of windows-only functions
// to allow for cross-platform compilation.

package shell

import (
	"os"
	"os/exec"
	"syscall"
)

const DETACHED_PROCESS = 0x00000008

func startDetachedProcess(listenCmd string) {
	devNull, _ := os.Open(os.DevNull)
	defer devNull.Close()

	cmd := exec.Command(os.Args[0], listenCmd)
	cmd.Stdout = devNull
	cmd.Stderr = devNull

	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | DETACHED_PROCESS,
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}
}
