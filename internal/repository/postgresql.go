package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"testApi/internal/model"
	"testApi/public_dto"
	"time"
)

type Pool struct {
	*pgxpool.Pool
}

func GetConnectionPool(connectionStr string) (*Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connPool, err := pgxpool.New(ctx, connectionStr)
	if err != nil {
		log.Println("Error creating connection pool:", err)
		return nil, err
	}
	return &Pool{connPool}, nil
}

func (pool *Pool) AddUser(user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pool.Pool.Exec(ctx, "INSERT INTO users (username,surname,age,nationality,gender) VALUES ($1,$2,$3,$4,$5);",
		user.UserName, user.Surname, user.Age, user.Nationality, user.Gender)
	if err != nil {
		return err
	}
	return nil
}

func (pool *Pool) UpdateUser(id int, user public_dto.NewUserData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exec, err := pool.Pool.Exec(ctx, "UPDATE users SET username = $1, surname = $2, age = $3, nationality = $4, gender = $5 WHERE id = $6",
		user.Name, user.Surname, user.Age, user.Nationality, user.Gender, id)
	if err != nil {
		return err
	}
	if exec.RowsAffected() == 0 {
		return errors.New("user does not exist")
	}
	return nil
}

func (pool *Pool) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exec, err := pool.Pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	if exec.RowsAffected() == 0 {
		return errors.New("user does not exist")
	}

	return nil
}

func (pool *Pool) GetFilteredUsers(filter public_dto.UserFilter) (public_dto.PagedUsers, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sqlString = "SELECT id,username,surname,age,nationality,gender FROM users WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM users WHERE 1=1"
	args := []interface{}{}
	argsIndex := 1

	if filter.Username != "" {
		sqlString += fmt.Sprintf(" AND username ILIKE $%d", argsIndex)
		args = append(args, "%"+filter.Username+"%")
		argsIndex++
	}
	if filter.Surname != "" {
		sqlString += fmt.Sprintf(" AND surname = $%d", argsIndex)
		args = append(args, filter.Surname)
		argsIndex++
	}
	if filter.Age != 0 {
		sqlString += fmt.Sprintf(" AND age = $%d", argsIndex)
		args = append(args, filter.Age)
		argsIndex++
	}
	if filter.Nationality != "" {
		sqlString += fmt.Sprintf(" AND nationality = $%d", argsIndex)
		args = append(args, filter.Nationality)
		argsIndex++
	}
	if filter.Gender != "" {
		sqlString += fmt.Sprintf(" AND gender = $%d", argsIndex)
		args = append(args, filter.Gender)
		argsIndex++
	}

	var total int
	err := pool.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return public_dto.PagedUsers{}, err
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	offset := (filter.Page - 1) * filter.PageSize

	sqlString += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argsIndex, argsIndex+1)
	argsIndex += 2
	args = append(args, filter.PageSize, offset)

	rows, err := pool.Pool.Query(ctx, sqlString, args...)
	if err != nil {
		return public_dto.PagedUsers{}, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.Id, &user.UserName, &user.Surname, &user.Age, &user.Nationality, &user.Gender)
		if err != nil {
			return public_dto.PagedUsers{}, err
		}
		users = append(users, user)
	}
	return public_dto.PagedUsers{
		Users:      users,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: (total + filter.PageSize - 1) / filter.PageSize,
	}, nil
}
