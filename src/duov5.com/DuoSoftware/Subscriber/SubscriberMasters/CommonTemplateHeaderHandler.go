package SubscriberMasters

import (
	"duov5.com/CommonServiceV5/DuoAuthorization"
	"duov5.com/CommonServiceV5/DuoAuthorization/CommonTools"
	"duov6.com/objectstore/client"
	"encoding/json"
	//"fmt"
	"strconv"
	"strings"
)

type CommonTemplateHeaderHandler struct {
}

func (c *CommonTemplateHeaderHandler) GetAll(userinfo DuoAuthorization.UserAuth) []map[string]interface{} {
	tmplist := make([]map[string]interface{}, 0)

	settings := make(map[string]interface{})
	settings["DB_Type"] = "MSSQL"
	settings["Username"] = "smsuser"
	settings["Password"] = "sms"
	settings["Server"] = "203.147.88.139"
	settings["Port"] = "1433"

	//Get Data from SMS__CommonTemplateHeader
	companyIds := ""
	for _, value := range userinfo.CompanyIDs.Int {
		companyIds += strconv.Itoa(value) + ","
	}
	companyIds = strings.TrimSuffix(companyIds, ",")

	viewObjectIds := ""
	for _, value := range userinfo.ViweObjectIDs.Int {
		viewObjectIds += strconv.Itoa(value) + ","
	}
	viewObjectIds = strings.TrimSuffix(viewObjectIds, ",")

	//hard coded.. remove when deploy
	//userinfo.TenantID = 3

	itemsQuery := "select * from [MBV5DBLive].[dbo].[SMS_CommonTemplateHeader] where CompanyId IN (" + companyIds + ") AND TenantId=" + strconv.Itoa(userinfo.TenantID) + " AND CommitStatus=" + strconv.Itoa(CommonTools.Commit) + " AND ViewObjectId IN (" + viewObjectIds + ");"

	itemData, _ := client.GoSmoothFlow("securityToken", "ignore", "ignore", settings).GetMany().ByQuerying(itemsQuery).Ok()

	templateHeaderItems := make([]map[string]interface{}, 0)
	_ = json.Unmarshal(itemData, &templateHeaderItems)

	//iterate thru all and read SMS_CommonTemplateDetails

	for x := 0; x < len(templateHeaderItems); x++ {
		itemID := templateHeaderItems[x]["ID"].(float64)
		//get records which matches ID from SMS_CommonTemplateDetails
		detailQuery := "select * from [MBV5DBLive].[dbo].[SMS_CommonTemplateDetail] where TemplateID=" + strconv.FormatFloat(itemID, 'E', -1, 64) + ";"

		detailData, _ := client.GoSmoothFlow("securityToken", "ignore", "ignore", settings).GetMany().ByQuerying(detailQuery).Ok()

		templateDetailItems := make([]map[string]interface{}, 0)
		_ = json.Unmarshal(detailData, &templateDetailItems)

		templateHeaderItems[x]["CommonTemplateDetailList"] = templateDetailItems

	}

	tmplist = templateHeaderItems

	return tmplist
}
