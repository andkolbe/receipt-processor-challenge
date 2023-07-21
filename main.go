package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/receipts/:id/points", getPoints)
	router.HandlerFunc(http.MethodPost, "/receipts/process", processReceipts)

	log.Print("Starting server on 4000")
	err := http.ListenAndServe("127.0.0.1:4000", router)
	log.Fatal(err)
}