package token

import "time"

type Maker interface {
	CreateToken(userId int64, duration time.Duration) (string, error)
	VerifyToken(token string) (int64, error)
}
