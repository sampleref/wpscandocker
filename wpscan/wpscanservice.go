package wpscan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const wpscan = "/usr/bin/wpscan"

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var reports = make(map[string]ReportFile)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type CheckUrl struct {
	Url string `json:"Url"`
	Id  string `json:"Id"`
}

type ReportFile struct {
	Id       string `json:"Id"`
	Created  int64  `json:"Created"`
	FilePath string `json:"FilePath"`
	Url      string `json:"Url"`
}

type ReportDetails struct {
	Id     string `json:"Id"`
	Report string `json:"Report"`
	Url    string `json:"Url"`
	Status string `json:"Status"` //Possible: NOT_FOUND | FOUND | ERROR
}

func readCommandOutput(command string) string {
	log.Printf("Running command %s \n", command)
	c, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		fmt.Printf("Error while running command %v \n", err)
		return ""
	}
	log.Printf("Returning output as %s \n", c)
	return string(c)
}

func UpdateDatabase(w http.ResponseWriter, r *http.Request) {
	//Update database
	command := wpscan + " --update "
	output := readCommandOutput(command)
	fmt.Printf("Update database completed with output %s \n", output)
}

func CheckSingleUrl(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var checkUrl CheckUrl
	err := json.Unmarshal(reqBody, &checkUrl)
	if err != nil {
		fmt.Printf("Error while parsing request json %v \n", err)
		return
	}
	//Create Report File Instance
	Id := RandStringRunes(8)
	checkUrl.Id = Id
	randFilePath := "/" + Id + ".json"
	reportFile := ReportFile{
		Id:       Id,
		Created:  time.Now().Unix(),
		FilePath: randFilePath,
		Url:      checkUrl.Url,
	}
	reports[Id] = reportFile

	//Async check URL
	go func(url string, filePath string) {
		fmt.Printf("Async Check URL Started for URL %s at file %s \n", url, filePath)
		//Perform check
		command := wpscan + " -e ap --ignore-main-redirect --url " + checkUrl.Url + " --output=" + randFilePath + " --format json"
		output := readCommandOutput(command)
		fmt.Printf("Async Check URL Completed with output %s for url %s \n", output, url)
	}(checkUrl.Url, randFilePath)

	err = json.NewEncoder(w).Encode(checkUrl)
	if err != nil {
		fmt.Printf("Error while parsing response json %v \n", err)
		return
	}
}

func GetAllReports(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting all available reports \n")
	err := json.NewEncoder(w).Encode(reports)
	if err != nil {
		fmt.Printf("Error while parsing response json %v \n", err)
		return
	}
}

func DeleteReport(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("Deleting report \n")
	var reportDet ReportDetails
	err := json.Unmarshal(reqBody, &reportDet)
	if err != nil {
		fmt.Printf("Error while parsing request json %v \n", err)
		return
	}
	if val, ok := reports[reportDet.Id]; ok {
		err := os.Remove(val.FilePath)
		if err != nil {
			fmt.Printf("Error removing in file json at %s \n", val.FilePath)
		}
		delete(reports, val.Id)
	} else {
		fmt.Printf("Report not found to delete for id %s \n", reportDet.Id)
	}
}

func GetReportById(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var reportDet ReportDetails
	err := json.Unmarshal(reqBody, &reportDet)
	if err != nil {
		fmt.Printf("Error while parsing request json %v \n", err)
		return
	}
	if val, ok := reports[reportDet.Id]; ok {
		bytes, err := ioutil.ReadFile(val.FilePath)
		if err != nil {
			fmt.Printf("Error in reading file for id %s ", reportDet.Id)
			reportDetails := ReportDetails{
				Id:     val.Id,
				Report: "",
				Url:    val.Url,
				Status: "ERROR",
			}
			err = json.NewEncoder(w).Encode(reportDetails)
			if err != nil {
				fmt.Printf("Error while parsing response json %v \n", err)
			}
			return
		}
		fmt.Printf("Reading success for file for id %s ", reportDet.Id)
		reportDetails := ReportDetails{
			Id:     val.Id,
			Report: string(bytes),
			Url:    val.Url,
			Status: "FOUND",
		}
		err = json.NewEncoder(w).Encode(reportDetails)
		if err != nil {
			fmt.Printf("Error while parsing response json %v \n", err)
			return
		}
	} else {
		fmt.Printf("File not found for id %s ", reportDet.Id)
		reportDetails := ReportDetails{
			Id:     val.Id,
			Report: "",
			Url:    val.Url,
			Status: "NOT_FOUND",
		}
		err = json.NewEncoder(w).Encode(reportDetails)
		if err != nil {
			fmt.Printf("Error while parsing response json %v \n", err)
			return
		}
	}
}
