package SubscriberMasters

import (
	"duov5.com/CommonServiceV5/DuoAuthorization"
	"duov6.com/objectstore/client"
	"encoding/json"
	"fmt"
)

type AttributeHandler struct {
}

func (s *AttributeHandler) GetAttributesByColumn(ColumnName, ColumnValue string, userInfo DuoAuthorization.UserAuth) []float64 {
	attList := make([]float64, 0)

	settings := make(map[string]interface{})
	settings["DB_Type"] = "MSSQL"
	settings["Username"] = "smsuser"
	settings["Password"] = "sms"
	settings["Server"] = "203.147.88.139"
	settings["Port"] = "1433"

	query := "select GURefID from [MBV5DBLive].[dbo].[SMS__Attributes] where ColumnName='" + ColumnName + "' AND ColumnValue='" + ColumnValue + "';"

	queryData, _ := client.GoSmoothFlow("securityToken", "ignore", "ignore", settings).GetMany().ByQuerying(query).Ok()

	arr := make([]map[string]interface{}, 0)
	erdd := json.Unmarshal(queryData, &arr)
	if erdd != nil {
		fmt.Println(erdd.Error())
	} else {
		for x := 0; x < len(arr); x++ {
			attList = append(attList, arr[x]["GURefID"].(float64))
		}
	}

	return attList
}
