package mem

import (
	"crypto/sha256"
	"fmt"
	"github.com/highlyunavailable/taggedimages/model"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type MemDataStore struct {
	BasePath string
	fileList map[string]*model.Image
	lock     sync.RWMutex
}

func NewMemDataStore(basePath string) *MemDataStore {
	return &MemDataStore{
		BasePath: basePath,
		fileList: make(map[string]*model.Image),
	}
}

func (ds *MemDataStore) PutImage(r io.Reader, tags []string) (*model.Image, error) {
	fm := &model.Image{}

	h := sha256.New()

	f, err := ioutil.TempFile("", "image")
	defer f.Close()

	if err != nil {
		defer os.Remove(f.Name())
		return nil, err
	}

	w := io.MultiWriter(h, f)

	if size, err := io.Copy(w, r); err == nil {
		fm.Id = fmt.Sprintf("%x", h.Sum(nil))
		fm.Size = uint64(size)
		f.Close()
	} else {
		defer os.Remove(f.Name())
		return nil, err
	}

	tempFile, err := os.Open(f.Name())
	defer tempFile.Close()
	if err == nil {
		if img, extension, err := image.Decode(tempFile); err == nil {
			fm.Width = img.Bounds().Max.X
			fm.Height = img.Bounds().Max.Y
			fm.Extension = extension
		} else {
			return nil, err
		}
		tempFile.Close()
	} else {
		return nil, err
	}

	ds.lock.Lock()
	defer ds.lock.Unlock()

	fm.Path = fmt.Sprintf("%s.%s", filepath.Join(ds.BasePath, fm.Id), fm.Extension)

	if err := os.Rename(f.Name(), fm.Path); err != nil {
		defer os.Remove(f.Name())
		if !os.IsExist(err) {
			return nil, err
		}
	}

	// Clean up and lower case all tags
	fm.Tags = tags
	for i, str := range fm.Tags {
		fm.Tags[i] = strings.ToLower(str)
	}
	sort.Strings(fm.Tags)

	ds.fileList[fm.Id] = fm

	return fm, nil
}

func (ds *MemDataStore) GetImage(id string) (*model.Image, error) {
	ds.lock.RLock()
	defer ds.lock.RUnlock()

	if image, ok := ds.fileList[id]; ok {
		return image, nil
	}
	return nil, model.ErrImageNotFound
}

func (ds *MemDataStore) GetImages(tags []string, page, pageSize int) []*model.Image {
	images := make([]*model.Image, 0)
	skip := (page - 1) * pageSize

imageLoop:
	for _, img := range ds.fileList {
		for _, tag := range tags {
			cleanTag := strings.ToLower(strings.TrimPrefix(tag, "-"))
			pos := sort.SearchStrings(img.Tags, cleanTag)
			if strings.HasPrefix(tag, "-") {
				if pos < len(img.Tags) && img.Tags[pos] == cleanTag {
					continue imageLoop
				}
			} else {
				if pos == len(img.Tags) || img.Tags[pos] != cleanTag {
					continue imageLoop
				}
			}
		}
		if skip > 0 {
			skip--
			continue imageLoop
		}
		images = append(images, img)
		if int(len(images)) == pageSize {
			break
		}
	}
	return images
}

func (ds *MemDataStore) DeleteImage(id string) error {
	if img, err := ds.GetImage(id); err == nil {
		ds.lock.Lock()
		defer ds.lock.Unlock()
		if rmerr := os.Remove(img.Path); rmerr != nil && !os.IsNotExist(rmerr) {
			return rmerr
		}
		delete(ds.fileList, id)
	} else {
		return err
	}
	return nil
}
