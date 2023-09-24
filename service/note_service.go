package service

import (
	"fmt"
	"strings"

	"github.com/iamunni/hugnin/model"
	"github.com/iamunni/hugnin/store"
)

type NoteService interface {
	Add(note model.Note) error
	View(note model.Note) error
}

type noteService struct {
	store store.Store
}

func NewNoteService(store store.Store) NoteService {
	return &noteService{
		store: store,
	}
}

func (n *noteService) Add(note model.Note) error {
	if len(note.Value) == 0 {
		return fmt.Errorf("%s", "note value not passed error")
	}
	var tags []string
	for _, tag := range strings.Split(note.Tag, ",") {
		tags = append(tags, strings.TrimSpace(tag))
	}
	err := n.store.Write(note.Value, tags)
	if err != nil {
		return err
	}
	return nil
}

func (n *noteService) View(note model.Note) error {
	result, err := n.store.Read(note)
	if err != nil {
		return err
	}
	print(result)
	return nil
}

func print(notes []model.Note) {
	fmt.Printf("%v", notes)
}
