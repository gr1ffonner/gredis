package db

import (
	"database/sql"
	"fmt"
	"gredis/internal/config"
	"gredis/pkg/logging"

	_ "github.com/lib/pq"
)

// DB is a struct to hold the database connection.
type DB struct {
	conn *sql.DB
}

// NewDB creates a new instance of DB with the provided connection string.
func NewDB(cfg config.Config, logger logging.Logger) (*DB, error) {
	psqlconn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	return &DB{conn: db}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}

// Ping pings the database to check the connection.
func (db *DB) Ping() error {
	return db.conn.Ping()
}

// CheckError checks for errors and panics if found.
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func (db *DB) CreateTable() {
	createStmt := `create table if not exists "Students" ("ID" serial primary key, "Name" varchar(255), "Roll" integer)`
	_, e := db.conn.Exec(createStmt)
	CheckError(e)
}

func (db *DB) InsertFake() {
	insertStmt := `insert into "Students"("Name", "Roll") values('John', 1)`
	_, e := db.conn.Exec(insertStmt)
	CheckError(e)
}

func (db *DB) GetFake() {
	getStmt := `select * from "Students"`
	rows, e := db.conn.Query(getStmt)
	CheckError(e)
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var roll int
		e = rows.Scan(&id, &name, &roll)
		CheckError(e)
		fmt.Println(id, name, roll)
	}
}
