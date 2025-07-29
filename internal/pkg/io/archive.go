package io

import (
	"archive/zip"
	"log"
	"os"
	"path/filepath"

	"github.com/Eizeed/2025-07-29/internal/pkg/archive"
	"github.com/google/uuid"
)

func zipDirPath() string {
	path := os.Getenv("ZIP_PATH")
	if path == "" {
		return getPwd()
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return getPwd()
	}

	err = os.MkdirAll(absPath, 0755)
	if err != nil {
		return getPwd()
	}

	return absPath
}

func getPwd() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println("Error occured while trying to get PWD:", err)
		panic("Where is PWD?")
	}

	return path
}

func ZipFromArchive(archive *archive.Archive) string {
	zipDirPath := zipDirPath()
	uuidStr := archive.UUID.String()
	if uuidStr == "" {
		uuidStr = uuid.New().String()
	}
	zipPath := filepath.Join(zipDirPath, uuidStr+".zip")

	zipFile, err := os.Create(zipPath)
	if err != nil {

	}

	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		if err := zipWriter.Close(); err != nil {

		}
	}()

	for _, content := range archive.Content {
		bytesContent, err := os.ReadFile(content)
		if err != nil {
		}

		relPath, err := filepath.Rel(zipDirPath, content)
		if err != nil {
		}

		writer, err := zipWriter.Create(relPath)
		if err != nil {
		}

		writer.Write(bytesContent)
	}

	return zipPath
}
