package cashdesk

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/nitrajka/paymentsFutured/postgres"
	"math"

	_ "github.com/lib/pq" // here
)

type dbCashDesk struct {
	db *postgres.Queries
}


func NewDBCashDesk() (*dbCashDesk, error) {
	conn, err := sql.Open(
		"postgres",
		"user=postgres password=password dbname=dev port=5432 host=database_postgres sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %v", err)
	}

	db := postgres.New(conn)
	return &dbCashDesk{db: db}, nil
}

func (d *dbCashDesk) GetPayment(ctx context.Context, id int32) (postgres.Payment, error) {
	p, err := d.db.GetPayment(ctx , id)
	if err != nil {
		return postgres.Payment{}, fmt.Errorf("error retrieving payment with id %v: %v", id, err)
	}
	return p, nil
}

func (d *dbCashDesk) GetPayments(ctx context.Context) ([]postgres.Payment, error) {
	p, err := d.db.GetPayments(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving payments: %v", err)
	}
	return p, nil
}

func (d *dbCashDesk) SavePayment(ctx context.Context, paymentParams postgres.CreatePaymentParams) (postgres.Payment, error) {
	p, err := d.db.CreatePayment(ctx, paymentParams)
	if err != nil {
		return postgres.Payment{}, fmt.Errorf("error creating payment: %v", err)
	}
	return p, nil
}

func (d *dbCashDesk) GetBalance(ctx context.Context) (float64, error) {
	balance, err := d.db.GetBalance(ctx)
	if err != nil {
		return 0, fmt.Errorf("error retrieving balance: %v", err)
	}

	return interfaceToFloat64(balance)
}

func interfaceToFloat64(i interface{}) (float64, error) {
	switch i := i.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	default:
		return math.NaN(), fmt.Errorf("getFloat: unknown value is of incompatible type")
	}
}