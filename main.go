package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"wpscandocker/wpscan"
)

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/checkurl", wpscan.CheckSingleUrl).Methods("POST")
	router.HandleFunc("/updatedb", wpscan.UpdateDatabase).Methods("GET")
	router.HandleFunc("/getreportbyid", wpscan.GetReportById).Methods("POST")
	router.HandleFunc("/getallreports", wpscan.GetAllReports).Methods("GET")
	router.HandleFunc("/deletereportbyid", wpscan.DeleteReport).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	log.Printf("Starting to serve on 8080...\n")
	handleRequests()
}
