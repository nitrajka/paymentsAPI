package server

import (
	"encoding/json"
	"fmt"
	"github.com/nitrajka/paymentsFutured/postgres"
	"net/http"
	"strconv"

	"github.com/gorilla/pat"
	"github.com/nitrajka/paymentsFutured/cashdesk"
)

type PaymentServer struct {
	cashDesk cashdesk.CashDesk
	http.Handler
}

func NewPaymentServer(cashDesk cashdesk.CashDesk) *PaymentServer {
	p := new(PaymentServer)
	p.cashDesk = cashDesk

	router := pat.New()
	router.Get("/payments/{id}", p.GetPayment)
	router.Get("/payments/", p.GetPayments)
	router.Post("/payments/", p.PostPayment)

	p.Handler = router
	return p
}

func NotFoundPaymentError(id int) string {
	return fmt.Sprintf("Payment with id: %v does not exist", id)
}

func InvalidBodyError(body string) string {
	return fmt.Sprintf("invalid parameters: %v\n", body)
}

func (p *PaymentServer) GetPayment(w http.ResponseWriter, r *http.Request) {
	paymentId := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(paymentId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, fmt.Sprintf("invalid id: %v", paymentId))
		return
	}

	payment, err := p.cashDesk.GetPayment(r.Context(), int32(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, NotFoundPaymentError(id))
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

func (p *PaymentServer) GetPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := p.cashDesk.GetPayments(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, fmt.Sprintf("could not response: %v", err))
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payments)
}

func (p *PaymentServer) PostPayment(w http.ResponseWriter, r *http.Request) {
	var paymParams postgres.CreatePaymentParams
	err := json.NewDecoder(r.Body).Decode(&paymParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, InvalidBodyError(fmt.Sprintf("check the fields of payment type: %v", err)))
		return
	}

	paym, err := p.cashDesk.SavePayment(r.Context(), paymParams)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paym)
}
