package DuoV5_DUOInventoryActivities_GetSerialInfoListBy

import (
	//"duov5.com/CommonServiceV5/DuoAuthorization"
	//"duov5.com/Duosoftware/Data/ACAM/CommonUpload"
	//templateClient "duov5.com/SubscriberManagementV5Services/DuoSoftware.Subscriber.Service/DuoSubscriberManagement/Masters"
	//"duov6.com/objectstore/client"
	"duov5.com/DuoSoftware/Subscriber/InventoryTransaction"
	//"duov5.com/DuoSoftware/Subscriber/SubscriberManagement"
	//"duov5.com/DuoSoftware/Subscriber/SubscriberMasters"
	"errors"
	"processengine/context"
	"processengine/logger"
	//"strconv"
	"fmt"
	//"strings"
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

	//ActivityID := 512

	//Out Arguments
	var OutActivityResult bool                             //Boolean, True is successful, False for failure
	var OutCustomerMessage string                          //Message corresponding to success or failure
	var OutLastReason string                               //Brief message indicating success or failure
	OutSerialContract := make([]map[string]interface{}, 0) //Decimal List containing ID's corresponding to the Column Name and Value given

	//In Arguments
	var InSearchOption string
	var InSearchValue string
	var InSecurityToken string
	var InSessionID string
	var InUsePersistance bool //True if The Activity needs to be linked with ARDS and persisted.
	//-False for none persisting activities.

	// var InSkillID int //Unique Skill ID used on ARDS
	// var InWorkflowServerID int
	// var InXMPPSendMessage bool   //{True - Send messages}  {False - Do not send xmpp message}
	// var InXMPPMessageType string //{"SendResult" - Send only the result}, {"SendError" - Send only errors}, {"SendStatus" - Send only the Status}, {"DonotSend" - Do not send messages}, {"SendAll" - Send all messages}
	// var InXMPPUserName string    //XMPPClient User Name
	// var InXMPPPassword string    //XMPPClient Password
	// var InXMPPServerName string  //XMPPClient User Name

	//---------delete later------
	_ = InSearchOption
	_ = InSearchValue
	_ = InSecurityToken
	_ = InSessionID
	//---------delete later------

	if FlowData["InSessionID"] != nil &&
		FlowData["InSearchOption"] != nil &&
		FlowData["InSearchValue"] != nil &&
		FlowData["InSecurityToken"] != nil &&
		FlowData["InUsePersistance"] != nil {
		InSessionID = FlowData["InSessionID"].(string)
		InSearchOption = FlowData["InSearchOption"].(string)
		InSearchValue = FlowData["InSearchValue"].(string)
		InSecurityToken = FlowData["InSecurityToken"].(string)
		InUsePersistance = FlowData["InUsePersistance"].(bool)

		logger.Log_ACT("Executing : DuoV5_DUOInventoryActivities_GetSerialInfoListBy Activity.", logger.Debug, InSessionID)

		if InUsePersistance {
			logger.Log_ACT("ARDS not implemented in SmoothFlow. Activity will be terminated with error.", logger.Debug, InSessionID)
			err = errors.New("Error : ARDS not implemented in SmoothFlow.")
		} else {
			//Do the magic
			serialHandler := InventoryTransaction.SerialNoHandler{}
			OutSerialContract = serialHandler.SearchSerialNoRecordBy(InSearchOption, InSearchValue, InSecurityToken)
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
		OutLastReason = "Workflow Session ID : " + InSessionID + " Completed without persisting."
		FlowData["OutActivityResult"] = OutActivityResult
		FlowData["OutCustomerMessage"] = OutCustomerMessage
		FlowData["OutLastReason"] = OutLastReason
		FlowData["OutSerialContract"] = OutSerialContract
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
		FlowData["OutSerialContract"] = OutSerialContract
	}

	return FlowData, activityContext
}
