package handlers

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/atotto/clipboard"

	"github.com/luca2040/clipse/config"
	"github.com/luca2040/clipse/shell"
	"github.com/luca2040/clipse/utils"
)

/*
runListener is essentially a while loop to be created as a system background process on boot.
	can be stopped at any time with:
		clipse -kill
		pkill -f clipse
		killall clipse
*/

var prevClipboardContent string // used to store clipboard content to avoid re-checking media data unnecessarily
var dataType string             // used to determine which poll interval to use based on current clipboard data format

func RunListener(displayServer string, imgEnabled bool) error {
	// Listen for SIGINT (Ctrl+C) and SIGTERM signals to properly close the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// channel to pass clipboard events to
	clipboardData := make(chan string, 1)

	// Goroutine to monitor clipboard
	go func() {
		for {
			input, err := clipboard.ReadAll()
			if err != nil {
				time.Sleep(1 * time.Second) // wait for boot
			}
			if input != prevClipboardContent {
				clipboardData <- input       // Pass clipboard data to main goroutine
				prevClipboardContent = input // update previous content
			}

			if runtime.GOOS == "windows" {
				hasImage, imageStr := utils.ClipboardHasImage()

				if hasImage && utils.HasClipboardContentChanged() {
					clipboardData <- imageStr
					prevClipboardContent = imageStr
				}
			}

			if dataType == Text {
				time.Sleep(defaultPollInterval)
				continue
			}

			time.Sleep(mediaPollInterval)
		}
	}()

MainLoop:
	for {
		select {
		case input := <-clipboardData:
			if input == "" {
				continue
			}
			dataType = utils.DataType(input)
			switch dataType {
			case Text:
				if err := config.AddClipboardItem(input, "null"); err != nil {
					utils.LogERROR(fmt.Sprintf("failed to add new item `( %s )` | %s", input, err))
				}
			case PNG, JPEG:
				if imgEnabled {
					fileName := fmt.Sprintf("%s-%s.%s", strconv.Itoa(len(input)), utils.GetTimeStamp(), dataType)
					itemTitle := fmt.Sprintf("%s %s", imgIcon, fileName)
					filePath := filepath.Join(config.ClipseConfig.TempDirPath, fileName)

					var saveErr error
					if runtime.GOOS == "windows" {
						saveErr = utils.SavePNGStringToFile(input, filePath)
					} else {
						saveErr = shell.SaveImage(filePath, displayServer)
					}

					if saveErr != nil {
						utils.LogERROR(fmt.Sprintf("failed to save image [%s] | %v", runtime.GOOS, saveErr))
						break
					}
					if err := config.AddClipboardItem(itemTitle, filePath); err != nil {
						utils.LogERROR(fmt.Sprintf("failed to add image to clipboard items | %s", err))
					}
				}
			}
		case <-interrupt:
			break MainLoop
		}
	}

	return nil
}
