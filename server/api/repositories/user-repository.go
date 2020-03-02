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

// CreateUser method
func (repo *UserRepository) CreateUser(user models.User) (models.User, error) {
	trns, er := json.Marshal(user.Transactions)
	if er != nil {
		return user, er
	}
	tr := string(trns)
	statement := `
    insert into users (name, number,amount,transactions)
    values ($1, $2, $3, $4)
    returning id
  `
	var id int64
	err := repo.db.QueryRow(statement, user.Name, user.Number, user.Amount, tr).Scan(&id)
	if err != nil {
		return user, err
	}
	createdUser := user
	createdUser.ID = id
	return createdUser, nil
}

// GetAllUsers method
func (repo *UserRepository) GetAllUsers() ([]models.User, error) {
	query := `
    select id,name, number,amount,transactions
    from users`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getUsersFromRows(rows)
}

// GetUserByNameAndNumber method
func (repo *UserRepository) GetUserByNameAndNumber(name, number string) (models.User, error) {
	query := `
    select id,name, number,amount,transactions
    from users
    where number = $2 and name = $1
  `
	row := repo.db.QueryRow(query, name, number)
	var user models.User
	var trans string
	err := row.Scan(&user.ID, &user.Name, &user.Number, &user.Amount, &trans)
	if err != nil {
		return models.User{}, err
	}
	transBytes := []byte(trans)
	_ = json.Unmarshal(transBytes, &user.Transactions)
	return user, nil
}

// GetUserByID method
func (repo *UserRepository) GetUserByID(id int64) (models.User, error) {
	query := `
    select id,name, number,amount,transactions
    from users
    where id= $1
  `
	row := repo.db.QueryRow(query, id)
	var user models.User
	var trans string
	err := row.Scan(&user.ID, &user.Name, &user.Number, &user.Amount, &trans)
	if err != nil {
		return models.User{}, err
	}
	transBytes := []byte(trans)
	_ = json.Unmarshal(transBytes, &user.Transactions)
	return user, nil
}

// UpdateUser method
func (repo *UserRepository) UpdateUser(id int64, user models.User) error {
	trns, er := json.Marshal(user.Transactions)
	if er != nil {
		return er
	}
	tr := string(trns)
	// set transactions = transactions || $2
	query := `
    update users
    set transactions = $2,amount = $3
    where id=$1
  `

	res, err := repo.db.Exec(query, id, tr, user.Amount)
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

// DeleteUser method
func (repo *UserRepository) DeleteUser(id int64) error {
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

func getUsersFromRows(rows *sql.Rows) ([]models.User, error) {
	users := []models.User{}
	for rows.Next() {
		var user models.User
		var trans string
		err := rows.Scan(&user.ID, &user.Name, &user.Number, &user.Amount, &trans)
		if err != nil {
			return nil, err
		}
		transBytes := []byte(trans)
		_ = json.Unmarshal(transBytes, &user.Transactions)
		users = append(users, user)
	}
	return users, nil
}
