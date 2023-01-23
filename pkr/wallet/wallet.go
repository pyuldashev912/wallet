package wallet

import (
	"errors"

	"github.com/pyuldashev912/wallet/pkr/types"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

type Service struct {
	accounts []*types.Account
}

// FindAccountById возвращает указатель на найденный аккаунт
func (s *Service) FindAccountById(accountId int64) (*types.Account, error) {
	for _, acc := range s.accounts {
		if accountId == acc.ID {
			return acc, nil
		}
	}

	return nil, ErrAccountNotFound
}
