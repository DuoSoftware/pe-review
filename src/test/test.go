package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	hmm = 11
)

func main() {
	err, body := HTTP_GET("https://duosoftware.atlassian.net/rest/api/2/myself")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(body))
	}
}

func HTTP_GET(url string) (err error, body []byte) {

	req, err := http.NewRequest("GET", url, nil)
	c1 := &http.Cookie{Name: "studio.crowd.tokenkey", Value: "Q40VYwmQUVC2n2gPw73YKw00", HttpOnly: false}
	c2 := &http.Cookie{Name: "JSESSIONID", Value: "E4C3BA7F11A93EE1CE354DBF5991D8F0", HttpOnly: false}
	//c3 := &http.Cookie{Name: "atlassian.xsrf.token", Value: "B04X-VUKS-X5QS-YFE2|7c26d8ab1d0decd37bbbcf0777f35c0a2324ace8|lin; cloud.session.token=eyJraWQiOiJzZXNzaW9uLXNlcnZpY2VcL3Nlc3Npb24tc2VydmljZSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI1NTcwNTg6ZjRkOWYxMzYtZDI0Yy00ZWRjLWE1ODYtM2M2NjM5ZjE2YTI1IiwiYXVkIjoiYXRsYXNzaWFuIiwiaW1wZXJzb25hdGlvbiI6W10sIm5iZiI6MTQ5MDAwNTcxNCwicmVmcmVzaFRpbWVvdXQiOjE0OTAwMDYzMTQsImlzcyI6InNlc3Npb24tc2VydmljZSIsInNlc3Npb25JZCI6ImMyOGE2NzdiLTFhMTYtNDg5YS1hMzJjLTFjZmMwYzE3MGU2MSIsImV4cCI6MTQ5MjU5NzcxNCwiaWF0IjoxNDkwMDA1NzE0LCJlbWFpbCI6InByYXNhZGpheWFzaGFua2FAZ21haWwuY29tIiwianRpIjoiYzI4YTY3N2ItMWExNi00ODlhLWEzMmMtMWNmYzBjMTcwZTYxIn0.xjVI6EKTlrfgKWqaVNAAfhtqeifVVH-wngrkK_KkxjzIq5bTgPESu4rSoa_HsbPPizDfB3_drjAb_yki-6jbyxJMFgismghaAvsOPc5GRPYkmrzoibE3UhcmQcHivcgCYIMcOPfpDz-7vMv1wsKz4KXrAmP876h9BZSA6JvGSclXZaK5D6Jb8VPvjMueskYTyxSkfAjDm3tdgp6Y35IJLLzI1_z4tdnIeS_G201J9Uw0uBWtOIFsEAB9i0YChwkfR-udzvZ8llQSGEncs-A8F-18l10FQuwz9hRVmLBjhlAwznZ-3WnMe9m0fphWxJpynM4h6FfaWGO9i8Hdv9lCIg; JSESSIONID=E4C3BA7F11A93EE1CE354DBF5991D8F0; __utma=174252215.298312510.1488184114.1491192622.1491300651.55; __utmb=174252215.1.10.1491300651; __utmc=174252215; __utmz=174252215.1488184114.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); studio.crowd.tokenkey=Q40VYwmQUVC2n2gPw73YKw00", HttpOnly: false}
	//c4 := &http.Cookie{Name: "cloud.session.token", Value: "eyJraWQiOiJzZXNzaW9uLXNlcnZpY2VcL3Nlc3Npb24tc2VydmljZSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI1NTcwNTg6ZjRkOWYxMzYtZDI0Yy00ZWRjLWE1ODYtM2M2NjM5ZjE2YTI1IiwiYXVkIjoiYXRsYXNzaWFuIiwiaW1wZXJzb25hdGlvbiI6W10sIm5iZiI6MTQ5MDAwNTcxNCwicmVmcmVzaFRpbWVvdXQiOjE0OTAwMDYzMTQsImlzcyI6InNlc3Npb24tc2VydmljZSIsInNlc3Npb25JZCI6ImMyOGE2NzdiLTFhMTYtNDg5YS1hMzJjLTFjZmMwYzE3MGU2MSIsImV4cCI6MTQ5MjU5NzcxNCwiaWF0IjoxNDkwMDA1NzE0LCJlbWFpbCI6InByYXNhZGpheWFzaGFua2FAZ21haWwuY29tIiwianRpIjoiYzI4YTY3N2ItMWExNi00ODlhLWEzMmMtMWNmYzBjMTcwZTYxIn0.xjVI6EKTlrfgKWqaVNAAfhtqeifVVH-wngrkK_KkxjzIq5bTgPESu4rSoa_HsbPPizDfB3_drjAb_yki-6jbyxJMFgismghaAvsOPc5GRPYkmrzoibE3UhcmQcHivcgCYIMcOPfpDz-7vMv1wsKz4KXrAmP876h9BZSA6JvGSclXZaK5D6Jb8VPvjMueskYTyxSkfAjDm3tdgp6Y35IJLLzI1_z4tdnIeS_G201J9Uw0uBWtOIFsEAB9i0YChwkfR-udzvZ8llQSGEncs-A8F-18l10FQuwz9hRVmLBjhlAwznZ-3WnMe9m0fphWxJpynM4h6FfaWGO9i8Hdv9lCIg", HttpOnly: false}
	req.AddCookie(c1)
	req.AddCookie(c2)
	//req.AddCookie(c3)
	//req.AddCookie(c4)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		err = errors.New("Connection Failed!")
	} else {
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			err = errors.New(string(body))
		}
	}

	return
}
