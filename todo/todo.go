package todo

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// Todo is a task to be done
type Todo struct {
	ID          int       `db:"id"`
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
func Create(db *sqlx.DB, id *int, title *string, desc *string, status *bool) (int, error) {
	var newTitle, newDesc string
	if title != nil {
		newTitle = *title
	}
	if desc != nil {
		newDesc = *desc
	}
	todo := new(newTitle, newDesc)
	sqlFields := []string{"title", "description", "created", "modified", "status"}
	if id != nil {
		sqlFields = append(sqlFields, "id")
		todo.ID = *id
	}
	if status != nil {
		todo.Status = *status
	}
	newID := -1
	// // Need to first ensure that the ID generator will give a valid ID.
	_, err := db.Exec("SELECT setval('todo_id_seq', (SELECT MAX(id) FROM todos));")
	if err != nil {
		return newID, err
	}
	sqlStatement := fmt.Sprintf(`
		INSERT INTO todos (%s)
		VALUES (%s)
		RETURNING id;`,
		strings.Join(sqlFields, ", "),
		strings.Join(prepend(sqlFields, ":"), ", "),
	)
	res, err := db.NamedQuery(sqlStatement, todo)
	if err != nil {
		return newID, err
	}
	for res.Next() {
		err = res.Scan(&newID)
		if newID == -1 {
			return newID, fmt.Errorf("unknown error creating todo:%s", err)
		}
	}
	return newID, nil
}

// Retrieve returns the To-Do with the supplied ID, if it exists
func Retrieve(db *sqlx.DB, id int) (*Todo, error) {
	todo := Todo{}
	sqlStatement := `
		SELECT * FROM todos WHERE id=$1`
	err := db.QueryRowx(sqlStatement, id).StructScan(&todo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &todo, nil
}

// Update modifies the To-Do that matches the supplied ID
func Update(db *sqlx.DB, id int, title *string, desc *string, status *bool) (int, error) {
	sqlStatement := `
		UPDATE todos 
		SET`

	updated := false
	update := Todo{ID: id}
	if title != nil {
		sqlStatement += " title=:title,"
		update.Title = *title
		updated = true
	}
	if desc != nil {
		sqlStatement += " description=:description,"
		update.Description = *desc
		updated = true
	}
	if status != nil {
		sqlStatement += " status=:status,"
		update.Status = *status
		updated = true
	}
	if !updated {
		return 0, nil
	}
	sqlStatement += " modified=:modified"
	update.Modified = time.Now()
	sqlStatement += " WHERE id=:id"
	res, err := db.NamedExec(sqlStatement, update)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	return int(count), err
}

func replace(db *sqlx.DB, id int, titlep *string, descp *string, statusp *bool) (int, error) {
	var title, desc string
	var status bool
	if titlep != nil {
		title = *titlep
	}
	if descp != nil {
		desc = *descp
	}
	if statusp != nil {
		status = *statusp
	}
	return Update(db, id, &title, &desc, &status)
}

// Upsert will update a To-Do, if it exists, and create it otherwise. If a new To-Do is
// created, it returns the ID. Otherwise, it returns -1.
func Upsert(db *sqlx.DB, id int, title *string, desc *string, status *bool) (int, error) {
	count, err := replace(db, id, title, desc, status)
	if err != nil {
		return -1, err
	}
	if count == 1 {
		// The To-Do exists
		return id, nil
	}
	return Create(db, &id, title, desc, status)
}

// Delete removes the To-Do that matches the supplied ID
func Delete(db *sqlx.DB, id int) (int, error) {
	sqlStatement := `
	DELETE FROM todos
	WHERE id = $1;`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	return int(count), err
}

// GetAll retrieves all To-Dos
func GetAll(db *sqlx.DB) ([]Todo, error) {
	var todos []Todo
	sqlStatement := `
		SELECT * FROM todos`
	err := db.Select(&todos, sqlStatement)
	return todos, err
}

func prepend(str []string, prefix string) []string {
	res := make([]string, len(str))
	for idx, s := range str {
		res[idx] = prefix + s
	}
	return res
}
