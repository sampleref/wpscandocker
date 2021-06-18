package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

const wpscan = "/usr/bin/wpscan"

type SingleUrl struct {
	Action string `json:"Action"`
	Url    string `json:"Url"`
	Enum   string `json:"Enum"`
	Report string `json:"Report"`
}

func readCommandOutput(command string) string {
	log.Printf("Running command %s \n", command)
	c, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		fmt.Printf("Error while running command %v", err)
	}
	log.Printf("Returning output as %s \n", c)
	return string(c)
}

func checkSingleUrl(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var singleUrl SingleUrl
	err := json.Unmarshal(reqBody, &singleUrl)
	if err != nil {
		fmt.Printf("Error while parsing request json %v", err)
		return
	}

	switch singleUrl.Action {
	case "check":
		//Perform check
		command := wpscan + " --url " + singleUrl.Url + " --enumerate " + singleUrl.Enum
		output := readCommandOutput(command)
		//Set to singleUrl Report
		singleUrl.Report = output
		break
	case "update":
		//Update database
		out, err := exec.Command("pwd").Output()
		if err != nil {
			fmt.Printf("Error while update %v", err)
		} else {
			fmt.Println("Update Successfully Executed")
			output := string(out[:])
			fmt.Println(output)
		}
		break
	default:
		fmt.Printf("No valid action\n")
	}
	err = json.NewEncoder(w).Encode(singleUrl)
	if err != nil {
		fmt.Printf("Error while parsing response json %v", err)
		return
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/urlcheck", checkSingleUrl).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	log.Printf("Starting to serve on 8080...\n")
	log.Printf("POST to Url: http://<>:8080/urlcheck | { \"Action\":\"check|update\", \"Url\":\"http://usablewp.com\", \"Enum\":\"args for enumerate\", \"Report\":\"\"} \n")
	log.Printf("Reply will be as { \"Url\":\"http://usablewp.com\", \"Report\":\"<JSON Report Updated here>\"} \n")
	handleRequests()
}
