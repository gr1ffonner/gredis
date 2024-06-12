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
	// sqlFile, err := os.Open("mock_data.sql")
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	// defer sqlFile.Close()
	//
	// sqlBytes, err := io.ReadAll(sqlFile)
	// if err != nil {
	// 	logger.Fatal(err)
	// }

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

	// _, err = db.Exec(string(sqlBytes))
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	stmt1, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS article(
		articleid INTEGER PRIMARY KEY,
		title VARCHAR(50) NOT NULL UNIQUE,
		text TEXT NOT NULL);
	`)

	stmt2, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS comment(
		commentid INTEGER PRIMARY KEY,
		text TEXT NOT NULL UNIQUE,
		rating INTEGER NOT NULL,
    articleid INTEGER NOT NULL REFERENCES article(articleid));
	`)

	_, err = stmt1.Exec()
	if err != nil {
		logger.Fatal(err)
	}

	_, err = stmt2.Exec()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("database is working and tables were created")

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

func (db *DB) GetArticleById(id string) Article {
	var a Article
	getStmt := `select * from "article" where "articleid" = $1`
	rows, e := db.conn.Query(getStmt, id)
	CheckError(e)
	defer rows.Close()
	for rows.Next() {
		e = rows.Scan(&a.ArticleID, &a.Title, &a.Text)
		CheckError(e)
	}
	return a
}
