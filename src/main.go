package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"cloud_backuper/internal/client/s3"
	"cloud_backuper/internal/configs"
	"cloud_backuper/internal/crypto"
	"cloud_backuper/internal/filesystem"
)

func init() {
	logFile, err := os.OpenFile("cloud_upload.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть файл логов: %v", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Println("начало загрузки бэкапов...")

	config := configs.LoadConfigs()
	
	s3Client := s3.New(&config.S3)

	files := filesystem.LS(config.LocalDirectoryPath)

	for _, file := range files {
		for _, dir := range config.DirectoryStruct {
			if strings.HasPrefix(file.Name(), dir.PrefixFile) {
				objectName := path.Join(dir.CloudDir, file.Name())
				fullPath := filepath.Join(config.LocalDirectoryPath, file.Name())
				hash := crypto.ComputeFileMD5(fullPath)
				
				if !s3Client.FileExists(hash, config.S3.Backet, objectName) {
					s3Client.UploadFile(config.S3.Backet, objectName, fullPath)
				}
				
				break
			}
		}

	}

	log.Println("конец загрузки бэкапов!")
}
