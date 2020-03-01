package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"notebook/api/models"
)

// UserRepository struct
type UserRepository struct {
	db *sql.DB
}

// Init method
func (repo *UserRepository) Init(db *sql.DB) {
	repo.db = db
	repo.db.Query("CREATE SCHEMA IF NOT EXISTS saiteja")
	repo.db.Query(`CREATE TABLE IF NOT EXISTS users
	(
		id SERIAL PRIMARY KEY,
		number VARCHAR(50) NOT NULL,
		name VARCHAR(200) NOT NULL,
		amount FLOAT default (0),
		transactions jsonb NOT NULL default('[]')
	)`)

}

// CreateTask method
func (repo *UserRepository) CreateTask(task models.User) (models.User, error) {
	trns, er := json.Marshal(task.Transactions)
	if er != nil {
		return task, er
	}
	tr := string(trns)
	statement := `
    insert into users (name, number,amount,transactions)
    values ($1, $2, $3, $4)
    returning id
  `
	var id int64
	err := repo.db.QueryRow(statement, task.Name, task.Number, task.Amount, tr).Scan(&id)
	if err != nil {
		return task, err
	}
	createdTask := task
	createdTask.ID = id
	return createdTask, nil
}

// GetAllTasks method
func (repo *UserRepository) GetAllTasks() ([]models.User, error) {
	query := `
    select *
    from users`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getTasksFromRows(rows)
}

// GetTaskByNameAndNumber method
func (repo *UserRepository) GetTaskByNameAndNumber(name, number string) (models.User, error) {
	query := `
    select *
    from users
    where number = $2 and name = $1
  `
	row := repo.db.QueryRow(query, name, number)
	var task models.User
	err := row.Scan(&task.ID, &task.Name, &task.Number, &task.Amount, &task.Transactions)
	if err != nil {
		return models.User{}, err
	}
	return task, nil
}

// GetTaskByID method
func (repo *UserRepository) GetTaskByID(id int64) (models.User, error) {
	query := `
    select *
    from users
    where id= $1
  `
	row := repo.db.QueryRow(query, id)
	var task models.User
	var trans string
	err := row.Scan(&task.ID, &task.Name, &task.Number, &task.Amount, &trans)
	if err != nil {
		return models.User{}, err
	}
	transBytes := []byte(trans)
	_ = json.Unmarshal(transBytes, &task.Transactions)
	return task, nil
}

// UpdateTask method
func (repo *UserRepository) UpdateTask(id int64, task models.User) error {
	trns, er := json.Marshal(task.Transactions)
	if er != nil {
		return er
	}
	tr := string(trns)
	// set transactions = transactions || $2
	query := `
    update users
    set transactions = $2
    where id=$1
  `

	res, err := repo.db.Exec(query, id, tr)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got updated for %d", id)
	}

	return nil
}

// DeleteTask method
func (repo *UserRepository) DeleteTask(id int64) error {
	query := `delete from users where id=$1`
	res, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("exactly 1 row is not impacted for %d", id)
	}

	return nil
}

func getTasksFromRows(rows *sql.Rows) ([]models.User, error) {
	tasks := []models.User{}
	for rows.Next() {
		var task models.User
		var trans string
		err := rows.Scan(&task.ID, &task.Name, &task.Number, &task.Amount, &trans)
		if err != nil {
			return nil, err
		}
		transBytes := []byte(trans)
		_ = json.Unmarshal(transBytes, &task.Transactions)
		tasks = append(tasks, task)
	}
	return tasks, nil
}
