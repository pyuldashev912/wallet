package types

// Money представляет собой денежную единицу в минимальных единицах
// (центы, копейки, дирамы и т.д.)
type Money int64

// PaymentCategory представляет собой категорию, в которой был совершён платёж
// (авто, аптеки, рестораны и т.д.)
type PaymentCategory string

// PaymentStatus представляет собой статус платежа
type PaymentStatus string

const (
	StatusOk         PaymentStatus = "OK"
	StatusFail       PaymentStatus = "Fail"
	StatusInProgress PaymentStatus = "INPROGRESS"
)

// Payment представляет информацию о платеже
type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

type Phone string

// Account представляет информацию о счёте пользователя
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}
