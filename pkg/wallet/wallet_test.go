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
