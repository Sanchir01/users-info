package app

import "github.com/Sanchir01/users-info/internal/feature/user"

type Services struct {
	UserService *user.Service
}

func NewServices(repos *Repositories, db *Database) *Services {
	return &Services{
		UserService: user.NewService(repos.UserRepository, db.PrimaryDB),
	}
}
