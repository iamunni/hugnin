package service

import (
	"fmt"
	"testing"

	"github.com/iamunni/hugnin/model"
	"github.com/iamunni/hugnin/store"
)

func (m *mockStore) Write(value string, tags []string) error {
	if len(value) == 0 {
		return fmt.Errorf("%s", "note value not passed error")
	}
	return nil
}

func (m *mockStore) Init(dbFile string) error {
	return nil
}

func (m *mockStore) Read(note model.Note) ([]model.Note, error) {
	return nil, nil
}

func (m *mockStore) Delete(note model.Note) error {
	return nil
}

func newMockStore() store.Store {
	return &mockStore{}
}

type mockStore struct{}

var mockStoreInstance = newMockStore()

func Test_noteService_Add(t *testing.T) {
	tests := []struct {
		name    string
		note    model.Note
		store   store.Store
		wantErr bool
	}{
		{
			name: "Empty value and tag",
			note: model.Note{
				Value: "",
				Tag:   "",
			},
			store:   mockStoreInstance,
			wantErr: true,
		},
		{
			name: "Empty value and non empty tag",
			note: model.Note{
				Value: "",
				Tag:   "sample tag",
			},
			store:   mockStoreInstance,
			wantErr: true,
		},
		{
			name: "non empty value and non empty tag",
			note: model.Note{
				Value: "sample value",
				Tag:   "sample tag",
			},
			store:   mockStoreInstance,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &noteService{
				store: tt.store,
			}
			if err := n.Add(tt.note); (err != nil) != tt.wantErr {
				t.Errorf("noteService.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_noteService_View(t *testing.T) {
	tests := []struct {
		name    string
		note    model.Note
		store   store.Store
		wantErr bool
	}{
		{
			name:    "View All Notes",
			note:    model.Note{},
			store:   mockStoreInstance,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &noteService{
				store: tt.store,
			}
			if err := n.View(tt.note); (err != nil) != tt.wantErr {
				t.Errorf("noteService.View() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_noteService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		note    model.Note
		store   store.Store
		wantErr bool
	}{
		{
			name:    "Delete All Notes",
			note:    model.Note{},
			store:   mockStoreInstance,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &noteService{
				store: tt.store,
			}
			if err := n.Delete(tt.note); (err != nil) != tt.wantErr {
				t.Errorf("noteService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
