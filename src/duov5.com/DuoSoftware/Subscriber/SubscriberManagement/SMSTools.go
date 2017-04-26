package SubscriberManagement

import (
	"duov5.com/CommonServiceV5/DuoAuthorization"
	"duov5.com/DuoSoftware/Subscriber/SubscriberMasters"
)

type SMSTools struct {
}

func (s *SMSTools) GetAttributeByColumn(ColumnName, ColumnValue, SecurityToken string) (outData []float64) {

	duoAuth := DuoAuthorization.DuoAuthorization{}

	userinfo, _ := duoAuth.GetAccess(SecurityToken, duoAuth.GetAccessCode(
		20000, //Service Code
		1000,  //Functionality Code
		"Write"))
	//Creating the Order Handler Object
	uomHler := SubscriberMasters.AttributeHandler{}
	//Saveing the Orders received
	outData = uomHler.GetAttributesByColumn(
		ColumnName, ColumnValue, // Order List sent to save
		userinfo, // User Profile Sent From the Identity Service
	)
	//retruning the saved Orders
	return

}
