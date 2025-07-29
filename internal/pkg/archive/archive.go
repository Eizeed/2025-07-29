package archive

import (
	"errors"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

const (
	ErrArchiveFull = "ErrArchiveFull"
)

type ArchiveRepo struct {
	inner  [3]Archive
	mutex  sync.Mutex
}

type Archive struct {
	UUID    uuid.UUID
	Content []string
}

func New() Archive {
	return Archive{
		UUID:    uuid.New(),
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
			log.Panicln("Unable to read a file in path: ", entry)
			continue
		}

		buffers = append(buffers, bytes)
	}

	return nil
}

func (archive *Archive) AddPath(path string) error {
	if len(archive.Content) > 3 {
		log.Fatalln("Assertion failed. len(Archive.content) expected < 3, but got", len(archive.Content))
	}

	if len(archive.Content) >= 3 {
		return errors.New(ErrArchiveFull)
	}

	archive.Content = append(archive.Content, path)
	return nil
}
