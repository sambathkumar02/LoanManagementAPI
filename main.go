package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var username string = os.Getenv("DB_USERNAME")
var pass string = os.Getenv("DB_PASS")

//set for connection URI
var client, err = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://" + username + ":" + pass + "@cluster0.ndaoh.mongodb.net/LoanManagement"))

//set the timeout limit
var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

var collection = client.Database("LoanManagement").Collection("LoanDetails")
var r = mux.NewRouter()

func main() {

	if err != nil {
		fmt.Print("Error Connecting Database!")

	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//Automatically close connection When exits
	defer client.Disconnect(ctx)

	//Get the List of all Databases
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Print("Database access sucess!")

	}

	loan := Loan{}

	r.HandleFunc("/loans", ListLoans).Methods("GET")
	r.HandleFunc("/loans", loan.CreateLoan).Methods("POST")
	r.HandleFunc("/loans/{id}", loan.GetLoanByID).Methods("GET")
	r.HandleFunc("/loans/{id}", loan.ChangeLoanStatus).Methods("PATCH")
	r.HandleFunc("/loans/{id}", loan.CancelLoan).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8088", r))

}
