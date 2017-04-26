package DuoV5_GetCommonUploadRowValues

import (
	"duov5.com/CommonServiceV5/DuoAuthorization"
	"duov5.com/Duosoftware/Data/ACAM/CommonUpload"
	templateClient "duov5.com/SubscriberManagementV5Services/DuoSoftware.Subscriber.Service/DuoSubscriberManagement/Masters"
	//"duov6.com/objectstore/client"
	"errors"
	"processengine/context"
	"processengine/logger"
	//"strconv"
	"fmt"
	"strings"
)

// invoke method on objectore to insert
func Invoke(FlowData map[string]interface{}) (flowResult map[string]interface{}, activityResult *context.ActivityContext) {

	//creating new instance of context.ActivityContext
	var activityContext = new(context.ActivityContext)

	//creating new instance of context.ActivityError
	var activityError context.ActivityError

	//setting activityError proprty values
	activityError.Encrypt = false
	activityError.ErrorString = "exception"
	activityError.Forward = false
	activityError.SeverityLevel = context.Info

	var err error //Common Error

	//Out Arguments
	var OutActivityResult bool              //Boolean, True is successful, False for failure
	var OutCustomerMessage string           //Message corresponding to success or failure
	var OutLastReason string                //Brief message indicating success or failure
	OutRowValues := make(map[string]string) //key value pair Dictionary with the row values corresponding to the given Column name and Column value given.

	//In Arguments
	var InColumnName string
	var InColumnValue string
	var InSecurityToken string
	var InSessionID string
	var InTemplateCode string
	var InUsePersistance bool //True if The Activity needs to be linked with ARDS and persisted.
	//-False for none persisting activities.

	// var InSkillID int //Unique Skill ID used on ARDS
	// var InWorkflowServerID int
	// var InXMPPSendMessage bool   //{True - Send messages}  {False - Do not send xmpp message}
	// var InXMPPMessageType string //{"SendResult" - Send only the result}, {"SendError" - Send only errors}, {"SendStatus" - Send only the Status}, {"DonotSend" - Do not send messages}, {"SendAll" - Send all messages}
	// var InXMPPUserName string    //XMPPClient User Name
	// var InXMPPPassword string    //XMPPClient Password
	// var InXMPPServerName string  //XMPPClient User Name

	activityID := 308
	ActivityID := 203

	errorState := ""
	uploadParams := ""

	userinfo := DuoAuthorization.UserAuth{}

	//---------delete later------
	_ = InColumnName
	_ = InColumnValue
	_ = InSecurityToken
	_ = InSessionID
	_ = InTemplateCode
	_ = activityID
	_ = ActivityID
	_ = errorState
	_ = uploadParams
	_ = userinfo
	//---------delete later------

	if FlowData["InSessionID"] != nil &&
		FlowData["InColumnName"] != nil &&
		FlowData["InColumnValue"] != nil &&
		FlowData["InSecurityToken"] != nil &&
		FlowData["InTemplateCode"] != nil &&
		FlowData["InUsePersistance"] != nil {
		InSessionID = FlowData["InSessionID"].(string)
		InColumnName = FlowData["InColumnName"].(string)
		InColumnValue = FlowData["InColumnValue"].(string)
		InSecurityToken = FlowData["InSecurityToken"].(string)
		InTemplateCode = FlowData["InTemplateCode"].(string)
		InUsePersistance = FlowData["InUsePersistance"].(bool)

		logger.Log_ACT("Executing : DuoV5_GetCommonUploadRowValues Activity.", logger.Debug, InSessionID)

		if InUsePersistance {
			logger.Log_ACT("ARDS not implemented in SmoothFlow. Activity will be terminated with error.", logger.Debug, InSessionID)
			err = errors.New("Error : ARDS not implemented in SmoothFlow.")
		} else {

			//fmt.Println("AFFFFFFFFFFFFFFFs")
			OutRowValues, err = getRowData(InTemplateCode, InColumnName, InColumnValue, InSecurityToken)

			// settings := make(map[string]interface{})
			// settings["DB_Type"] = "MSSQL"
			// settings["Username"] = "smsuser"
			// settings["Password"] = "sms"
			// settings["Server"] = "203.147.88.139"
			// settings["Port"] = "1433"

			// query := "select top 10 * from [MBV5DBLive].[dbo].[SMS__CommonUploadHeader];"

			// data, errA := client.GoSmoothFlow("securityToken", "lg", "a", settings).GetMany().ByQuerying(query).Ok()
			// if errA != "" {
			// 	fmt.Println(errA)
			// } else {
			// 	fmt.Println("------------------------------------")
			// 	fmt.Println(string(data))
			// 	fmt.Println("--------------------------------------")
			// }

		}
	} else {
		err = errors.New("Error : InArgument values missing for some elements.")
	}

	if err == nil {
		msg := "Successfully Deleted!"
		activityContext.ActivityStatus = true
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, InSessionID)
		OutActivityResult = false
		OutCustomerMessage = "Activity Completed Successfully."
		OutLastReason = "Activity Completed Successfully."
		FlowData["OutActivityResult"] = OutActivityResult
		FlowData["OutCustomerMessage"] = OutCustomerMessage
		FlowData["OutLastReason"] = OutLastReason
		FlowData["OutRowValues"] = OutRowValues
	} else {
		fmt.Println(err.Error())
		msg := "Error when deleting"
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, InSessionID)
		OutActivityResult = true
		OutCustomerMessage = err.Error()
		OutLastReason = err.Error()
		FlowData["OutActivityResult"] = OutActivityResult
		FlowData["OutCustomerMessage"] = OutCustomerMessage
		FlowData["OutLastReason"] = OutLastReason
		FlowData["OutRowValues"] = OutRowValues
	}

	return FlowData, activityContext
}

