package service

import (
	"fmt"
	"testing"

	"github.com/iamunni/hugnin/writer"
)

func (m *mockWriter) Write(value string, tag []string) error {
	if len(value) == 0 || len(tag) == 0 {
		return fmt.Errorf("%s", "error")
	}
	return nil
}

func (m *mockWriter) Init(dbFile string) error {
	return nil
}

func newMockWriter() writer.Writer {
	return &mockWriter{}
}

type mockWriter struct{}

var mockWriterInstance = newMockWriter()

func Test_noteService_Add(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		tag     string
		writer  writer.Writer
		wantErr bool
	}{
		{
			name:    "Empty value and tag",
			value:   "",
			tag:     "",
			writer:  mockWriterInstance,
			wantErr: true,
		},
		{
			name:    "Empty value and non empty tag",
			value:   "",
			tag:     "test tag",
			writer:  mockWriterInstance,
			wantErr: true,
		},
		{
			name:    "non empty value and non empty tag",
			value:   "test value",
			tag:     "test tag",
			writer:  mockWriterInstance,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &noteService{
				writer: tt.writer,
			}
			if err := n.Add(tt.value, tt.tag); (err != nil) != tt.wantErr {
				t.Errorf("noteService.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
