package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

type Loan struct {
	LoanId       string  `json:"id"`
	Customername string  `json:"customername"`
	Phoneno      string  `json:"phoneno"`
	Email        string  `json:"email"`
	LoanAmount   float64 `json:"loanamount"`
	Status       string  `json:"status"`
	CreditScore  int     `json:"creditscore"`
}

var status_list []string = []string{"New", "Approved", "Rejected", "Cancelled"}

func GenerateID() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, 10)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func (loan Loan) IsValidStatus(data string) bool {
	for _, i := range status_list {
		if i == data {
			return true
		}
	}
	return false

}

func (loan Loan) CreateLoan(response http.ResponseWriter, request *http.Request) {
	request_data, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(request_data, &loan)

	if loan.Customername == "" || loan.Phoneno == "" || loan.Email == "" || loan.CreditScore == 0 || loan.LoanAmount == 0 {
		http.Error(response, "Bad Request", http.StatusBadRequest)
	}
	loan.LoanId = GenerateID()
	loan.Status = "New"
	_, err := collection.InsertOne(context.TODO(), loan)
	if err != nil {
		http.Error(response, "Loan Creation Failed", http.StatusNotAcceptable)

	}
	response.WriteHeader(http.StatusCreated)
}

func (loan Loan) ChangeLoanStatus(response http.ResponseWriter, request *http.Request) {
	var statusdata map[string]string
	request_data, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(request_data, &statusdata)
	vars := mux.Vars(request)
	filter := bson.M{"loanid": vars["id"]}

	result := loan.IsValidStatus(statusdata["status"])
	if !result {
		http.Error(response, "Bad Request", http.StatusBadRequest)
	}
	staus_update := bson.M{"status": statusdata["status"]}
	update := bson.M{"$set": staus_update}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(response, "Status Change Failed", http.StatusNotModified)
	}
	response.WriteHeader(http.StatusOK)
}

func (loan Loan) CancelLoan(response http.ResponseWriter, request *http.Request) {
	var statusdata map[string]string
	request_data, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(request_data, &statusdata)
	vars := mux.Vars(request)
	filter := bson.M{"loanid": vars["id"]}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		http.Error(response, "Not Modified", http.StatusNotModified)
	}
}

func (loan Loan) GetLoanByID(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	filter := bson.M{"loanid": vars["id"]}
	collection.FindOne(context.TODO(), filter).Decode(&loan)
	js, _ := json.Marshal(loan)
	response.Header().Set("Content-Type", "application/json")
	response.Write(js)

}

func ListLoans(response http.ResponseWriter, request *http.Request) {

	var loandata []Loan
	var filter bson.M

	query_status := request.URL.Query().Get("status")
	query_AmountGreater, _ := strconv.Atoi(request.URL.Query().Get("greaterthan"))

	if query_status == "" && query_AmountGreater == 0 {
		filter = bson.M{}

	} else if query_status == "" && query_AmountGreater != 0 {
		greater := bson.M{"$gt": query_AmountGreater}
		filter = bson.M{"loanamount": greater}

	} else if query_status != "" && query_AmountGreater == 0 {
		filter = bson.M{"status": query_status}
	} else {
		greater := bson.M{"$gt": query_AmountGreater}
		filter = bson.M{"status": query_status, "loanamount": greater}

	}
	list_cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		http.Error(response, "Unable to Fetch Data", http.StatusNotFound)
	}
	err = list_cursor.All(ctx, &loandata)
	if err != nil || len(loandata) <= 0 {
		http.Error(response, "Unable to Fetch Data", http.StatusNotFound)
	}
	js, _ := json.Marshal(&loandata)
	response.Header().Set("Content-Type", "application/json")
	response.Write(js)
}
