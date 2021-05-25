package mocks

import (
	"database/sql"
	"encoding/json"
	"instant-messaging-platform-backend/config"
)

func GetLoginJson(username string, password string) ([]byte, error){
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