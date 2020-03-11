package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/nitrajka/paymentsFutured/cashdesk"
	"github.com/nitrajka/paymentsFutured/server"
)

func main() {
	port := flag.String("dbport", "5431", "use 5432 when running go ap in docker, else 5431")
	host := flag.String("dbhost", "localhost", "use database_postgres when running app in docker, else localhost")
	flag.Parse()

	dbCashDesk, err := cashdesk.NewDBCashDesk(*port, *host)
	if err != nil {
		exit(fmt.Sprintf("could not establish connection with cahdesk: %v", err))
	}
	fmt.Printf("connection with db established on port: %v on host: %v\n", *port, *host)

	s := server.NewPaymentServer(dbCashDesk)
	if err := http.ListenAndServe(":5000", s); err != nil {
		exit(fmt.Sprintf("could not listen on port 5000 %v", err))
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
