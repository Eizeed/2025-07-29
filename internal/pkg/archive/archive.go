package archive

import (
	"errors"

	"github.com/Eizeed/2025-07-29/internal/pkg/constants"
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

func (archive *Archive) AddPath(path string) error {
	if len(archive.Content) >= constants.URL_LIMIT {
		return errors.New(ErrArchiveFull)
	}

	archive.Content = append(archive.Content, path)
	return nil
}
