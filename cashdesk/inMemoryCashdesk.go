package cashdesk

import (
	"fmt"

	"github.com/nitrajka/paymentsFutured/postgres"
)

func NewInMemoryCashDesk() *cashDesk {
	return &cashDesk{payments: make(map[int32]postgres.Payment), balance: 0}
}

type cashDesk struct {
	balance  float64
	payments map[int32]postgres.Payment
}

func (c *cashDesk) GetPayment(id int32) (postgres.Payment, error) {
	if val, ok := c.payments[id]; ok {
		return val, nil
	}
	return postgres.Payment{}, fmt.Errorf("payment %v does not exist", id)
}

func (c *cashDesk) GetPayments() []postgres.Payment {
	var res []postgres.Payment
	for key := range c.payments {
		res = append(res, c.payments[key])
	}
	return res
}

func (c *cashDesk) SavePayment(paymentParams postgres.CreatePaymentParams) postgres.Payment {
	payment := postgres.Payment{
		ID:          int32(len(c.payments)),
		Amount:      paymentParams.Amount,
		Description: paymentParams.Description,
		Sender:      paymentParams.Sender,
		Datetime:    paymentParams.Datetime,
	}
	c.payments[int32(len(c.payments))] = payment
	c.balance += payment.Amount
	return payment
}

func (c *cashDesk) GetBalance() float64 {
	return c.balance
}
