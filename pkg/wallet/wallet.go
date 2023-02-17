package wallet

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/pyuldashev912/wallet/pkg/types"
)

const (
	sep string = "/"
)

var (
	ErrPhoneRegistered      = errors.New("phone has already registered")
	ErrAmountMustBePositive = errors.New("amount must be positive")
	ErrNotEnoughBalance     = errors.New("not enough balance")
	ErrAccountNotFound      = errors.New("account not found")
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrFavoriteNotFound     = errors.New("favorite payment not found")
)

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

// RegisterAccount возвращает зарегистрированный аккаунт
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

// Deposit поплняет баланс пользователя
func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount < 0 {
		return ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}

// Pay производит платеж
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount < 0 {
		return nil, ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	payment := &types.Payment{
		ID:        uuid.New().String(),
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.StatusInProgress,
	}

	s.payments = append(s.payments, payment)
	return payment, nil
}

// FindAccountByID возвращает указатель на найденный аккаунт
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, acc := range s.accounts {
		if accountID == acc.ID {
			return acc, nil
		}
	}

	return nil, ErrAccountNotFound
}

// Reject отменяет платеж
func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}

	payment.Status = types.StatusFail
	// Пропускаем проверку на ошибку, так как платеж не может быть совершен
	// с несуществующего аккаунта
	account, _ := s.FindAccountByID(payment.AccountID)
	account.Balance += payment.Amount

	return nil
}

// FindPaymentByID возвращает указанный платеж
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, currPayment := range s.payments {
		if currPayment.ID == paymentID {
			return currPayment, nil
		}
	}

	return nil, ErrPaymentNotFound
}

// Repeat повторяет платеж
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return nil, err
	}

	if account.Balance < payment.Amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= payment.Amount
	result := &types.Payment{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Status:    payment.Status,
		Category:  payment.Category,
	}
	s.payments = append(s.payments, result)

	return result, nil
}

// FavoritePayment создает избранный платеж из ранее сделанного платежа
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	targetPayment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := &types.Favorite{
		ID:        uuid.New().String(),
		AccountID: targetPayment.AccountID,
		Name:      name,
		Amount:    targetPayment.Amount,
		Category:  targetPayment.Category,
	}

	s.favorites = append(s.favorites, favorite)

	return favorite, nil
}

func (s *Service) findPaymentFavorite(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}

	return nil, ErrFavoriteNotFound
}

// PayFromFavorite
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favPayment, err := s.findPaymentFavorite(favoriteID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(
		favPayment.AccountID, favPayment.Amount, favPayment.Category,
	)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *Service) ExportToFile(path string) error {
	file, err := prepareFoldersAndFile(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, account := range s.accounts {
		res := fmt.Sprintf("%d;%s;%d\n", account.ID, account.Phone, account.Balance)
		_, err := file.Write([]byte(res))
		if err != nil {
			return err
		}
	}

	return nil
}

func prepareFoldersAndFile(path string) (*os.File, error) {
	pathSlice := strings.Split(path, sep)
	if len(pathSlice) > 1 {
		folders := strings.Join(pathSlice[:len(pathSlice)-1], sep)
		if err := os.MkdirAll(folders, os.ModePerm); err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return nil, err
	}

	return file, nil
}
