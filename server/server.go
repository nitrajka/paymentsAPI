package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/pat"
	"github.com/nitrajka/paymentsFutured/cashdesk"
	"github.com/nitrajka/paymentsFutured/payment"
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
	return fmt.Sprintf("invalid parameters %v:", body)
}

func (p *PaymentServer) GetPayment(w http.ResponseWriter, r *http.Request) {
	paymentId := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(paymentId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, fmt.Sprintf("invalid id: %v", paymentId))
		return
	}

	payment, err := p.cashDesk.GetPayment(id)
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
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p.cashDesk.GetPayments())
}

func (p *PaymentServer) PostPayment(w http.ResponseWriter, r *http.Request) {
	var paym payment.Payment
	err := json.NewDecoder(r.Body).Decode(&paym)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, InvalidBodyError("check the fields of payment type"))
		return
	}

	paym = p.cashDesk.SavePayment(paym)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paym)
}