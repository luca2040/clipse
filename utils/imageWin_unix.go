//go:build !windows
// +build !windows

// This file provides stub implementations of windows-only functions
// to allow for cross-platform compilation.

package utils

func HasClipboardContentChanged() bool {
	return false
}

func ClipboardHasImage() (bool, string) {
	return false, ""
}

func SavePNGStringToFile(pngStr string, filePath string) error {
	return nil
}

func CopyImageToClipboard(imagePath string) error {
	return nil
}
