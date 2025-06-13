//go:build !windows
// +build !windows

// This file provides stub implementations of windows-only functions
// to allow for cross-platform compilation.

package shell

func startDetachedProcess(listenCmd string) {
}
