package backend

import (
	"fmt"
	"github.com/rapid7/go-get-proxied/proxy"
	"io"
	"net/http"
	"os" // 用于创建临时目录和文件
	"path/filepath"
	"time"
)

func getSystemProxyClient() http.Client {
	systemProxy := proxy.NewProvider("").GetProxy("https", "https://chalaoshi.click")
	proxyTransport := &http.Transport{}
	if systemProxy != nil && systemProxy.String() != "" {
		proxyTransport.Proxy = http.ProxyURL(systemProxy.URL())
	}
	client := http.Client{Transport: proxyTransport}
	return client
}

func testIfConnected() bool {
	client := getSystemProxyClient()
	client.Timeout = 3 * time.Second // 设置超时时间为5秒
	resp, err := client.Get("https://chalaoshi.click")

	if err == nil && resp.StatusCode == http.StatusOK {
		connected = true
		return true
	}

	connected = false
	return false
}

func downloadFromSource(zipURL string) (string, error) {
	fmt.Printf("Downloading zip file from: %s\n", zipURL)

	client := getSystemProxyClient()
	resp, err := client.Get(zipURL)
	if err != nil {
		return "", fmt.Errorf("failed to download zip file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return "", fmt.Errorf("failed to download zip file: received status code %d", resp.StatusCode)
	}

	zipFile, err := createZipFile()
	if err != nil {
		return "", fmt.Errorf("failed to create zip file")
	}

	_, err = io.Copy(zipFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write downloaded content to temporary file: %w", err)
	}

	// 4. 关闭 zip 文件以便重新打开并读取
	err = zipFile.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close temporary zip file: %w", err)
	}

	return zipFile.Name(), nil
}

func GetTempDir() (string, error) {
	tempDir := os.TempDir()
	tempDirPath := filepath.Join(tempDir, "chalaoshi")

	absPath, err := filepath.Abs(tempDirPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %s: %w", tempDirPath, err)
	}

	if _, err := os.Stat(tempDirPath); err == nil {
		return absPath, nil
	}

	err = os.MkdirAll(tempDirPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", tempDirPath, err)
	}

	return absPath, nil
}

func createZipFile() (*os.File, error) {
	filePath, err := GetTempDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get temporary directory: %w", err)
	}
	// 如果目录存在，删除里面的文件
	dirEntries, err := os.ReadDir(filePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read 'files' directory: %w", err)
	}
	for _, entry := range dirEntries {
		err := os.RemoveAll(filepath.Join(filePath, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to remove file or directory '%s': %w", entry.Name(), err)
		}
	}

	// 1. 创建文件来存储下载的 zip 内容
	fileName := fmt.Sprintf("chalaoshi-%s.zip", time.Now().Format("20060102"))
	fullFilePath := filepath.Join(filePath, fileName)
	zipFile, err := os.Create(fullFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary zip file: %w", err)
	}

	return zipFile, nil
}
