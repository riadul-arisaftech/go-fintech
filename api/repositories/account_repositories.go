package repositories

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"codemead.com/go_fintech/fintech_backend/api/models"
	db "codemead.com/go_fintech/fintech_backend/db/sqlc"
	"codemead.com/go_fintech/fintech_backend/token"
	"codemead.com/go_fintech/fintech_backend/utils"
)

type AccountRepository struct {
	DB         *db.Store
	Config     *utils.Config
	TokenMaker token.Maker
}

func NewAccountRepository(queries *db.Store, config *utils.Config, maker token.Maker) *AccountRepository {
	return &AccountRepository{
		DB:         queries,
		Config:     config,
		TokenMaker: maker,
	}
}

func (r *AccountRepository) GetUserAccounts(userId int64) ([]db.Account, error) {
	accounts, err := r.DB.GetAccountByUserID(context.Background(), int32(userId))
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *AccountRepository) CreateAccount(userId int64, request models.AccountRequest) (*db.Account, error) {
	arg := db.CreateAccountParams{
		UserID:   int32(userId),
		Currency: request.Currency,
		Balance:  0,
	}

	account, err := r.DB.CreateAccount(context.Background(), arg)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) CreateTransfer(userId int64, request models.TransferRequest) (*db.TransferTxResponse, int, error) {
	account, err := r.DB.GetAccountByID(context.Background(), int64(request.FromAccountID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusBadRequest, errors.New("couldn't get account")
		}
		return nil, http.StatusInternalServerError, err
	}

	if account.UserID != int32(userId) {
		return nil, http.StatusBadRequest, errors.New("couldn't get account")
	}

	toAccount, err := r.DB.GetAccountByID(context.Background(), int64(request.ToAccountID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusBadRequest, errors.New("couldn't find account to send to")
		}
		return nil, http.StatusInternalServerError, err
	}

	if toAccount.Currency != account.Currency {
		return nil, http.StatusBadRequest, errors.New("currencies of accounts do not match")
	}

	toArg := db.CreateTransferParams{
		FromAccountID: request.FromAccountID,
		ToAccountID:   request.ToAccountID,
		Amount:        request.Amount,
	}
	tx, err := r.DB.TransferTx(context.Background(), toArg)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("encountered issue with transaction")
	}

	return &tx, http.StatusCreated, nil
}
