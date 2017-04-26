package Masters

import (
	"duov5.com/CommonServiceV5/DuoAuthorization"
	"duov5.com/DuoSoftware/Subscriber/SubscriberMasters"
	//"fmt"
)

func GetAllCommonTemplateHeader(SecurityToken string) []map[string]interface{} {
	duoAuth := DuoAuthorization.DuoAuthorization{}

	userinfo, _ := duoAuth.GetAccess(SecurityToken, duoAuth.GetAccessCode(
		10000, //Service Code
		1000,  //Functionality Code
		"Read"))

	//fmt.Println(userinfo)

	handler := SubscriberMasters.CommonTemplateHeaderHandler{}
	return handler.GetAll(userinfo)

}
