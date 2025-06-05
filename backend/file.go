package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func fileService() {
	for {
		available, filePathTemp := checkIfLatest()
		if !available {
			status = 1
			if testIfConnected() {
				fmt.Println("Internet connection detected. Preparing file.")
				err := fetchFileFromSource(url)
				if err != nil {
					time.Sleep(3 * time.Second)
					continue
				}
			} else {
				fmt.Println("No internet connection. Skipping file download.")
				err := downloadFileFromPeer()
				if err != nil {
					time.Sleep(3 * time.Second)
					continue
				}
			}
		} else {
			filePath = filePathTemp
			time.Sleep(10 * time.Minute) // 每10分钟检查一次
		}
	}
}

func fetchFileFromSource(url string) error {
	var err error
	status = 2
	filePath, err = downloadFromSource(url)
	if err != nil {
		return fmt.Errorf("failed to download file from source: %w", err)
	}
	status = 0
	return nil
}

func checkIfLatest() (bool, string) {
	// check if files dir exists; zip file exists; timestamp is today's
	filePath, err := GetTempDir()
	if err != nil {
		fmt.Println("Failed to get temporary directory:", err)
		return false, ""
	}

	// list files in the directory
	dirEntries, err := os.ReadDir(filePath)
	if err != nil {
		return false, ""
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		if strings.HasPrefix(entry.Name(), "chalaoshi-") && strings.HasSuffix(entry.Name(), ".zip") {
			// extract the timestamp from the file name
			timeStr := strings.TrimSuffix(strings.TrimPrefix(entry.Name(), "chalaoshi-"), ".zip")
			timeFile, err := time.Parse("20060102", timeStr)

			if err != nil {
				fmt.Println("Failed to parse time from file name:", err)
				continue
			}

			nowYear, nowMonth, nowDay := time.Now().Date()
			fileYear, fileMonth, fileDay := timeFile.Date()
			if nowYear == fileYear && nowMonth == fileMonth && nowDay == fileDay {
				status = 0
				return true, filepath.Join(filePath, entry.Name())
			}
		}
	}

	return false, ""
}

func getFileBytes(filePath string) []byte {
	if filePath == "" {
		return nil
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}

	return fileBytes
}
