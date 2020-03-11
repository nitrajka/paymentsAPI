package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strconv"

	"github.com/nitrajka/paymentsFutured/postgres"

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
	router.Get("/payments/{id}", logHandler(p.GetPayment))
	router.Get("/payments/", logHandler(p.GetPayments))
	router.Post("/payments/", logHandler(p.PostPayment))

	p.Handler = router
	return p
}

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		log.Println(fmt.Sprintf("%q", x))
		rec := httptest.NewRecorder()
		fn(rec, r)
		log.Println(fmt.Sprintf("%q", rec.Body))

		// this copies the recorded response to the response writer
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)
	}
}

func NotFoundPaymentError(err error) string {
	return fmt.Sprintf("Payment with id: %v does not exist.\n", err)
}

func InvalidBodyError(err error) string {
	return fmt.Sprintf("Could not create payment, invalid body parameters: %v.\n", err)
}

func InternalServerError(err error) string {
	return fmt.Sprintf("oops, something went wrong, try later: %v.\n", err)
}

func (p *PaymentServer) GetPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(paymentID)

	if p.checkErrAndMaybeFailResponse(
		err, NotFoundPaymentError(fmt.Errorf("payment with id %v not found", paymentID)),
		http.StatusNotFound, w) {
		return
	}

	payment, err := p.cashDesk.GetPayment(r.Context(), int32(id))
	if !p.checkErrAndMaybeFailResponse(
		err, NotFoundPaymentError(fmt.Errorf("payment with id %v not found", paymentID)),
		http.StatusNotFound, w) {
		p.encodeJSONAndMaybeSucceed(payment, w)
	}
}

func (p *PaymentServer) GetPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := p.cashDesk.GetPayments(r.Context())
	if p.checkErrAndMaybeFailResponse(err, InternalServerError(err), http.StatusInternalServerError, w) {
		return
	}

	p.encodeJSONAndMaybeSucceed(payments, w)
}

func (p *PaymentServer) PostPayment(w http.ResponseWriter, r *http.Request) {
	var paymParams postgres.CreatePaymentParams
	err := json.NewDecoder(r.Body).Decode(&paymParams)
	if p.checkErrAndMaybeFailResponse(err, InvalidBodyError(err), http.StatusBadRequest, w) {
		return
	}

	payment, err := p.cashDesk.SavePayment(r.Context(), paymParams)
	if !p.checkErrAndMaybeFailResponse(err, InternalServerError(err), http.StatusInternalServerError, w) {
		p.encodeJSONAndMaybeSucceed(payment, w)
	}
}

func (p *PaymentServer) encodeJSONAndMaybeSucceed(i interface{}, w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
	err := json.NewEncoder(w).Encode(i)

	if !p.checkErrAndMaybeFailResponse(
		err, InternalServerError(fmt.Errorf("could not encode response: %v", err)),
		http.StatusInternalServerError, w) {
		w.WriteHeader(http.StatusOK)
	}
}

func (p *PaymentServer) checkErrAndMaybeFailResponse(err error, responseError string, status int, w http.ResponseWriter) (failed bool) {
	if err != nil {
		w.WriteHeader(status)
		fmt.Fprintf(w, responseError)
		return true
	}
	return false
}
