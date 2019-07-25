package main

import (
	"log"
	"fmt"
	"net/http"
	"encoding/json"
 
    
)


type Country struct{
    Name        string      `json:"name"`   
    Capital     string      `json:"capital"` 
}
 

type AddressResponse struct{
	FinalBalance      int64   `json:"final_balance"`  
	TotalReceived     int64   `json:"total_received"`
}

func (adRes *AddressResponse) ToString() string {
   return fmt.Sprintf("{ 'final_balance': %d, 'total_received': %d }", adRes.FinalBalance, adRes.TotalReceived)
}


func Call(address string) bool{

	client := &http.Client{}
    url := "https://blockchain.info/rawaddr/" + address

    // build a new request, but not doing the POST yet
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
    	log.Println(err)
      return false
    }
    req.Header.Add("Content-Type", "application/json")

    // now POST it
    resp, errResp := client.Do(req)
    if errResp != nil  || resp.StatusCode != http.StatusOK {
        log.Printf("Status error: %v", resp.StatusCode)
        return false
    }


    var addressResponse AddressResponse
    if errOrder := json.NewDecoder(resp.Body).Decode(&addressResponse); errOrder != nil {
        log.Printf("Parse error: %v", errOrder)
        return false
    }

    //log.Println(addressResponse.ToString())
   // return (addressResponse.FinalBalance > 0 || addressResponse.TotalReceived > 0)
    return (addressResponse.FinalBalance > 0)

}
 
func getCountries() []Country{
    list := make([]Country, 0)

    client := &http.Client{}
    url := "https://restcountries.eu/rest/v2/all"

    // build a new request, but not doing the POST yet
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Println(err)
    }
    req.Header.Add("Content-Type", "application/json")

    // now POST it
    resp, errResp := client.Do(req)
    if errResp != nil  || resp.StatusCode != http.StatusOK {
        log.Printf("Status error: %v", resp.StatusCode)
    }

    if errOrder := json.NewDecoder(resp.Body).Decode(&list); errOrder != nil {
        log.Printf("Parse error: %v", errOrder)
    }

    return list
}
