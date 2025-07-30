package archive

import (
	"errors"
	"log"
	"os"

	"github.com/Eizeed/2025-07-29/pkg/uuid"
)

const (
	ErrArchiveFull = "ErrArchiveFull"
)

type Archive struct {
	UUID    uuid.UUID
	Content []string
}

func NewArchive() Archive {
	return Archive{
		UUID:    uuid.NewV4(),
		Content: []string{},
	}
}

func (archive *Archive) SaveToPath(path string) error {
	if len(archive.Content) > 3 {
		log.Fatalln("Assertion failed. len(Archive.content) expected < 3, but got", len(archive.Content))
	}

	buffers := [][]byte{}

	for _, entry := range archive.Content {
		bytes, err := os.ReadFile(entry)
		if err != nil {
			log.Println("Unable to read a file in path: ", entry)
			continue
		}

		buffers = append(buffers, bytes)
	}

	return nil
}

func (archive *Archive) AddPath(path string) error {
	if len(archive.Content) > 3 {
		log.Fatalln("Assertion failed. len(Archive.content) expected <= 3, but got", len(archive.Content))
	}

	if len(archive.Content) >= 3 {
		return errors.New(ErrArchiveFull)
	}

	archive.Content = append(archive.Content, path)
	return nil
}
