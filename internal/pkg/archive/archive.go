package archive

import (
	"errors"
	"log"
	"os"
	"sync"

	"github.com/Eizeed/2025-07-29/pkg/uuid"
)

const (
	ErrArchiveFull      = "ErrArchiveFull"
	ErrArchiveQueueFull = "ErrArchiveQueueFull"
)

type ArchiveRepo struct {
	Inner []Archive
	mutex sync.Mutex
}

func NewArchiveRepo() ArchiveRepo {
	return ArchiveRepo{
		Inner: []Archive{},
		mutex: sync.Mutex{},
	}
}

func (ar *ArchiveRepo) Push(archive Archive) error {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()

	if len(ar.Inner) >= 3 {
		return errors.New(ErrArchiveQueueFull)
	}

	ar.Inner = append(ar.Inner, archive)

	return nil
}

func (ar *ArchiveRepo) RemoveByUUID(uuid uuid.UUID) (Archive, bool) {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()

	for i, archive := range ar.Inner {
		if archive.UUID == uuid {
			ar.Inner = append(ar.Inner[:i], ar.Inner[i+1:]...)
			return archive, true
		}
	}

	return Archive{}, false
}

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
