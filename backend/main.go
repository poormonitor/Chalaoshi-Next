package backend

import (
	"encoding/base64"
)

var (
	filePath  = ""
	status    = 1
	connected = false
	url       = "https://chalaoshi.click/static/download/chalaoshi_csv20250502_5399305_2696_26893D_sha256.zip"
)

func Init() {
	filePath = ""
	go fileService()

	addrLocal := getAddrLocal()
	if len(addrLocal) > 0 {
		go announcePeer(addrLocal)
		runServer(addrLocal[0].Port)
	}
}

func GetFileBase64() string {
	return base64.StdEncoding.EncodeToString(getFileBytes(filePath))
}

func GetStatus() int {
	return status
}
