package main

import (
	"bytes"
	"duov5.com/DuoAuthorization"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	dd := DuoAuthorization.DuoAuthorization{}
	accessCode := dd.GetAccessCode(10004, 1000, "Write")
	fmt.Println(dd.GetAccess("3ccc98d7b4fa61d3a3fa162f95f79316", accessCode))
}

func dd(query string) {

	url := "http://192.168.1.194/DuoSubscribe5/CommonServices/Authorization/auth.svc"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(query)))
	req.Method = "POST"

	if err != nil {
		fmt.Println("E1 : " + err.Error())
		return
	}

	req.Header.Set("Content-Type", "text/xml")

	//Login
	//req.Header.Set("SOAPAction", "http://tempuri.org/Iauth/login")
	//GetAccess
	req.Header.Set("SOAPAction", "http://tempuri.org/Iauth/GetAccess")
	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("E2 : " + err.Error())
		return
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("E3 : " + err.Error())
		return
	}

	fmt.Println(string(b))
}
