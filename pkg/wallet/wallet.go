package wallet

import (
	"errors"

	"github.com/pyuldashev912/wallet/pkg/types"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrPaymentNotFound = errors.New("payment not found")
)

type Service struct {
	accounts []*types.Account
	payments []*types.Payment
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

// Reject ...
func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentById(paymentID)
	if err != nil {
		return err
	}

	payment.Status = types.StatusFail
	// Пропускаем проверку на ошибку, так как платеж не может быть совершен
	// с несуществующего аккаунта
	account, _ := s.FindAccountById(payment.AccountID)
	account.Balance += payment.Amount

	return nil
}

// FindPaymentById
func (s *Service) FindPaymentById(paymentID string) (*types.Payment, error) {
	for _, currPayment := range s.payments {
		if currPayment.ID == paymentID {
			return currPayment, nil
		}
	}

	return nil, ErrPaymentNotFound
}
