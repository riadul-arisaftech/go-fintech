package models

type AccountRequest struct {
	Currency string `json:"currency" validate:"required,currency"`
}

type TransferRequest struct {
	FromAccountID int32   `json:"from_account_id" validate:"required"`
	ToAccountID   int32   `json:"to_account_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required"`
}
