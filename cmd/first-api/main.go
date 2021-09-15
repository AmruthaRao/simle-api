package main

import (
	"io"
	"log"
	"fmt"
	"net/http"
	"strconv"
)

func main() {

	http.HandleFunc("/voting/eligibility", voteEligibilityHandler)
	log.Println("Listing for requests at http://localhost:8000/voting/eligibility")
	log.Fatal(http.ListenAndServe(":8000", nil))
}


func voteEligibilityHandler(w http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	name := queryParams["name"][0]
	age, _ :=  strconv.Atoi(queryParams["age"][0])
	validationResult := make(chan bool)
	message := make(chan string)
	go firstValidator(name, validationResult)
	go secondValidator(age, validationResult)
	go thirdValidator("blah", validationResult)
	x, y, z := <-validationResult, <-validationResult, <-validationResult
	log.Println(x, y, z)
	if x && y && z {
		io.WriteString(w, "All the validations passed\n")
	} else {
		io.WriteString(w, "One of the validations failed\n")
		return
	}
	go greet(name, message)
	go checkEligibility(age, message)
	msg2, msg1 := <-message, <-message
	go buildStatusReport(message, msg1, msg2)
	m := <- message
	io.WriteString(w, m)
}
func firstValidator(name string, c chan bool) {
	log.Println("inside validator 1")
	//should only have alphabets
	c <- len(name) < 200
}

func secondValidator(age int, c chan bool) {
	log.Println("inside validator 2")
	c <- age < 100
}

func thirdValidator(email string, c chan bool) {
	log.Println("inside validator 3")
	//check ends with .com
	c <- true
}

func greet(name string, c chan string) {
	c <- fmt.Sprintf("Hello %s\n", name)
}

func checkEligibility(age int, c chan string) {
	switch {
	case age > 17:
		c <- "You are eligible to vote\n"
	default:
		c <- fmt.Sprintf("You have to wait for %d more year/s to vote\n", 18-age)
	}
}

func buildStatusReport(c chan string, messages ...string) {
	header := "Voting Eligibility Status:\n"
	var body string
	for _, m := range messages {
		body += m
	}
	c <- fmt.Sprintf("%s%s", header, body)
}
