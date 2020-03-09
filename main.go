package main

import (
	"log"
	"net/http"

	"github.com/nitrajka/paymentsFutured/cashdesk"
	"github.com/nitrajka/paymentsFutured/server"
)

func main() {
	s := server.NewPaymentServer(cashdesk.NewInMemoryCashDesk())
	if err := http.ListenAndServe(":5000", s); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

//přijímat a evidovat platby.
//vylistovat všechny platby, jednotlivou platbu a přijímat platby
//GET /payments/     -> list all payments
//GET /payments/{id} -> list 1 payment
//POST payment       -> save to db and returning its id

//todo: use context, set up db, atomic, payment timestamp, vypisovat vsetko na konzolu, klient