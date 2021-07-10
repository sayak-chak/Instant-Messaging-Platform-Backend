package authentication

import (
	"errors"
	"github.com/go-pg/pg/v10"
	_ "github.com/lib/pq"
	"instant-messaging-platform-backend/config"
	"instant-messaging-platform-backend/database/model"
)

func IsAuthValid(authToken string) (bool, error) {

	db := pg.Connect(&pg.Options{
		User:     config.User,
		Database: config.DbName,
	})
	defer db.Close()

	noOfUsersWithThisAuth, err := db.Model(&model.CredsTable{}).Where("uuid = ?", authToken).Count()
	if err != nil {
		return false, err
	}
	if noOfUsersWithThisAuth == 1 {
		return true, nil
	}

	return false, errors.New("invalid credentials")
}
