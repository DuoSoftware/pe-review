package objectstore

import (
	"bytes"
	"net/http"
	"processengine/context"
)

type Delete struct {
}

func (d Delete) Invoke(parameters map[string]interface{}) *context.ActivityContext {

	securityToken := parameters["securityToken"].(string)
	/*log := parameters["log"].(string)*/
	domain := parameters["namespace"].(string)
	class := parameters["class"].(string)
	JSON_Document := parameters["JSON"].(string)

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

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer([]byte(JSON_Document)))
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
		activityContext.ActivityStatus = true
		activityContext.Message = "Data Successfully Deleted!"
		activityContext.ErrorState = activityError
	}
	defer resp.Body.Close()
	return activityContext
}
