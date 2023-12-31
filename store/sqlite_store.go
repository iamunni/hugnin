package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/iamunni/hugnin/model"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	dbConn *sql.DB
}

func NewSQLiteStore() Store {
	db, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return &SQLiteStore{
		dbConn: db,
	}
}

func (s *SQLiteStore) Write(value string, tags []string) error {
	defer s.dbConn.Close()
	err := insertNote(s.dbConn, value, tags)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) Read(note model.Note) ([]model.Note, error) {
	defer s.dbConn.Close()

	var sb strings.Builder
	sb.WriteString("SELECT * FROM notes")
	if !reflect.DeepEqual(note, model.Note{}) {
		sb.WriteString(" WHERE")
		if len(note.Value) > 0 {
			sb.WriteString(fmt.Sprintf(" note LIKE '%s'", note.Value))
		}
		if len(note.Tag) > 0 {
			if len(note.Value) > 0 {
				sb.WriteString(" AND")
			}
			tags := strings.Split(note.Tag, ",")
			quotedTags := "'" + strings.Join(tags, "','") + "'"
			sb.WriteString(fmt.Sprintf(" tags IN (%s)", quotedTags))
		}
	}

	sb.WriteString(";")
	stmt := sb.String()

	var result []model.Note
	rows, err := s.dbConn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var note model.Note
		err = rows.Scan(&note.Id, &note.Value, &note.Tag)
		if err != nil {
			return nil, err
		}
		result = append(result, note)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SQLiteStore) Init(dbFile string) error {
	defer s.dbConn.Close()
	err := createDatabase(dbFile)
	if err != nil {
		return err
	}
	err = createTable(s.dbConn)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) Search(keyword string) ([]model.Note, error) {
	defer s.dbConn.Close()
	rows, err := s.dbConn.Query("SELECT * FROM notes")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []model.Note
	note := model.Note{}

	for rows.Next() {
		note = model.Note{}

		err = rows.Scan(&note.Id, &note.Value, &note.Tag)

		if strings.Contains(note.Value, keyword) {
			result = append(result, note)
		} else if strings.Contains(note.Tag, keyword) {
			result = append(result, note)
		}

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *SQLiteStore) Delete(note model.Note) error {
	defer s.dbConn.Close()
	if note.Id == -1 {
		statement, err := s.dbConn.Prepare("DELETE FROM notes") // Prepare SQL Statement
		if err != nil {
			return err
		}
		statement.Exec() // Execute SQL Statements
		return nil
	} else if note.Id != 0 {
		stmt := fmt.Sprintf("DELETE FROM notes WHERE Id=%d", note.Id)
		statement, err := s.dbConn.Prepare(stmt) // Prepare SQL Statement
		if err != nil {
			return err
		}
		statement.Exec() // Execute SQL Statements
		return nil
	}
	if note.Tag != "" {
		fmt.Println("removing notes with tags")
		stmt := fmt.Sprintf("DELETE FROM notes WHERE tags LIKE '%s'", note.Tag)
		statement, err := s.dbConn.Prepare(stmt) // Prepare SQL Statement
		if err != nil {
			return err
		}
		statement.Exec() // Execute SQL Statements
		return nil
	}
	return nil
}

func createDatabase(dbFile string) error {
	file, err := os.Create(dbFile) // Create SQLite file
	if err != nil {
		return err
	}
	file.Close()
	log.Println("sqlite-database.db created")
	return nil
}

func createTable(dbConn *sql.DB) error {
	createNotesTableSQL := `CREATE TABLE IF NOT EXISTS notes
		(id INTEGER PRIMARY KEY AUTOINCREMENT,
		note TEXT,
		tags TEXT);` // SQL Statement for Create Table

	log.Println("Create Notes table...")
	statement, err := dbConn.Prepare(createNotesTableSQL) // Prepare SQL Statement
	if err != nil {
		return err
	}
	statement.Exec() // Execute SQL Statements
	log.Println("notes table created")
	return nil
}

func insertNote(dbConn *sql.DB, value string, tags []string) error {
	log.Println("Inserting notes record ...")

	tx, err := dbConn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO notes (note, tags) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, tag := range tags {
		_, err = stmt.Exec(value, tag)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
