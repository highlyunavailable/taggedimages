package taggedimages

import (
	"github.com/highlyunavailable/taggedimages/model"
	"io"
)

type Datastore interface {
	PutImage(r io.Reader, tags []string) (*model.Image, error)
	GetImage(id string) (*model.Image, error)
	GetImages(tags []string, page, pageSize int) []*model.Image
	DeleteImage(id string) error
	// TODO: Write a method that scans all files on disk and deletes ones that
	// have no metadata entry
	//CollectGarbage() error
}
