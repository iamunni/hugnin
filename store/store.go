package store

import (
	"github.com/iamunni/hugnin/model"
)

type Store interface {
	Init(string) error
	Write(value string, tags []string) error
	Read(note model.Note) ([]model.Note, error)
	Delete(note model.Note) error
	Search(keyword string) ([]model.Note, error)
}
