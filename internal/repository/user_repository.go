package repository

import (
	"testApi/internal/model"
	"testApi/public_dto"
)

type UserRepository interface {
	AddUser(user model.User) error
	UpdateUser(id int, user public_dto.NewUserData) error
	DeleteUser(id int) error
	GetFilteredUsers(filter public_dto.UserFilter) (public_dto.PagedUsers, error)
}
