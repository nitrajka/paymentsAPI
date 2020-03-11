package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nitrajka/paymentsFutured/postgres"

	"github.com/nitrajka/paymentsFutured/cashdesk"
)

type stubStore struct {
	balance  float64
	payments map[int]postgres.Payment
}

func (s *stubStore) GetPayment(ctx context.Context, id int32) (postgres.Payment, error) {
	if val, ok := s.payments[int(id)]; ok {
		return val, nil
	}
	return postgres.Payment{}, fmt.Errorf("payment %v does not exist", id)
}

func (s *stubStore) GetPayments(ctx context.Context) ([]postgres.Payment, error) {
	var res []postgres.Payment
	for key := range s.payments {
		res = append(res, s.payments[key])
	}
	return res, nil
}

func (s *stubStore) SavePayment(ctx context.Context, paymentParams postgres.CreatePaymentParams) (postgres.Payment, error) {
	payment := postgres.Payment{
		ID:          int32(len(s.payments)),
		Amount:      paymentParams.Amount,
		Description: paymentParams.Description,
		Sender:      paymentParams.Sender,
		Datetime:    paymentParams.Datetime,
	}
	s.payments[len(s.payments)] = payment
	s.balance += payment.Amount
	return payment, nil
}

func (s *stubStore) GetBalance(ctx context.Context) (float64, error) {
	return s.balance, nil
}

func TestGetPayment(t *testing.T) {
	store := &stubStore{
		payments: map[int]postgres.Payment{
			0: {ID: 0, Amount: 10, Description: "fst payment", Sender: "anonymous", Datetime: time.Now()},
			1: {ID: 1, Amount: 10, Description: "snd payment", Sender: "anonymous", Datetime: time.Now()},
		}, balance: 0}

	server := NewPaymentServer(store)

	t.Run("test get existing payment", func(t *testing.T) {
		request := newGetPaymentRequest(t, 0)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		payment := getPaymentFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertPayment(t, payment, store.payments[0])
	})

	t.Run("test get non-existing payment", func(t *testing.T) {
		request := newGetPaymentRequest(t, len(store.payments))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("test get all payments", func(t *testing.T) {
		request := newGetPaymentsRequest(t)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertPayments(t, response.Body, store)
	})
}

func TestPostPayment(t *testing.T) {
	cashD := &stubStore{payments: map[int]postgres.Payment{
		0: {ID: 0, Amount: 10, Description: "fst payment", Sender: "anonymous", Datetime: time.Now()},
	}, balance: 0}
	server := NewPaymentServer(cashD)

	t.Run("test get created payment in response", func(t *testing.T) {
		tm := time.Now()
		p := postgres.CreatePaymentParams{Amount: 10, Description: "snd payment", Sender: "anonymous", Datetime: tm}
		expected := postgres.Payment{ID: 1, Amount: 10, Description: "snd payment", Sender: "anonymous", Datetime: tm}
		request := newPostPaymentRequest(t, p)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertPayment(t, getPaymentFromResponse(t, response.Body), expected)
		assertBalance(t, cashD.balance, 10)
	})

	t.Run("test save second payment", func(t *testing.T) {
		tm := time.Now()
		p := postgres.CreatePaymentParams{Amount: 10, Description: "trd payment", Sender: "anonymous", Datetime: tm}
		expected := postgres.Payment{ID: 2, Amount: 10, Description: "trd payment", Sender: "anonymous", Datetime: tm}
		request := newPostPaymentRequest(t, p)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertPayment(t, getPaymentFromResponse(t, response.Body), expected)

		request = newGetPaymentsRequest(t)
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertPayments(t, response.Body, cashD)
		assertBalance(t, cashD.balance, 20)
	})

	t.Run("test assert invalid payment when save", func(t *testing.T) {
		tm := time.Now().String() + "Z"
		p := struct {
			Amount float64
			Description int
			Sender string
			Datetime string
		}{Amount: 10, Description: 43, Sender: "anonymous", Datetime: tm}

		request := newPostPaymentRequestInvalid(t, p)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusBadRequest)
	})
}

func getPaymentFromResponse(t *testing.T, r io.Reader) (payment postgres.Payment) {
	t.Helper()
	err := json.NewDecoder(r).Decode(&payment)
	if err != nil {
		t.Errorf("error while decoding a payment: %v", err)
	}

	return
}

func assertStatus(t *testing.T, actual, expected int) {
	t.Helper()
	if actual != expected {
		t.Errorf("did not get correct status, actual %d, expected %d", actual, expected)
	}
}

func assertBalance(t *testing.T, actual, expected float64) {
	t.Helper()
	if actual != expected {
		t.Errorf("balance not right: expected %v, actual: %v", expected, actual)
	}
}

func assertPayments(t *testing.T, body io.Reader, store cashdesk.CashDesk) {
	t.Helper()
	var payments []postgres.Payment

	expected, err := store.GetPayments(context.Background())
	if err != nil {
		t.Errorf("could not fetch payments: %v", err)
	}

	err = json.NewDecoder(body).Decode(&payments)
	if err != nil {
		t.Errorf("could not decode received payments: %v", err)
	}

	if len(payments) != len(expected) {
		t.Errorf("did not receive all payments")
	}

	for i := range expected {
		found := false
		for j := range payments {
			if paymentsEqual(expected[i], payments[j]) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("did not receive all payments")
		}
	}
}

func paymentsEqual(p1, p2 postgres.Payment) bool {
	return p1.ID == p2.ID && p1.Description == p2.Description &&
		p1.Sender == p2.Sender && p1.Amount == p2.Amount &&
		p1.Datetime.Local().String() == p2.Datetime.Local().String()
}

func assertPayment(t *testing.T, actual, expected postgres.Payment) {
	t.Helper()
	if !paymentsEqual(actual, expected) {
		t.Errorf("payments are not equal: actual %v, expected %v", actual, expected)
	}
}

func newPostPaymentRequest(t *testing.T, p postgres.CreatePaymentParams) *http.Request {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(p)
	if err != nil {
		t.Errorf("something went wrong creating a request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/payments/", &body)
	if err != nil {
		t.Errorf("something went wrong creating a request: %v", err)
	}
	return req
}

func newPostPaymentRequestInvalid(t *testing.T, i interface{}) *http.Request {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(i)
	if err != nil {
		t.Errorf("something went wrong creating a request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/payments/", &body)
	if err != nil {
		t.Errorf("something went wrong creating a request: %v", err)
	}
	return req
}

func newGetPaymentRequest(t *testing.T, id int) *http.Request {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/payments/%v", id), nil)
	if err != nil {
		t.Errorf("something went wrong creating a request: %v", err)
	}
	return req
}

func newGetPaymentsRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "/payments/", nil)
	if err != nil {
		t.Errorf("something went wrong creating a request: %v", err)
	}
	return req
}
