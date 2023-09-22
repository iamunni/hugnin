package service

import (
	"fmt"
	"strings"

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
	var tags []string
	for _, tag := range strings.Split(tag, ",") {
		tags = append(tags, strings.TrimSpace(tag))
	}
	err := n.writer.Write(value, tags)
	if err != nil {
		return err
	}
	return nil
}
