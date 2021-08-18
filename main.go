package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// encode the data
	postBody, _ := json.Marshal(map[string]string{
		"email":    "foobar@gmail.com",
		"password": "1234",
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://localhost:3000/authenticate", "application/json", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	// client must close the response body when finished
	defer resp.Body.Close()

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300

	// read the response body
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	} else if !statusOK {
		log.Println("Whoops! Authentication token was not granted [Status Code", resp.StatusCode, "].")
	} else {
		var result map[string]string
		json.Unmarshal(body, &result)

		authToken := string(result["auth_token"])
		log.Println("Success! Your authentication token is", authToken)
	}
}
