package model

type UsersTable struct {
	Username string `pg:",pk"`
	Password string
}

type CredsTable struct {
	UUID string `pg:",pk"`
	Username string
}

type ChatTable struct {
	Sender string
	Room string
	Message string
}