package taggedimages

import (
	"github.com/highlyunavailable/taggedimages/model"
	"io"
)

type ImageStore struct {
	Datastore Datastore
}

func New(ds Datastore) *ImageStore {
	return &ImageStore{Datastore: ds}
}

func (b *ImageStore) PutImage(r io.Reader, tags []string) (img *model.Image, err error) {
	if img, err := b.Datastore.PutImage(r, tags); err != nil {
		return img, err
	}
	return
}

func (b *ImageStore) GetImage(id string) (img *model.Image, err error) {
	if img, err = b.Datastore.GetImage(id); err != nil {
		return nil, err
	}
	return
}

func (b *ImageStore) GetImages(tags []string, page, pageSize int) []*model.Image {
	return b.Datastore.GetImages(tags, page, pageSize)
}

func (b *ImageStore) DeleteImage(id string) error {
	if err := b.Datastore.DeleteImage(id); err != nil {
		return err
	}
	return nil
}
