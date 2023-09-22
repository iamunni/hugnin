package writer

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
)

func newMockWriter(t *testing.T) *mockWriter {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return &mockWriter{
		dbConn: db,
		mock:   mock,
	}
}

type mockWriter struct {
	dbConn *sql.DB
	mock   sqlmock.Sqlmock
}

func TestSQLiteWriter_Write(t *testing.T) {
	tests := []struct {
		name    string
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWriterInstance := newMockWriter(t)
			s := &SQLiteWriter{
				dbConn: mockWriterInstance.dbConn,
			}
			mockWriterInstance.mock.ExpectBegin()
			mockWriterInstance.mock.ExpectPrepare("INSERT INTO notes")
			for _, tag := range tt.tags {
				mockWriterInstance.mock.ExpectExec("INSERT INTO notes").WithArgs(tt.value, tag).WillReturnResult(sqlmock.NewResult(1, 1))
			}
			mockWriterInstance.mock.ExpectCommit()
			if err := s.Write(tt.value, tt.tags); (err != nil) != tt.wantErr {
				t.Errorf("SQLiteWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
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
			mockWriterInstance := newMockWriter(t)
			s := &SQLiteWriter{
				dbConn: mockWriterInstance.dbConn,
			}
			mockWriterInstance.mock.ExpectPrepare("CREATE TABLE IF NOT EXISTS notes")
			mockWriterInstance.mock.ExpectExec("CREATE TABLE IF NOT EXISTS notes").WillReturnResult(sqlmock.NewResult(1, 1))
			if err := s.Init(dbFile); (err != nil) != tt.wantErr {
				t.Errorf("SQLiteWriter.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
