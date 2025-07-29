package io

import (
	"archive/zip"
	"log"
	"os"
	"path/filepath"

	"github.com/Eizeed/2025-07-29/internal/pkg/archive"
	"github.com/google/uuid"
)

func SaveToFileDir(name string, bytes []byte) (string, error) {
	dirPath, err := fileDirPath()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(dirPath, name)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	file.Write(bytes)

	return filePath, nil
}

func fileDirPath() (string, error) {
	path := os.Getenv("FILE_PATH")
	if path == "" {
		return defaulFileDir()
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return defaulFileDir()
	}

	err = os.MkdirAll(absPath, 0755)
	if err != nil {
		return defaulFileDir()
	}

	return absPath, nil
}

func ZipDirPath() (string, error) {
	path := os.Getenv("ZIP_PATH")
	if path == "" {
		return defaulZigDir()
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return defaulZigDir()
	}

	err = os.MkdirAll(absPath, 0755)
	if err != nil {
		return defaulZigDir()
	}

	return absPath, nil
}

func defaulFileDir() (string, error) {
	fileDirPath := filepath.Join(getPwd(), "files")

	err := os.MkdirAll(fileDirPath, 0755)
	if err != nil {
		return "", err
	}

	return fileDirPath, nil
}

func defaulZigDir() (string, error) {
	zipDirPath := filepath.Join(getPwd(), "zip")

	err := os.MkdirAll(zipDirPath, 0755)
	if err != nil {
		return "", err
	}

	return zipDirPath, nil
}

func getPwd() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println("Error occured while trying to get PWD:", err)
		panic("Where is PWD?")
	}

	return path
}

func ZipFromArchive(archive *archive.Archive) (string, error) {
	zipDirPath, err := ZipDirPath()
	if err != nil {
		return "", err
	}

	uuidStr := archive.UUID.String()
	if uuidStr == "" {
		uuidStr = uuid.New().String()
	}
	zipPath := filepath.Join(zipDirPath, uuidStr+".zip")

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}

	zipWriter := zip.NewWriter(zipFile)

	for _, content := range archive.Content {
		bytesContent, err := os.ReadFile(content)
		if err != nil {
			return "", err
		}

		base := filepath.Base(content)

		writer, err := zipWriter.Create(base)
		if err != nil {
			return "", err
		}

		writer.Write(bytesContent)
	}

	if err := zipWriter.Close(); err != nil {
		return "", err
	}

	return zipPath, nil
}
