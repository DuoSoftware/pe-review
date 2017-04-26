package objectstore

import (
	"bytes"
	"net/http"
	"processengine/context"
)

type Update struct {
}

// method used to Update data present in objectstore
func (u Update) Invoke(parameters map[string]interface{}) *context.ActivityContext {

	// getting and assign the values from the MAP to internal variables
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

	// preparing the data relevent to make the objectstore API
	//url := "http://" + Common.GetObjectstoreIP() + "/" + domain + "/" + class
	url := "http://" + domain + "/data/" + domain + "/" + class + "?securityToken=" + securityToken

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(JSON_Document)))
	/*req.Header.Set("securityToken", securityToken)
	req.Header.Set("log", log)*/
	client := &http.Client{}
	resp, err := client.Do(req)
	// once the request is made according to the response the following is done
	if err != nil {
		activityError.ErrorString = err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = "Connection to server failed! Check connectivity."
		activityContext.ErrorState = activityError
	} else {
		activityContext.ActivityStatus = true
		activityContext.Message = "Data Successfully Updated!"
		activityContext.ErrorState = activityError
	}
	defer resp.Body.Close()
	// once the functionality of the method completes it returns the processed data
	return activityContext
}
