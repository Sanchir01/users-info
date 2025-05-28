package app

import "github.com/Sanchir01/users-info/internal/feature/user"

type Repositories struct {
	UserRepository *user.Repository
}

func NewRepositories(databases *Database) *Repositories {
	return &Repositories{
		UserRepository: user.NewRepository(databases.PrimaryDB),
	}
}
