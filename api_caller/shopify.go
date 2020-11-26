package api_caller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ListOfCustomers struct{
	Name string
	Email string
}

func ShopifyCall() {
	// make a request
	query :=
		"query{" +
			"customers(first: 50){" +
			"edges{" +
			"node{" +
			"firstName\n" +
			"lastName\n" +
			"email\n" +
			"}" +
			"}" +
			"}" +
			"}"

	request, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(query)))
	//set auth headers
	request.Header.Set("Content-Type", "application/graphql")
	request.Header.Set("X-Shopify-Access-Token", "")
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	responseData, _ := ioutil.ReadAll(response.Body)

	// fmt.Println(string(responseData))

	type CustomerList struct {
		Data struct {
			Customers struct {
				Edges []struct {
					Node struct {
						FirstName string `json:"firstName"`
						LastName  string `json:"lastName"`
						Email     string `json:"email"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"customers"`
		} `json:"data"`
	}

	fmt.Println("#########")
	var customerlist CustomerList

	error := json.Unmarshal(responseData, &customerlist)
	if error != nil {
		fmt.Println("Failed to put into struct : \n", error)
	}
	var list [50]ListOfCustomers
	//Loops trough all customers in struct
	for i := range customerlist.Data.Customers.Edges {
		list[i] = ListOfCustomers{
			Name : customerlist.Data.Customers.Edges[i].Node.FirstName +
				" " + customerlist.Data.Customers.Edges[i].Node.LastName,
			Email: customerlist.Data.Customers.Edges[i].Node.Email,
		}
	}
	placeHolderFunction(list)
}


func placeHolderFunction(list [50]ListOfCustomers){

	for i := range list {
		fmt.Println(list[i])
	}
}