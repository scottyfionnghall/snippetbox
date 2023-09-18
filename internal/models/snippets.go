package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type to hold the data for an individual snippet.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type wich wraps a sql.DB connection pool.
type SnippetModel struct {
	DB         *sql.DB
	InserStmt  *sql.Stmt
	GetStmt    *sql.Stmt
	LatestStmt *sql.Stmt
	DeleteStmt *sql.Stmt
}

// Creates a constructor for a SnippetModel, which includes prepared statements.
// This is need so we can reuse this statements and not recrete them on each call.
func NewSnippetModel(db *sql.DB) (*SnippetModel, error) {
	insertStmt, err := db.Prepare(`INSERT INTO snippets (title, content, created, expires)
	VALUES(?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))`)
	if err != nil {
		return nil, err
	}
	getStmt, err := db.Prepare(`SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`)
	if err != nil {
		return nil, err
	}
	latestStmt, err := db.Prepare(`SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`)
	if err != nil {
		return nil, err
	}
	deleteStmt, err := db.Prepare(`DELETE FROM snippets WHERE id=?`)
	if err != nil {
		return nil, err
	}
	return &SnippetModel{db, insertStmt, getStmt, latestStmt, deleteStmt}, nil
}

// Closes all the prepared statements to ensuare that it is properly closed
// before main function terminates
func (s *SnippetModel) CloseAll() error {
	err := s.InserStmt.Close()
	if err != nil {
		return err
	}
	s.GetStmt.Close()
	if err != nil {
		return err
	}
	s.DeleteStmt.Close()
	if err != nil {
		return err
	}
	s.LatestStmt.Close()
	if err != nil {
		return err
	}
	return nil
}

// Function to insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// Use the Exec() method on the embedded connection pool to execute
	// the statement
	result, err := m.InserStmt.Exec(title, content, expires)
	if err != nil {
		return 0, err
	}
	// Use the LatestInserID() method on the result to get the ID of our newly
	// inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// The ID returned has the type int64, so we convert it to an int type
	// before returning
	return int(id), nil
}

// Function to return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	// Use the QueryRow() method on the connection pool to execure our
	// SQL statement, passing in the untrasted id variable as the value for
	// the placeholder patameter. This returns a pointer to a sql.Row object
	// wich holds the result from the database.
	row := m.GetStmt.QueryRow(id)
	// Initialize a pointer to a new zeroed Snippet struct.
	s := &Snippet{}
	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// Function to return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	rows, err := m.LatestStmt.Query()
	if err != nil {
		return nil, err
	}
	// Defer rows.Close() to ensure the sql.Rows resultest is always
	// properly cloesd before the Latest() method returns.
	defer rows.Close()
	// Initialize an empty slice to hold the Snippet sturcts.
	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

func (m *SnippetModel) Delete(id int) error {
	_, err := m.DeleteStmt.Query(id)
	if err != nil {
		return ErrNoRecord
	} else {
		return nil
	}

}
