package backend

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func fileService() {
	for {
		if !checkIfLatest() {
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
			filePath = fmt.Sprintf("./files/chalaoshi-%s.zip", time.Now().Format("20060102"))
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

func checkIfLatest() bool {
	// check if files dir exists; zip file exists; timestamp is today's
	if _, err := os.Stat("./files"); os.IsNotExist(err) {
		return false
	}
	// list files in the directory
	dirEntries, err := os.ReadDir("./files")
	if err != nil {
		return false
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), "chalaoshi-") && strings.HasSuffix(entry.Name(), ".zip") {
			// extract the timestamp from the file name
			timeStr := strings.TrimSuffix(strings.TrimPrefix(entry.Name(), "chalaoshi-"), ".zip")
			timeFile, _ := time.Parse("20060102", timeStr)

			nowYear, nowMonth, nowDay := time.Now().Date()
			fileYear, fileMonth, fileDay := timeFile.Date()
			if nowYear == fileYear && nowMonth == fileMonth && nowDay == fileDay {
				status = 0
				return true
			}
		}
	}

	return false

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
