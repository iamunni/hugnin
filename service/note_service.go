package service

import (
	"fmt"

	"github.com/iamunni/hugnin/writer"
)

type NoteService interface {
	Add(value string, tag string) error
}

type noteService struct {
	writer writer.Writer
}

func NewNoteService(writer writer.Writer) NoteService {
	return &noteService{
		writer: writer,
	}
}

func (n *noteService) Add(value string, tag string) error {
	fmt.Printf("%+v\n", n)
	err := n.writer.Write(value, tag)
	if err != nil {
		return err
	}
	return nil
}