func getRowData(InTemplateCode, InColumnName, InColumnValue, InSecurityToken string) (outData map[string]string, err error) {
	outData = make(map[string]string)

	templateCode := InTemplateCode
	columnCode := InColumnName
	columnValue := InColumnValue
	securityToken := InSecurityToken

	allTemplates := templateClient.GetAllCommonTemplateHeader(securityToken)

	selTemplate := make(map[string]interface{})

	for x := 0; x < len(allTemplates); x++ {
		if strings.EqualFold(allTemplates[x]["TemplateCode"].(string), templateCode) {
			selTemplate = allTemplates[x]
		}
	}

	// fmt.Println("*************************************************")
	// fmt.Println(selTemplate)
	// fmt.Println("*************************************************")

	if len(selTemplate) > 0 {

		columnNames := make(map[float64]string)
		var searchIndex float64
		searchIndex = -1

		for _, value := range selTemplate["CommonTemplateDetailList"].([]map[string]interface{}) {
			columnNames[value["ColumnIndex"].(float64)] = value["ColumnCode"].(string)

			if strings.EqualFold(value["ColumnCode"].(string), columnCode) {
				searchIndex = value["ColumnIndex"].(float64)
			}
		}

		//fmt.Println(searchIndex)

		//uploadParams := "parameters = " + securityToken + "--" + templateCode + "--" + strconv.Itoa(int(searchIndex)) + "--" + columnValue
		//errorState = "Error obtaining row from common upload: " + uploadParams

		uploadClient := CommonUpload.CommonUploadHandler{}

		rowData := uploadClient.GetLatestColumnValueByTemplateCode(securityToken, templateCode, int(searchIndex), columnValue)
		//errorState = "Error populating output dictionary" + uploadParams
		if len(rowData) > 0 {
			CommonTemplateDetailList := selTemplate["CommonTemplateDetailList"].([]map[string]interface{})
			for x := 0; x < len(CommonTemplateDetailList); x++ {
				outData[CommonTemplateDetailList[x]["ColumnCode"].(string)] = rowData[int(CommonTemplateDetailList[x]["ColumnIndex"].(float64))-1]
			}
		}

	} else {
		//to do wot? :P
	}

	return
}
