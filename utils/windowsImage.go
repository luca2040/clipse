package utils

import (
	"bytes"
	"errors"
	"image/png"
	"os"
	"syscall"

	"github.com/gonutz/w32/v2"
	"golang.design/x/clipboard"
)

var (
	user32                         = syscall.NewLazyDLL("user32.dll")
	procGetClipboardSequenceNumber = user32.NewProc("GetClipboardSequenceNumber")
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procGlobalSize                 = kernel32.NewProc("GlobalSize")
)

func GetClipboardSequenceNumber() uint32 {
	ret, _, _ := procGetClipboardSequenceNumber.Call()
	return uint32(ret)
}

func GlobalSize(h uintptr) uint32 {
	ret, _, _ := procGlobalSize.Call(h)
	return uint32(ret)
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
				LogERROR("Clipboard image is empty or unavailable")
				return false, ""
			}

			// Decode PNG bytes into an image.Image
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil {
				LogERROR("Failed to decode PNG from clipboard data")
				return false, ""
			}

			// Encode back to PNG in memory buffer
			var pngBuf bytes.Buffer
			if err := png.Encode(&pngBuf, img); err != nil {
				LogERROR("Failed to encode PNG to buffer")
				return false, ""
			}

			return true, pngBuf.String()
		}
	}

	return false, ""
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

func SavePNGStringToFile(pngStr string, filePath string) error {
	// Convert string back to []byte
	data := []byte(pngStr)

	// Create or truncate the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the byte data to the file
	n, err := file.Write(data)
	if err != nil {
		return err
	}

	// Sanity check: did we write all data?
	if n != len(data) {
		return errors.New("incomplete write to file")
	}

	return nil // Success
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
