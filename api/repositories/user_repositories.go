package repositories

import (
	"context"
	"database/sql"
	"errors"

	"codemead.com/go_fintech/fintech_backend/api/models"
	db "codemead.com/go_fintech/fintech_backend/db/sqlc"
	"codemead.com/go_fintech/fintech_backend/token"
	"codemead.com/go_fintech/fintech_backend/utils"
)

type UserRepository struct {
	DB         *db.Store
	Config     *utils.Config
	TokenMaker token.Maker
}

func NewUserRepository(queries *db.Store, config *utils.Config, maker token.Maker) *UserRepository {
	return &UserRepository{
		DB:         queries,
		Config:     config,
		TokenMaker: maker,
	}
}

func (r *UserRepository) ListUsers() ([]models.ResponseUser, error) {
	arg := db.ListUsersParams{
		Offset: 0,
		Limit:  10,
	}

	users, err := r.DB.ListUsers(context.Background(), arg)
	if err != nil {
		return nil, err
	}

	newUsers := []models.ResponseUser{}
	for _, v := range users {
		n := models.ToUserResponse(&v)
		newUsers = append(newUsers, *n)
	}
	return newUsers, nil
}

func (r *UserRepository) GetProfile(userId int64) (*models.ResponseUser, error) {
	user, err := r.DB.GetUserByID(context.Background(), userId)
	if err == sql.ErrNoRows {
		return nil, errors.New("not authorized to access resources")
	} else if err != nil {
		return nil, err
	}
	return models.ToUserResponse(&user), nil
}
