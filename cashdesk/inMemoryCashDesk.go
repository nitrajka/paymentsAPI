package cashdesk

import (
	"fmt"

	"github.com/nitrajka/paymentsFutured/payment"
)

func NewInMemoryCashDesk() *cashDesk {
	return &cashDesk{payments: make(map[int]payment.Payment), balance: 0}
}

type cashDesk struct {
	balance  float64
	payments map[int]payment.Payment
}

func (c *cashDesk) GetPayment(id int) (payment.Payment, error) {
	if val, ok := c.payments[id]; ok {
		return val, nil
	}
	return payment.Payment{}, fmt.Errorf("payment %v does not exist", id)
}

func (c *cashDesk) GetPayments() []payment.Payment {
	var res []payment.Payment
	for key := range c.payments {
		res = append(res, c.payments[key])
	}
	return res
}

func (c *cashDesk) SavePayment(payment payment.Payment) payment.Payment {
	payment.Id = len(c.payments)
	c.payments[len(c.payments)] = payment
	c.balance += payment.Amount
	return payment
}

func (c *cashDesk) GetBalance() float64 {
	return c.balance
}