package postdb

import (
	"Auth/internal/auth"
	"Auth/internal/middleware"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

const connStr = "postgres://postgres:1234@postgres:5432/mydb?sslmode=disable"

const initTable = `
		CREATE TABLE IF NOT EXISTS guidbase (
			id SERIAL PRIMARY KEY,
			ip CHARACTER VARYING(15) NOT NULL,
			guid CHARACTER VARYING(255) NOT NULL,
			refresh_token CHARACTER VARYING(255) NOT NULL,
		    expireat bigint NOT NULL,
		    createAt timestamp with time zone NOT NULL DEFAULT now()
		);
	`

const addUser = `INSERT INTO guidbase (ip, guid, refresh_token, expireat) VALUES ($1, $2, $3, $4)`

const getUser = `SELECT ip, refresh_token, expireat FROM guidbase WHERE guid = $1`

const deleteUser = `DELETE FROM guidbase WHERE guid = $1`

const updateUser = `UPDATE guidbase SET ip = $1, refresh_token = $2, expireat = $3 WHERE guid = $4`

// Инициализация БД
func init() {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(initTable)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

// Добавлени данных юзера в БД
func AddUser(ip string, dataRefresh *auth.RefreshToken, refreshToken string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return middleware.ErrDB
	}
	defer db.Close()
	_, err = db.Exec(addUser, ip, dataRefresh.Guid.GUID, refreshToken, dataRefresh.ExpireAt)
	if err != nil {
		return middleware.ErrDB
	}
	return nil
}

// Получение данных юзера из БД
func GetUser(GUID string) (*auth.ClientRefreshToken, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, middleware.ErrDB
	}
	defer db.Close()

	refresh := auth.ClientRefreshToken{}
	err = db.QueryRow(getUser, GUID).Scan(&refresh.Ip, &refresh.RefreshToken, &refresh.ExpireAt)
	if err != nil {
		return nil, middleware.ErrDB
	}
	return &refresh, nil
}

// Проверка есть ли юзер в БД
func VerifyUser(GUID string) (*auth.ClientRefreshToken, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, middleware.ErrDB
	}
	defer db.Close()

	refresh := auth.ClientRefreshToken{}
	err = db.QueryRow(getUser, GUID).Scan(&refresh.Ip, &refresh.RefreshToken, &refresh.ExpireAt)
	if refresh.RefreshToken == "" {
		return &refresh, nil
	} else {
		return &refresh, middleware.ErrHaveRefreshToken
	}
}

// Удаление юзера из БД
func DeleteUser(GUID string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(deleteUser, GUID)
	if err != nil {
		return middleware.ErrDB
	}
	return nil
}

func UpdateUser(ip string, dataRefresh *auth.RefreshToken, refreshToken string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return middleware.ErrDB
	}
	defer db.Close()
	_, err = db.Exec(updateUser, ip, refreshToken, dataRefresh.ExpireAt, dataRefresh.Guid.GUID)
	if err != nil {
		return middleware.ErrDB
	}
	return nil
}
