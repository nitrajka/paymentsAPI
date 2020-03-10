package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nitrajka/paymentsFutured/cashdesk"
	"github.com/nitrajka/paymentsFutured/server"
)

func main() {
	dbCashDesk, err := cashdesk.NewDBCashDesk()
	if err != nil {
		fmt.Printf("could not establish connection with cahdesk: %v", err)
		os.Exit(1)
	}
	fmt.Println("connection with db established")
	s := server.NewPaymentServer(dbCashDesk)
	if err := http.ListenAndServe(":5000", s); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

//todo: vypisovat vsetko na konzolu