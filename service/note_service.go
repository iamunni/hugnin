package service

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/iamunni/hugnin/model"
	"github.com/iamunni/hugnin/store"
	"github.com/olekukonko/tablewriter"
)

type NoteService interface {
	Add(note model.Note) error
	View(note model.Note) error
	Delete(note model.Note) error
	Search(keyword string) error
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

func (n *noteService) Search(keyword string) error {
	result, err := n.store.Search(keyword)
	if err != nil {
		return err
	}
	print(result)
	return nil
}

func (n *noteService) Delete(note model.Note) error {
	err := n.store.Delete(note)
	if err != nil {
		return err
	}
	return nil
}

func print(notes []model.Note) {
	var data = [][]string{}

	for _, note := range notes {
		data = append(data, []string{strconv.Itoa(int(note.Id)), note.Value, note.Tag})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Note", "Tag"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
