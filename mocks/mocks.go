package mocks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"instant-messaging-platform-backend/config"
)

func AuthValidator(authToken string) (bool, error) { return true, nil }

func GetLoginJson(username string, password string) ([]byte, error) {
	jsonString := map[string]interface{}{
		"username": username,
		"password": password,
	}
	return json.Marshal(jsonString)
}

func FillAuthTokenInDB(uuid string, username string) error {
	db, err := sql.Open("postgres", config.PostgresConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("insert into " + config.CredsTable + " values ('" + uuid + "','" + username + "')")
	return err
}

func RegisterUser(username string) error {
	db, err := sql.Open("postgres", config.PostgresConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("insert into " + config.UsersTable + " values ('" + username + "','" + "MOCK_PASSWORD" + "')")
	return err
}

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

func TearDown() {
	db, err := sql.Open("postgres", config.PostgresConfig)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	_, err = db.Exec("drop table " + config.CredsTable)
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("drop table " + config.UsersTable)
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("drop table " + config.ChatTable)
	if err != nil {
		fmt.Println(err)
	}
}
