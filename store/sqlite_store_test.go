package store

import (
	"database/sql"
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/iamunni/hugnin/model"
	_ "github.com/mattn/go-sqlite3"
)

func newMockStore(t *testing.T) *mockStore {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return &mockStore{
		dbConn: db,
		mock:   mock,
	}
}

type mockStore struct {
	dbConn *sql.DB
	mock   sqlmock.Sqlmock
}

func TestSQLiteStore_Write(t *testing.T) {
	tests := []struct {
		name    string
		note    model.Note
		value   string
		tags    []string
		wantErr bool
	}{
		{
			name:    "sample value and tag",
			value:   "sample value",
			tags:    []string{"sample tag"},
			wantErr: false,
		},
		{
			name:    "note with value only",
			value:   "sample value",
			tags:    []string{},
			wantErr: false,
		},
		{
			name:    "note with tag only",
			tags:    []string{"sample tag"},
			wantErr: false,
		},
		{
			name:    "empty note",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStoreInstance := newMockStore(t)
			s := &SQLiteStore{
				dbConn: mockStoreInstance.dbConn,
			}
			mockStoreInstance.mock.ExpectBegin()
			mockStoreInstance.mock.ExpectPrepare("INSERT INTO notes")
			for _, tag := range tt.tags {
				mockStoreInstance.mock.ExpectExec("INSERT INTO notes").WithArgs(tt.value, tag).WillReturnResult(sqlmock.NewResult(1, 1))
			}
			mockStoreInstance.mock.ExpectCommit()
			if err := s.Write(tt.value, tt.tags); (err != nil) != tt.wantErr {
				t.Errorf("SQLiteStore.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLiteWriter_Init(t *testing.T) {
	dbFile := "sqlite-database-test.db"
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Table Creation Success",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStoreInstance := newMockStore(t)
			s := &SQLiteStore{
				dbConn: mockStoreInstance.dbConn,
			}
			mockStoreInstance.mock.ExpectPrepare("CREATE TABLE IF NOT EXISTS notes")
			mockStoreInstance.mock.ExpectExec("CREATE TABLE IF NOT EXISTS notes").WillReturnResult(sqlmock.NewResult(1, 1))
			if err := s.Init(dbFile); (err != nil) != tt.wantErr {
				t.Errorf("SQLiteWriter.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLiteStore_ReadWithoutAnyFilter(t *testing.T) {
	tests := []struct {
		name    string
		note    model.Note
		want    []model.Note
		wantErr bool
	}{
		{
			name: "read with no filter",
			note: model.Note{},
			want: []model.Note{
				{
					Id:    1,
					Value: "note1",
					Tag:   "tag1",
				},
				{
					Id:    2,
					Value: "note2",
					Tag:   "tag1,tag2",
				},
				{
					Id:    3,
					Value: "note2",
					Tag:   "tag1,tag3",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStoreInstance := newMockStore(t)
			rows := sqlmock.NewRows([]string{"id", "note", "tags"}).
				AddRow(1, "note1", "tag1").
				AddRow(2, "note2", "tag1,tag2").
				AddRow(3, "note2", "tag1,tag3")
			mockStoreInstance.mock.ExpectQuery("SELECT (.+) FROM notes;$").WillReturnRows(rows)
			s := &SQLiteStore{
				dbConn: mockStoreInstance.dbConn,
			}
			got, err := s.Read(tt.note)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLiteStore.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLiteStore.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLiteStore_ReadWithNoteFilter(t *testing.T) {
	tests := []struct {
		name    string
		note    model.Note
		want    []model.Note
		wantErr bool
	}{
		{
			name: "read with value filter",
			note: model.Note{
				Value: "note1",
			},
			want: []model.Note{
				{
					Id:    1,
					Value: "note1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStoreInstance := newMockStore(t)
			rows := sqlmock.NewRows([]string{"id", "note", "tags"}).
				AddRow(1, "note1", "")
			mockStoreInstance.mock.ExpectQuery("^SELECT (.+) FROM notes WHERE note LIKE (.+');$").WillReturnRows(rows)
			s := &SQLiteStore{
				dbConn: mockStoreInstance.dbConn,
			}
			got, err := s.Read(tt.note)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLiteStore.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLiteStore.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLiteStore_ReadWithTagFilter(t *testing.T) {
	tests := []struct {
		name    string
		note    model.Note
		want    []model.Note
		wantErr bool
	}{
		{
			name: "read with tag filter",
			note: model.Note{
				Tag: "tag1",
			},
			want: []model.Note{
				{
					Id:  1,
					Tag: "tag1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStoreInstance := newMockStore(t)
			rows := sqlmock.NewRows([]string{"id", "note", "tags"}).
				AddRow(1, "", "tag1")
			mockStoreInstance.mock.ExpectQuery("^SELECT (.+) FROM notes WHERE tags IN ((.+));").WillReturnRows(rows)
			s := &SQLiteStore{
				dbConn: mockStoreInstance.dbConn,
			}
			got, err := s.Read(tt.note)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLiteStore.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLiteStore.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLiteStore_ReadWithNoteAndTagFilter(t *testing.T) {
	tests := []struct {
		name    string
		note    model.Note
		want    []model.Note
		wantErr bool
	}{
		{
			name: "read with tag filter",
			note: model.Note{
				Value: "note1",
				Tag:   "tag1",
			},
			want: []model.Note{
				{
					Id:    1,
					Value: "note1",
					Tag:   "tag1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStoreInstance := newMockStore(t)
			rows := sqlmock.NewRows([]string{"id", "note", "tags"}).
				AddRow(1, "note1", "tag1")
			mockStoreInstance.mock.ExpectQuery("^SELECT (.+) FROM notes WHERE note LIKE (.+) AND tags IN ((.+));").WillReturnRows(rows)
			s := &SQLiteStore{
				dbConn: mockStoreInstance.dbConn,
			}
			got, err := s.Read(tt.note)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLiteStore.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLiteStore.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLiteStore_DeleteAll(t *testing.T) {
	tests := []struct {
		name    string
		note    model.Note
		wantErr bool
	}{
		{
			name: "delete all notes",
			note: model.Note{
				Id: -1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStoreInstance := newMockStore(t)
			mockStoreInstance.mock.ExpectPrepare("DELETE FROM notes")
			mockStoreInstance.mock.ExpectExec("DELETE FROM notes")
			s := &SQLiteStore{
				dbConn: mockStoreInstance.dbConn,
			}
			if err := s.Delete(tt.note); (err != nil) != tt.wantErr {
				t.Errorf("SQLiteStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLiteStore_DeleteBasedOnTag(t *testing.T) {
	tests := []struct {
		name    string
		note    model.Note
		wantErr bool
	}{
		{
			name: "delete notes based on Tag",
			note: model.Note{
				Tag: "tag1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStoreInstance := newMockStore(t)
			mockStoreInstance.mock.ExpectPrepare("DELETE FROM notes WHERE tags LIKE (.+')$")
			mockStoreInstance.mock.ExpectExec("DELETE FROM notes WHERE tags LIKE (.+')$")
			s := &SQLiteStore{
				dbConn: mockStoreInstance.dbConn,
			}
			if err := s.Delete(tt.note); (err != nil) != tt.wantErr {
				t.Errorf("SQLiteStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLiteStore_DeleteBasedOnId(t *testing.T) {
	tests := []struct {
		name    string
		note    model.Note
		wantErr bool
	}{
		{
			name: "delete notes based on Tag",
			note: model.Note{
				Id: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStoreInstance := newMockStore(t)
			mockStoreInstance.mock.ExpectPrepare("DELETE FROM notes WHERE Id=(.+)$")
			mockStoreInstance.mock.ExpectExec("DELETE FROM notes WHERE Id=(.+)$")
			s := &SQLiteStore{
				dbConn: mockStoreInstance.dbConn,
			}
			if err := s.Delete(tt.note); (err != nil) != tt.wantErr {
				t.Errorf("SQLiteStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLiteStore_SearchAvailableNoteValue(t *testing.T) {
	tests := []struct {
		name    string
		keyword string
		want    []model.Note
		wantErr bool
	}{
		{
			name:    "keyword in note value",
			keyword: "test",
			want: []model.Note{
				{
					Id:    1,
					Value: "testnote",
					Tag:   "tag1",
				},
			},
		},
		{
			name:    "keyword in tag value",
			keyword: "2",
			want: []model.Note{
				{
					Id:    2,
					Value: "note2",
					Tag:   "tag2",
				},
			},
		},
		{
			name:    "keyword in multiple tag values",
			keyword: "tag",
			want: []model.Note{
				{
					Id:    1,
					Value: "testnote",
					Tag:   "tag1",
				},
				{
					Id:    2,
					Value: "note2",
					Tag:   "tag2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStoreInstance := newMockStore(t)
			rows := mockStoreInstance.mock.NewRows([]string{"id", "note", "tags"}).
				AddRow(1, "testnote", "tag1").
				AddRow(2, "note2", "tag2")
			mockStoreInstance.mock.ExpectQuery("SELECT (\\*+) FROM notes").WillReturnRows(rows)
			s := &SQLiteStore{
				dbConn: mockStoreInstance.dbConn,
			}
			got, err := s.Search(tt.keyword)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLiteStore.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLiteStore.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
