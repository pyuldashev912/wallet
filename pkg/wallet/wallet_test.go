package wallet

import (
	"errors"
	"testing"

	"github.com/pyuldashev912/wallet/pkg/types"
)

func TestService_FindAccountById_positive(t *testing.T) {
	s := Service{
		accounts: []*types.Account{
			{ID: 1, Phone: "+992927777777", Balance: 1549},
			{ID: 2, Phone: "+992550149898", Balance: 11549},
			{ID: 3, Phone: "+992929608008", Balance: 91514},
			{ID: 4, Phone: "+992929750076", Balance: 64752},
		},
	}

	result, _ := s.FindAccountById(2)
	expect := s.accounts[1].Phone

	if result.Phone != expect {
		t.Errorf("Expect %s, got %s", expect, result.Phone)
	}
}

func TestService_FindAccountById_negative(t *testing.T) {
	s := Service{
		accounts: []*types.Account{
			{ID: 1, Phone: "+992927777777", Balance: 1549},
			{ID: 2, Phone: "+992550149898", Balance: 11549},
			{ID: 3, Phone: "+992929608008", Balance: 91514},
			{ID: 4, Phone: "+992929750076", Balance: 64752},
		},
	}

	_, err := s.FindAccountById(5)

	if !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("Expect '%v' error, got '%v' error", ErrAccountNotFound, err)
	}
}

func TestService_FindPaymentById_positive(t *testing.T) {
	s := Service{
		payments: []*types.Payment{
			{ID: "1", AccountID: 1, Amount: 47575, Category: "Home", Status: types.StatusInProgress},
			{ID: "2", AccountID: 2, Amount: 25475, Category: "Drugstore", Status: types.StatusInProgress},
			{ID: "3", AccountID: 3, Amount: 6855, Category: "Coffee", Status: types.StatusInProgress},
		},
	}

	expectedID := "2"
	got, _ := s.FindPaymentById(expectedID)
	if got.ID != expectedID {
		t.Errorf("Expect %s, got %s", expectedID, got.ID)
	}
}

func TestService_FindPaymentById_negative(t *testing.T) {
	s := Service{
		payments: []*types.Payment{
			{ID: "1", AccountID: 1, Amount: 47575, Category: "Home", Status: types.StatusInProgress},
			{ID: "2", AccountID: 2, Amount: 25475, Category: "Drugstore", Status: types.StatusInProgress},
			{ID: "3", AccountID: 3, Amount: 6855, Category: "Coffee", Status: types.StatusInProgress},
		},
	}

	_, err := s.FindPaymentById("5")
	if !errors.Is(err, ErrPaymentNotFound) {
		t.Errorf("Expect '%v' error, got '%v' error", ErrPaymentNotFound, err)
	}
}

func TestService_Reject_positive(t *testing.T) {
	s := Service{
		accounts: []*types.Account{
			{ID: 1, Phone: "+992927777777", Balance: 1549},
			{ID: 2, Phone: "+992550149898", Balance: 11549},
			{ID: 3, Phone: "+992929608008", Balance: 91514},
		},
		payments: []*types.Payment{
			{ID: "1", AccountID: 1, Amount: 47575, Category: "Home", Status: types.StatusInProgress},
			{ID: "2", AccountID: 2, Amount: 25475, Category: "Drugstore", Status: types.StatusInProgress},
			{ID: "3", AccountID: 3, Amount: 6855, Category: "Coffee", Status: types.StatusInProgress},
		},
	}

	s.Reject("3")
	expectedBalance := types.Money(98369)
	if expectedBalance != s.accounts[2].Balance {
		t.Errorf("Expected balance %d, got %d", expectedBalance, s.accounts[2].Balance)
	}

	if s.payments[2].Status != types.StatusFail {
		t.Errorf("Expected payment status %v, got %v", types.StatusFail, s.payments[2].Status)
	}
}

func TestService_Reject_negative(t *testing.T) {
	s := Service{
		accounts: []*types.Account{
			{ID: 1, Phone: "+992927777777", Balance: 1549},
			{ID: 2, Phone: "+992550149898", Balance: 11549},
			{ID: 3, Phone: "+992929608008", Balance: 91514},
		},
		payments: []*types.Payment{
			{ID: "1", AccountID: 1, Amount: 47575, Category: "Home", Status: types.StatusInProgress},
			{ID: "2", AccountID: 2, Amount: 25475, Category: "Drugstore", Status: types.StatusInProgress},
			{ID: "3", AccountID: 3, Amount: 6855, Category: "Coffee", Status: types.StatusInProgress},
		},
	}

	err := s.Reject("5")
	if !errors.Is(err, ErrPaymentNotFound) {
		t.Errorf("Expect '%v' error, got '%v' error", ErrPaymentNotFound, err)
	}
}
