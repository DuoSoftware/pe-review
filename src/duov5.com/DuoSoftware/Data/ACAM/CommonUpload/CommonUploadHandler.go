package CommonUpload

import (
	"duov5.com/CommonServiceV5/DuoAuthorization"
	"duov6.com/objectstore/client"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type CommonUploadHandler struct {
}

func (c *CommonUploadHandler) GetLatestColumnValueByTemplateCode(securityToken string, templateCode string, columnIndex int, columnValue string) (OutData []string) {
	duoAuth := DuoAuthorization.DuoAuthorization{}
	userinfo, _ := duoAuth.GetAccess(securityToken, duoAuth.GetAccessCode(
		10000, //Service Code
		1000,  //Functionality Code
		"Read"))

	OutData = c.GetLatestColumnValueByTemplateCode2(userinfo, templateCode, columnIndex, columnValue)
	return
}

func (c *CommonUploadHandler) GetLatestColumnValueByTemplateCode2(userInfo DuoAuthorization.UserAuth, templateCode string, columnIndex int, columnValue string) (OutData []string) {
	//fmt.Println(userInfo)
	OutData = make([]string, 0)

	//selectedIds := make([]int64, 0)

	settings := make(map[string]interface{})
	settings["DB_Type"] = "MSSQL"
	settings["Username"] = "smsuser"
	settings["Password"] = "sms"
	settings["Server"] = "203.147.88.139"
	settings["Port"] = "1433"

	// from header in uploadModel.SMS__CommonUploadHeaders where header.TemplateCode == templateCode &&
	//  header.CompanyId == userInfo.CompanyID && header.TenantId == userInfo.TenantID &&
	//   userInfo.viweObjectIDs.Contains((int)header.ViewObjectId) select header

	inQueryArgs := ""

	for _, value := range userInfo.ViweObjectIDs.Int {
		inQueryArgs += strconv.Itoa(value) + ","
	}
	inQueryArgs = strings.TrimSuffix(inQueryArgs, ",")

	//testing only. remove when live
	//userInfo.TenantID = 3
	//userInfo.CompanyID = 2

	query := "select HeaderGuRefId from [MBV5DBLive].[dbo].[SMS__CommonUploadHeader] where TemplateCode='" + templateCode + "' AND 	CompanyID=" + strconv.Itoa(userInfo.CompanyID) + " AND  TenantId=" + strconv.Itoa(userInfo.TenantID) + " AND viewObjectID IN (" + inQueryArgs + ");"
	//fmt.Println(query)
	data, errA := client.GoSmoothFlow("securityToken", "ignore", "ignore", settings).GetMany().ByQuerying(query).Ok()
	if errA != "" {
		fmt.Println(errA)
	} else {
		//fmt.Println(string(data))
		arr := make([]map[string]interface{}, 0)
		erdd := json.Unmarshal(data, &arr)
		if erdd != nil {
			fmt.Println(erdd.Error())
		} else {

			for x := 0; x < len(arr); x++ {
				OutData = append(OutData, strconv.FormatFloat(arr[x]["HeaderGuRefId"].(float64), 'f', -1, 64))
			}
		}
	}

	return
}
