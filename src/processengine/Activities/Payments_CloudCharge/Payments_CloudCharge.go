package Payments_CloudCharge

import "processengine/context"
import "net/http"
import "bytes"
import "processengine/logger"
import "io/ioutil"
import "encoding/json"

// the struct relavent for the object accepted through the API
type PayItem struct {
	ItemRefID   string `json:"ItemRefID"`
	ItemType    string `json:"ItemType"`
	Description string `json:"Description"`
	UnitPrice   string `json:"UnitPrice"`
	UOM         string `json:"UOM"`
	Qty         string `json:"Qty"`
	Subtotal    string `json:"Subtotal"`
	Discount    string `json:"Discount"`
	Tax         string `json:"Tax"`
	TotalPrice  string `json:"TotalPrice"`
	AccountId   string `json:"AccountId"`
	TaxDetails  string `json:"TaxDetails"`
}

// the main method for cloud charge method
func Invoke(FlowData map[string]interface{}) (flowResult map[string]interface{}, activityResult *context.ActivityContext) {

	logger.Log_ACT("Starting cloudcharge request.", logger.Debug, FlowData["InSessionID"].(string))
	//creating new instance of context.ActivityContext
	var activityContext = new(context.ActivityContext)

	//creating new instance of context.ActivityError
	var activityError context.ActivityError

	//setting activityError proprty values
	activityError.Encrypt = false
	activityError.ErrorString = "exception"
	activityError.Forward = false
	activityError.SeverityLevel = context.Info

	// getting the details from the input arguments
	Namespace := FlowData["Namespace"].(string)
	AccountId := FlowData["AccountId"].(string)
	CardNo := FlowData["CardNo"].(string)
	Name := FlowData["Name"].(string)
	CardType := FlowData["CardType"].(string)
	DeliveryAddress := FlowData["DeliveryAddress"].(string)
	BillingAddress := FlowData["BillingAddress"].(string)
	CSV := FlowData["CSV"].(string)
	ExpiryYear := FlowData["ExpiryYear"].(string)
	ExpiryMonth := FlowData["ExpiryMonth"].(string)
	Active := FlowData["Active"].(string)
	Verified := FlowData["Verified"].(string)

	// preparing the list of items to be passed to the API
	var ItemList []PayItem
	if jsonParseErr := json.Unmarshal([]byte(FlowData["ItemArray"].(string)), &ItemList); jsonParseErr != nil {
		logger.Log_ACT("The JSON input is not in correct format.", logger.Debug, FlowData["InSessionID"].(string))
	}
	// convert the struct to JSON
	jsonOut, MarshalError := json.Marshal(ItemList)
	if MarshalError != nil {
		logger.Log_ACT("marshal error", logger.Debug, FlowData["InSessionID"].(string))
	}

	// preparing the request object
	requestStr := `
	{
		"AccountId": "` + AccountId + `",
		"Cards": 
		{
			"CardNo": "` + CardNo + `",
			"Name": "` + Name + `",
			"CardType": "` + CardType + `",
			"DeliveyAddress": "` + DeliveryAddress + `",
			"BillingAddress": "` + BillingAddress + `",
			"CSV": "` + CSV + `",
			"ExpiryYear": "` + ExpiryYear + `",
			"ExpiryMonth": "` + ExpiryMonth + `",
			"Active": ` + Active + `,
			"verified": ` + Verified + `
		}
		,
		"Items": ` + string(jsonOut) + `
		}`

	// preparing the URL for the API
	url := "http://" + Namespace + "/payapi/transaction/paystrip"

	// go the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestStr)))

	client := &http.Client{}
	resp, err := client.Do(req)
	// once the request if completed the following is done according to the response
	if err != nil {
		logger.Log_ACT("There was an error while making the HTTP request for Cloudcharge.", logger.Debug, FlowData["InSessionID"].(string))
		//activityError.ErrorString = err.Error().(string)
		activityContext.ActivityStatus = false
		activityContext.Message = "Connection to server failed! Check connectivity."
		activityContext.ErrorState = activityError
	} else {
		logger.Log_ACT("Cloudcharge request was successfull.", logger.Debug, FlowData["InSessionID"].(string))
		activityContext.ActivityStatus = true
		activityContext.Message = "Cloudcharge request was successfull."
		activityContext.ErrorState = activityError
	}

	// get the content returned from the API request
	bs, bserror := ioutil.ReadAll(resp.Body)
	if bserror != nil {
		msg := "There was an error reading the responce body."
		logger.Log_ACT(msg, logger.Debug, FlowData["InSessionID"].(string))
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityError.ErrorString = bserror.Error()
		activityContext.ErrorState = activityError
	}
	// the response is added to the return MAP object
	FlowData["Response"] = string(bs)

	// close the body once the request is completed.
	defer resp.Body.Close()
	return FlowData, activityContext
}
