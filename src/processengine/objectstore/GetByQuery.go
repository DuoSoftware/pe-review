package objectstore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"processengine/context"
)

type GetByQuery struct {
}

func (g GetByQuery) Invoke(parameters map[string]interface{}) *context.ActivityContext {

	securityToken := parameters["securityToken"].(string)
	/*log := parameters["log"].(string)*/
	domain := parameters["namespace"].(string)
	class := parameters["class"].(string)
	JSON_Document := parameters["query"].(string)
	queryString := "{\"Query\":{\"Type\":\"Query\", \"Parameters\":\"" + JSON_Document + " \"}}"

	//creating new instance of context.ActivityContext
	var activityContext = new(context.ActivityContext)

	//creating new instance of context.ActivityError
	var activityError context.ActivityError

	//setting activityError proprty values
	activityError.Encrypt = false
	activityError.ErrorString = "No Exception"
	activityError.Forward = false
	activityError.SeverityLevel = context.Info

	//url := "http://" + Common.GetObjectstoreIP() + "/" + domain + "/" + class
	url := "http://" + domain + "/data/" + domain + "/" + class + "?securityToken=" + securityToken
	fmt.Println(url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(queryString)))

	/*req.Header.Set("securityToken", securityToken)
	req.Header.Set("log", log)*/

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		activityError.ErrorString = err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = "Connection to server failed! Check connectivity."
		activityContext.ErrorState = activityError
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		/*var objArray []map[string]interface{}
		json.Unmarshal(body, &objArray)*/

		var response = new(context.ObjectStoreResponse)
		json.Unmarshal(body, &response)

		if len(response.Result) > 0 {
			activityContext.ActivityStatus = true
			activityContext.Message = "Data Successfully Retireved!"
			activityContext.ErrorState = activityError

			obj, err := json.Marshal(response.Result)
			if err != nil {
				fmt.Println(err.Error())
			}
			activityContext.SharedContext = obj
		} else {
			activityContext.ActivityStatus = false
			activityContext.Message = "Records Not Found!"
			activityContext.ErrorState = activityError
		}
	}
	defer resp.Body.Close()

	return activityContext
}
