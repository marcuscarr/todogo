package todo

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Todo is a task to be done
type Todo struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Status      bool      `db:"status"`
	Created     time.Time `db:"created"`
	Modified    time.Time `db:"modified"`
}

func new(title string, desc string) Todo {
	now := time.Now()
	return Todo{
		Title:       title,
		Description: desc,
		Status:      false,
		Created:     now,
		Modified:    now,
	}
}

// Create adds a new To-Do to the database and returns its id
func Create(db *sqlx.DB, title string, desc string) (int64, error) {
	todo := new(title, desc)
	sqlStatement := `
		INSERT INTO todos (title, description, created, modified)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	var id int64 = -1
	err := db.QueryRow(sqlStatement, todo.Title, todo.Description, todo.Created, todo.Modified).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

// Retrieve returns the To-Do with the supplied ID, if it exists
func Retrieve(db *sqlx.DB, id int64) (*Todo, error) {
	var todo *Todo
	sqlStatement := `
		SELECT * FROM todos WHERE id=$1`
	err := db.Get(todo, sqlStatement, id)
	if err == nil {
		return todo, nil
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nil, err
}

// Update modifies the To-Do that matches the supplied ID
func Update(ctx context.Context, db *sqlx.DB, id int64, title *string, desc *string, status *bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	updated := false
	if title != nil {
		err = updateTitle(tx, id, *title)
		if err != nil {
			return err
		}
		updated = true
	}
	if desc != nil {
		err = updateDesc(tx, id, *desc)
		if err != nil {
			return err
		}
		updated = true
	}
	if status != nil {
		err = updateStatus(tx, id, *status)
		if err != nil {
			return err
		}
		updated = true
	}
	if updated {
		err = updateModified(tx, id, time.Now())
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

func updateTitle(tx *sql.Tx, id int64, title string) error {
	sqlStatement := `
	UPDATE todos
	SET title = $2
	WHERE id = $1;`
	_, err := tx.Exec(sqlStatement, id, title)
	return err
}

func updateDesc(tx *sql.Tx, id int64, desc string) error {
	sqlStatement := `
	UPDATE todos
	SET desc = $2
	WHERE id = $1;`
	_, err := tx.Exec(sqlStatement, id, desc)
	return err
}

func updateStatus(tx *sql.Tx, id int64, status bool) error {
	sqlStatement := `
	UPDATE todos
	SET status = $2
	WHERE id = $1;`
	_, err := tx.Exec(sqlStatement, id, status)
	return err
}

func updateModified(tx *sql.Tx, id int64, mod time.Time) error {
	sqlStatement := `
	UPDATE todos
	SET modified = $2
	WHERE id = $1;`
	_, err := tx.Exec(sqlStatement, id, mod)
	return err
}

// Delete removes the To-Do that matches the supplied ID
func Delete(db *sqlx.DB, id int64) (int64, error) {
	sqlStatement := `
	DELETE FROM todos
	WHERE id = $1;`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	return count, err
}
