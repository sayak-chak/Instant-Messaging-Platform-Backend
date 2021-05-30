package database

import (
	"database/sql"
	"instant-messaging-platform-backend/config"
	_ "github.com/lib/pq"
)

func SetupDataBase() error {
	db, err := sql.Open("postgres", config.PostgresConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("create table if not exists " + config.UsersTable + "(username varchar(20), password varchar(20), PRIMARY KEY (username))")
	if err != nil {
		return err
	}
	_, err = db.Exec("create table if not exists " + config.CredsTable + "(uuid varchar(36), username varchar(20), PRIMARY KEY (uuid))") //varchar 36 to store uuid
	if err != nil {
		return err
	}
	_, err = db.Exec("create table if not exists " + config.ChatTable + "(sender varchar(20), room varchar(50), message varchar(500))") 
	return err
}
