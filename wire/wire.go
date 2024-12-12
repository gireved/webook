//go:build wireinject

package wire

import (
	"geektime-basic-go/wire/repository"
	"geektime-basic-go/wire/repository/dao"
	"github.com/google/wire"
)

func InitUserRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, InitDB, dao.NewUserDAO)
	return &repository.UserRepository{}
}
