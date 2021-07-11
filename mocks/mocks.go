package mocks

import (
	"encoding/json"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	_ "github.com/lib/pq"
	"instant-messaging-platform-backend/config"
	"instant-messaging-platform-backend/database/model"
)

const MOCK_PASSWORD = "MOCK_PASSWORD"
func AuthValidator(authToken string) (bool, error) { return true, nil }

func GetLoginJson(username string, password string) ([]byte, error) {
	jsonString := map[string]interface{}{
		"username": username,
		"password": password,
	}
	return json.Marshal(jsonString)
}

func FillAuthTokenInDB(uuid string, username string) error {
	db := pg.Connect(&pg.Options{
		User:     config.User,
		Database: config.DbName,
	})
	defer db.Close()

	_, err := db.Model(&model.CredsTable{
		UUID:     uuid,
		Username: username,
	}).
		Insert()
	//Exec("insert into " + config.CredsTable + " values ('" + uuid + "','" + username + "')")
	return err
}

func RegisterUser(username string) error {
	db := pg.Connect(&pg.Options{
		User:     config.User,
		Database: config.DbName,
	})
	defer db.Close()

	_, err := db.Model(&model.UsersTable{
		username,
		MOCK_PASSWORD,
	}).Insert()
	//Exec("insert into " + config.UsersTable + " values ('" + username + "','" + "MOCK_PASSWORD" + "')")
	return err
}

func SetupDataBase() error {
	db := pg.Connect(&pg.Options{
		User:     config.User,
		Database: config.DbName,
	})
	defer db.Close()

	modelList := []interface{}{
		(*model.UsersTable)(nil),
		(*model.CredsTable)(nil),
		(*model.ChatTable)(nil),
	}

	for _, tableModel := range modelList {
		err := db.Model(tableModel).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func TearDown() error {
	db := pg.Connect(&pg.Options{
		User:     config.User,
		Database: config.DbName,
	})
	defer db.Close()

	modelList := []interface{}{
		(*model.UsersTable)(nil),
		(*model.CredsTable)(nil),
		(*model.ChatTable)(nil),
	}

	for _, tableModel := range modelList {
		err := db.Model(tableModel).DropTable(&orm.DropTableOptions{
			Cascade: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
