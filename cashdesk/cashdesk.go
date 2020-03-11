package cashdesk

import (
	"context"

	"github.com/nitrajka/paymentsFutured/postgres"
)

type CashDesk interface {
	GetPayment(ctx context.Context, id int32) (postgres.Payment, error)
	GetPayments(ctx context.Context) ([]postgres.Payment, error)
	SavePayment(ctx context.Context, payment postgres.CreatePaymentParams) (postgres.Payment, error)
	GetBalance(ctx context.Context) (float64, error)
}
