package postdb

import (
	"Auth/auth"
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
			ip CHARACTER VARYING(15) NOT NULL,
			guid CHARACTER VARYING(255) NOT NULL,
			refresh_token CHARACTER VARYING(255) NOT NULL,
		    expireat bigint NOT NULL,
		    createAt timestamp with time zone NOT NULL DEFAULT now()
		);
	`)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func AddUser(ip string, dataRefresh *auth.RefreshToken, refreshToken string) error {
	const connStr = "postgres://postgres:2412@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	// Добавление данных в таблицу
	_, err = db.Exec(
		`INSERT INTO guidbase (ip, guid, refresh_token, expireat) 
            VALUES ($1, $2, $3, $4)
            `, ip, dataRefresh.Guid.GUID, refreshToken, dataRefresh.ExpireAt)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(GUID string) (auth.ClientRefreshToken, error) {
	const connStr = "postgres://postgres:2412@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return auth.ClientRefreshToken{}, err
	}
	defer db.Close()

	refresh := auth.ClientRefreshToken{}
	err = db.QueryRow(`SELECT ip, refresh_token, expireat FROM guidbase WHERE guid = $1`, GUID).Scan(&refresh.Ip, &refresh.RefreshToken, &refresh.ExpireAt)
	if err != nil {
		return auth.ClientRefreshToken{}, err
	}
	return refresh, nil
}

func DeleteUser(GUID string) error {
	const connStr = "postgres://postgres:2412@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM guidbase WHERE guid = $1`, GUID)
	if err != nil {
		return err
	}
	return nil
}
