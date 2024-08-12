package postdb

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func init() {
	const connStr = "postgres://postgres:2412@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS guidbase (
			id SERIAL PRIMARY KEY,
			GUID CHARACTER VARYING(255) NOT NULL,
			refresh_token CHARACTER VARYING(255) NOT NULL
		);
	`)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func AddUser(GUID, refreshToken string) {
	const connStr = "postgres://postgres:2412@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Добавление данных в таблицу
	_, err = db.Exec(
		`INSERT INTO guidbase (GUID, refresh_token) 
            VALUES ($1, $2)
            `, GUID, refreshToken)
	if err != nil {
		panic(err)
	}
}

func GetUser(GUID string) (string, error) {
	const connStr = "postgres://postgres:2412@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var refreshToken string
	err = db.QueryRow(`SELECT refresh_token FROM guidbase WHERE GUID = $1`, GUID).Scan(&refreshToken)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func DeleteUser(GUID string) {
	const connStr = "postgres://postgres:2412@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM guidbase WHERE GUID = $1`, GUID)
	if err != nil {
		panic(err)
	}
}
