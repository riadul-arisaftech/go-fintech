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

type AuthRepository struct {
	DB         *db.Store
	Config     *utils.Config
	TokenMaker token.Maker
}

func NewAuthRepository(queries *db.Store, config *utils.Config, maker token.Maker) *AuthRepository {
	return &AuthRepository{
		DB:         queries,
		Config:     config,
		TokenMaker: maker,
	}
}

func (r *AuthRepository) Login(request models.RequestUserParams) (*models.ResponseAuthUser, error) {
	user, err := r.DB.GetUserByEmail(context.Background(), request.Email)
	if err == sql.ErrNoRows {
		return nil, errors.New("incorrect email or password")
	} else if err != nil {
		return nil, err
	}

	if err := utils.CheckPassword(request.Password, user.HashedPassword); err != nil {
		return nil, errors.New("incorrect password")
	}

	token, err := r.TokenMaker.CreateToken(user.ID, r.Config.ExpDuration)
	if err != nil {
		return nil, errors.New("Failed to create token :" + err.Error())
	}
	return &models.ResponseAuthUser{
		ID:        user.ID,
		Email:     user.Email,
		Token:     token,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *AuthRepository) Register(request models.RequestUserParams) (*models.ResponseUser, error) {
	hashedPassword, err := utils.HashedPassword(request.Password)
	if err != nil {
		return nil, err
	}

	arg := db.CreateUserParams{
		Email:          request.Email,
		HashedPassword: hashedPassword,
	}

	user, err := r.DB.CreateUser(context.Background(), arg)
	if err != nil {
		return nil, err
	}

	return models.ToUserResponse(&user), nil
}
