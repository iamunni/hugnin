package writer

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteWriter struct {
	dbConn *sql.DB
}

func NewSQLiteWriter() Writer {
	db, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return &SQLiteWriter{
		dbConn: db,
	}
}

func (s *SQLiteWriter) Write(value string, tags []string) error {
	defer s.dbConn.Close()
	err := insertNote(s.dbConn, value, tags)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteWriter) Init(dbFile string) error {
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
