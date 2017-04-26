package main

import (
	"bytes"
	"fmt"
	"github.com/clbanning/x2j"
	"io/ioutil"
	"net/http"
	"reflect"
)

func main() {
	query := `<x:Envelope xmlns:x="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tem="http://tempuri.org/">
    <x:Header/>
    <x:Body>
        <tem:GetAccess>
            <tem:SecurityToken>37c146758bfa95df13fc8d2544f39962</tem:SecurityToken>
            <tem:AccessCode></tem:AccessCode>
        </tem:GetAccess>
    </x:Body>
</x:Envelope>`
	GetSoapEnvelope(query)
}

const url = "http://192.168.1.194/DuoSubscribe5/CommonServices/Authorization/auth.svc"

func GetSoapEnvelope(query string) {
	httpClient := new(http.Client)
	resp, err := httpClient.Post(url, "application/soap+xml", bytes.NewBufferString(query))
	if err != nil {
		fmt.Println(err.Error())
	}
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		fmt.Println(e.Error())

	} else {

		data := make(map[string]interface{})

		doc, err4 := x2j.DocToMap(string(b))
		if err4 != nil {
			fmt.Println(err4.Error())
		} else {
			// for key, value := range doc["Envelope"].(map[string]interface{}) {
			// 	if key == "Body" {
			// 		for k1, v1 := range value.(map[string]interface{}) {
			// 			if k1 == "GetAccountInfoByGuAccountIDResponse" {
			// 				for k2, v2 := range v1.(map[string]interface{}) {
			// 					if k2 == "GetAccountInfoByGuAccountIDResult" {
			// 						for k3, v3 := range v2.(map[string]interface{}) {
			// 							if reflect.TypeOf(v3).String() == "map[string]interface {}" && v3.(map[string]interface{})["-nil"] != nil {
			// 								data[k3] = ""
			// 							} else {
			// 								data[k3] = v3
			// 							}

			// 						}
			// 					}
			// 				}
			// 			}
			// 		}
			// 	}
			// }
			fmt.Println(data)
		}

	}

	resp.Body.Close()
}
