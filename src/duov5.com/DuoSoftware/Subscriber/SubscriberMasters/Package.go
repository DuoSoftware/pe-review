package SubscriberMasters

import (
	"duov5.com/CommonServiceV5/DuoAuthorization"
	"duov6.com/objectstore/client"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Package struct {
}

func (p *Package) GetPackageList(guPackageIDList []float64, SecurityToken string) []map[string]interface{} {
	outData := make([]map[string]interface{}, 0)
	duoAuth := DuoAuthorization.DuoAuthorization{}

	userinfo, _ := duoAuth.GetAccess(SecurityToken, duoAuth.GetAccessCode(
		10000, //Service Code
		1000,  //Functionality Code
		"Write"))

	//Saveing the Orders received
	outData = p.GetPackageListByAuth(guPackageIDList, userinfo)
	//retruning the saved Orders
	return outData

}

func (p *Package) GetPackageListByAuth(guPackageIDList []float64, userinfo DuoAuthorization.UserAuth) []map[string]interface{} {
	outData := make([]map[string]interface{}, 0)
	// var result = context.SMS_PackageMasters.Where(search => guPackageIDList.Contains(search.GUPackageID) &&
	// 	userinfo.CompanyIDs.Contains(search.CompanyID) &&
	//     userinfo.TenantIDs.Contains(search.TenantID) && userinfo.viweObjectIDs.Contains(search.ViewObjectID)).ToList();
	// //return result.Select(select => new PackageMaster()
	// {
	//     CommitStatus = select.CommitStatus.Value,
	//     CompanyID = select.CompanyID,
	//     CreateDate = select.CreateDate.Value,
	//     CreateUser = select.CreateUser,
	//     GUPackageID = select.GUPackageID,
	//     PackageCategory = select.PackageCategory.ToString(),
	//     PackageClass = select.PackageClass,
	//     PackageCode = select.PackageCode,
	//     PackageType = select.PackageType,
	//     Rental = select.Rental.Value,
	//     Version = select.Version.Value,
	//     PackageDescription = select.PackageDescription,
	//     PrintDescription = select.PrintDescription
	// }).ToList();

	settings := make(map[string]interface{})
	settings["DB_Type"] = "MSSQL"
	settings["Username"] = "smsuser"
	settings["Password"] = "sms"
	settings["Server"] = "203.147.88.139"
	settings["Port"] = "1433"

	//Create guPackageIDList IN query
	guPackageIDListInQueryPart := ""
	for _, value := range guPackageIDList {
		guPackageIDListInQueryPart += strconv.FormatFloat(value, 'f', -1, 64) + ","
	}
	guPackageIDListInQueryPart = strings.TrimSuffix(guPackageIDListInQueryPart, ",")

	//Create companyIDs IN query
	companyIDsInQueryPart := ""
	for _, value := range userinfo.CompanyIDs.Int {
		companyIDsInQueryPart += strconv.Itoa(value) + ","
	}
	companyIDsInQueryPart = strings.TrimSuffix(companyIDsInQueryPart, ",")

	//Create guPackageIDList IN query
	tenantIDsInQueryPart := ""
	for _, value := range userinfo.TenantIDs.Int {
		tenantIDsInQueryPart += strconv.Itoa(value) + ","
	}
	tenantIDsInQueryPart = strings.TrimSuffix(tenantIDsInQueryPart, ",")

	//Create guPackageIDList IN query
	viweObjectIDsInQueryPart := ""
	for _, value := range userinfo.ViweObjectIDs.Int {
		viweObjectIDsInQueryPart += strconv.Itoa(value) + ","
	}
	viweObjectIDsInQueryPart = strings.TrimSuffix(viweObjectIDsInQueryPart, ",")

	query := "select * from [MBV5DBLive].[dbo].[SMS_PackageMaster] where GUPackageID IN (" + guPackageIDListInQueryPart + ") AND CompanyID IN (" + companyIDsInQueryPart + ") AND TenantID IN (" + tenantIDsInQueryPart + ") AND ViewObjectID IN (" + viweObjectIDsInQueryPart + ");"
	//fmt.Println(query)
	queryData, _ := client.GoSmoothFlow("securityToken", "ignore", "ignore", settings).GetMany().ByQuerying(query).Ok()

	erdd := json.Unmarshal(queryData, &outData)
	if erdd != nil {
		fmt.Println(erdd.Error())
	}
	//fmt.Println(string(queryData))
	return outData
}
