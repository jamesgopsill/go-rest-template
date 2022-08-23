package db

type User struct {
	ID           string
	Name         string
	Email        string
	Scopes       GormStringArray
	PasswordHash string
	Thumbnail    string
}
