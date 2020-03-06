package cashdesk

import (
	"github.com/nitrajka/paymentsFutured/payment"
)

type CashDesk interface {
	GetPayment(id int) (payment.Payment, error)
	GetPayments() []payment.Payment
	SavePayment(payment payment.Payment) payment.Payment
	GetBalance() float64
}
