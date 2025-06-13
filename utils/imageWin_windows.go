//go:build windows
// +build windows

package utils

import (
	"errors"
	"os"
	"syscall"

	"github.com/gonutz/w32/v2"
	"golang.design/x/clipboard"
)

var (
	user32                         = syscall.NewLazyDLL("user32.dll")
	procGetClipboardSequenceNumber = user32.NewProc("GetClipboardSequenceNumber")
)

func GetClipboardSequenceNumber() uint32 {
	ret, _, _ := procGetClipboardSequenceNumber.Call()
	return uint32(ret)
}

var prevClipboardNum uint32 = 0

func HasClipboardContentChanged() bool {
	current := GetClipboardSequenceNumber()
	if current != prevClipboardNum {
		prevClipboardNum = current
		return true
	}
	return false
}

func ClipboardHasImage() (bool, string) {
	if !w32.OpenClipboard(0) {
		return false, ""
	}
	defer w32.CloseClipboard()

	for _, format := range []uint{8, 17} { // Only CF_DIB and CF_DIBV5
		h := w32.GetClipboardData(format)
		if h != 0 {
			data := clipboard.Read(clipboard.FmtImage)
			if len(data) == 0 {
				LogERROR("clipboard image is empty or unavailable")
				return false, ""
			}

			return true, string(data)
		}
	}

	return false, ""
}

func SavePNGStringToFile(pngStr string, filePath string) error {
	// Convert string back to []byte
	data := []byte(pngStr)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Write(data)
	if err != nil {
		return err
	}

	// Check if file has been written completely
	if n != len(data) {
		LogERROR("cannot write image file")
		return errors.New("incomplete write to file")
	}

	return nil
}

func CopyImageToClipboard(imagePath string) error {
	// Read the file into a []byte
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return err
	}

	// Write it to the clipboard as image data
	clipboard.Write(clipboard.FmtImage, data)

	return nil
}
