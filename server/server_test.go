package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nitrajka/paymentsFutured/cashdesk"
	"github.com/nitrajka/paymentsFutured/payment"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nitrajka/paymentsFutured/cashdesk"
	"github.com/nitrajka/paymentsFutured/payment"
)

type stubStore struct {
	balance  float64
	payments map[int]payment.Payment
}

func (s *stubStore) GetPayment(id int) (payment.Payment, error) {
	if val, ok := s.payments[id]; ok {
		return val, nil
	}
	return payment.Payment{}, fmt.Errorf("payment %v does not exist", id)
}

func (s *stubStore) GetPayments() []payment.Payment {
	var res []payment.Payment
	for key := range s.payments {
		res = append(res, s.payments[key])
	}
	return res
}

func (s *stubStore) SavePayment(payment payment.Payment) payment.Payment {
	payment.Id = len(s.payments)
	s.payments[len(s.payments)] = payment
	s.balance += payment.Amount
	return payment
}

func (s *stubStore) GetBalance() float64 {
	return s.balance
}

func TestGetPayment(t *testing.T) {
	store := &stubStore{
		payments: map[int]payment.Payment{
			0: payment.NewPayment(0, "fst payment", "anonymous", 10),
			1: payment.NewPayment(1, "snd payment", "anonymous", 10),
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
		assertStatus(t, response.Code, http.StatusBadRequest)
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
	cashD := &stubStore{payments: map[int]payment.Payment{
		0: payment.NewPayment(0, "fst payment", "anonymous", 10),
	}, balance: 0}
	server := NewPaymentServer(cashD)

	t.Run("test get created payment in response", func(t *testing.T) {
		expected := payment.NewPayment(1, "snd payment", "anonymous", 10)
		request := newPostPaymentRequest(t, expected)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertPayment(t, getPaymentFromResponse(t, response.Body), expected)
		assertBalance(t, cashD.balance, 10)
	})

	t.Run("test save second payment", func(t *testing.T) {
		p := payment.NewPayment(0, "trd payment", "anonymous", 10)
		expected := payment.NewPayment(2, "trd payment", "anonymous", 10)
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
}

func getPaymentFromResponse(t *testing.T, r io.Reader) (payment payment.Payment) {
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
	var payments []payment.Payment
	expected := store.GetPayments()
	err := json.NewDecoder(body).Decode(&payments)
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

func paymentsEqual(p1 payment.Payment, p2 payment.Payment) bool {
	return p1.Id == p2.Id && p1.Description == p2.Description && p1.Sender == p2.Sender && p1.Amount == p2.Amount
}
func assertPayment(t *testing.T, actual, expected payment.Payment) {
	t.Helper()
	if !paymentsEqual(actual, expected) {
		t.Errorf("payments are not equal: actual %v, expected %v", actual, expected)
	}
}

func newPostPaymentRequest(t *testing.T, p payment.Payment) *http.Request {
	var body bytes.Buffer
	json.NewEncoder(&body).Encode(p)
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