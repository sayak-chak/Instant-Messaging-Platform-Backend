package database

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	_ "github.com/lib/pq"
	"instant-messaging-platform-backend/config"
	"instant-messaging-platform-backend/database/model"
)

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
			Temp:          false,
			IfNotExists:   true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
