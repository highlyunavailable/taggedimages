package model

import "errors"

var ErrImageNotFound = errors.New("the image ID was not found in the store")

type Image struct {
	Id            string
	Tags          []string
	Path          string
	Extension     string
	Height, Width int
	Size          uint64
}
