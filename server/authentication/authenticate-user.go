package authentication

import (
	"database/sql"
	"errors"
	"instant-messaging-platform-backend/config"
	_ "github.com/lib/pq"
)


func IsAuthValid(authToken string) (bool, error) {
	
	db, err := sql.Open("postgres", config.PostgresConfig)
	if err != nil {
		return false, err
	}
	defer db.Close()

	sqlRes, err := db.Exec("select * from " + config.CredsTable + " where uuid='" + authToken + "'")
	if err != nil {
		return false, err
	}
	noOfUsersWithThisAuth, err := sqlRes.RowsAffected()
	if err != nil {
		return false, err
	}
	if noOfUsersWithThisAuth == 1 {
		return true, nil
	}
	
	return false, errors.New("Invalid credentials")
}
