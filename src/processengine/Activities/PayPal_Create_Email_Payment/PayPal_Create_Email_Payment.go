package PayPal_Create_Email_Payment

import "processengine/context"
import "processengine/logger"
import "github.com/logpacker/PayPal-Go-SDK"
import "errors"

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

	var clientID string
	var secretID string
	var emailSubject string
	emailSubject = "Payment From SmoothFlow."
	var receiver string
	var value string
	var currency string
	var optionalNote string
	var optionalItemTrackingID string

	var payoutResp *paypalsdk.PayoutResponse
	var err error

	InSessionID := FlowData["InSessionID"].(string)

	if FlowData["clientID"] != nil && FlowData["secretID"] != nil && FlowData["value"] != nil && FlowData["receiver"] != nil && FlowData["currency"] != nil {
		clientID = FlowData["clientID"].(string)
		secretID = FlowData["secretID"].(string)
		receiver = FlowData["receiver"].(string)
		value = FlowData["value"].(string)
		currency = FlowData["currency"].(string)
	} else {
		err = errors.New("Error! Required Fields Error : Check ClientID, SecretID, value, reciever and currecnty inputs.")
	}

	if FlowData["emailSubject"] != nil {
		emailSubject = FlowData["emailSubject"].(string)
	}

	if FlowData["optionalNote"] != nil {
		emailSubject = FlowData["optionalNote"].(string)
	}

	if FlowData["optionalItemTrackingID"] != nil {
		emailSubject = FlowData["optionalItemTrackingID"].(string)
	}

	if err == nil {
		c, err := paypalsdk.NewClient(clientID, secretID, paypalsdk.APIBaseSandBox)
		if err == nil {
			payout := paypalsdk.Payout{
				SenderBatchHeader: &paypalsdk.SenderBatchHeader{
					EmailSubject: emailSubject,
				},
				Items: []paypalsdk.PayoutItem{
					paypalsdk.PayoutItem{
						RecipientType: "EMAIL",
						Receiver:      receiver,
						Amount: &paypalsdk.AmountPayout{
							Value:    value,
							Currency: currency,
						},
						Note:         optionalNote,
						SenderItemID: optionalItemTrackingID,
					},
				},
			}

			payoutResp, err = c.CreateSinglePayout(payout)
		}
	}

	if err != nil {
		msg := err.Error()
		activityContext.ActivityStatus = false
		activityContext.Message = msg
		activityContext.ErrorState = activityError
		logger.Log_ACT(msg, logger.Debug, InSessionID)
		FlowData["custMsg"] = msg
		FlowData["status"] = "false"
	} else {
		if len(payoutResp.Items) > 0 {
			msg := "Successfully Created PayPal Payment by Email."
			activityContext.ActivityStatus = true
			activityContext.Message = msg
			activityContext.ErrorState = activityError
			logger.Log_ACT(msg, logger.Debug, InSessionID)
			FlowData["custMsg"] = msg
			FlowData["status"] = "true"
			FlowData["PayoutItemID"] = payoutResp.Items[0].PayoutItemID
		} else {
			msg := "Payout Items Not Completed succesfully at Paypal Server!"
			activityContext.ActivityStatus = false
			activityContext.Message = msg
			activityContext.ErrorState = activityError
			logger.Log_ACT(msg, logger.Debug, InSessionID)
			FlowData["custMsg"] = msg
			FlowData["status"] = "false"
		}
	}

	return FlowData, activityContext
}
