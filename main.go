package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type JsonResp struct {
	Success Success `json:"quote"`
	Error   Error   `json:"error"`
}

type Success struct {
	Ticker       string   `json:"ticker"`
	TickerColor  string   `json:"ticker_color"`
	CompanyName  string   `json:"company_name"`
	OpenPrice    float32  `json:"open_price"`
	Delta        float32  `json:"delta"`
	CurrentPrice float32  `json:"current_price"`
	Tags         []string `json:"tags"`
}

type Error struct {
	Title  string `json:"title"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func main() {
	stocks := os.Args[1:]

	if len(stocks) == 0 {
		log.Fatalln("Please enter a ticker symbol (e.g. AAPL, AMZN).")
	}

	symbol := stocks[0]

	authToken, err := getAuthToken()
	if err != nil {
		log.Fatalln(err)
	}

	host := goDotEnvVariable("RG_DASHBOARD_API_HOST")

	client := &http.Client{}
	req, _ := http.NewRequest("GET", host+"/api/v1/quotes/"+symbol, nil)
	req.Header.Add("Authorization", authToken)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// parse response body into data object
	var jsonResp JsonResp
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode == 200 {
		log.Println("Request for", symbol, "quote was successful!")
		j, _ := json.MarshalIndent(jsonResp, "", "    ")
		fmt.Println(string(j))
	} else {
		log.Println("Error:", jsonResp.Error.Detail, "[Status Code "+jsonResp.Error.Code+"].")
	}
}

func getAuthToken() (string, error) {
	email := goDotEnvVariable("RG_DASHBOARD_EMAIL")
	pswd := goDotEnvVariable("RG_DASHBOARD_PASSWORD")
	host := goDotEnvVariable("RG_DASHBOARD_API_HOST")

	// encode the data
	postBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": pswd,
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(host+"/authenticate", "application/json", responseBody)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	// client must close the response body when finished
	defer resp.Body.Close()

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300

	// read the response body
	body, err := ioutil.ReadAll(resp.Body)
	var authToken string

	if err != nil {
		log.Fatalln(err)
	} else if !statusOK {
		return "", errors.New("Whoops! Authentication token was not granted [Status Code " + fmt.Sprintf("%v", resp.StatusCode) + "].")
	} else {
		var result map[string]string
		json.Unmarshal(body, &result)
		authToken = string(result["auth_token"])
		log.Println("Your authentication token is", authToken)
	}

	return authToken, nil
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
