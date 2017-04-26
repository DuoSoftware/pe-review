package InventoryTransaction

import (
	"duov5.com/CommonServiceV5/DuoAuthorization"
	"duov6.com/objectstore/client"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type SerialNoHandler struct {
}

func (s *SerialNoHandler) SearchSerialNoRecordBy(SearchBy string, Value string, securityToken string) []map[string]interface{} {
	outData := make([]map[string]interface{}, 0)

	duoAuth := DuoAuthorization.DuoAuthorization{}

	userinfo, _ := duoAuth.GetAccess(securityToken, duoAuth.GetAccessCode(
		10000, //Service Code
		1000,  //Functionality Code
		"Write"))

	outData = s.SearchSerialNoRecordByAuth(SearchBy, Value, userinfo)

	return outData
}

func (s *SerialNoHandler) SearchSerialNoRecordByAuth(SearchBy string, Value string, user DuoAuthorization.UserAuth) []map[string]interface{} {

	outData := make([]map[string]interface{}, 0)

	SNoQueryGbl := make([]map[string]interface{}, 0)
	SNoQueryGbldetail := make([]map[string]interface{}, 0)

	visitHeader := false
	visitDetail := false

	switch SearchBy {
	case "SerialNo":
		//	visitHeader = true
		tenantIDInQueryPart := s.GetINQueryPart(user.TenantIDs.Int, "int")
		companyIDInQueryPart := s.GetINQueryPart(user.CompanyIDs.Int, "int")
		query := "select * from [MBV5DBLive].[dbo].[INV_SerialNumberHeader] where SerialNumber='" + Value + "' AND TenantID IN (" + tenantIDInQueryPart + ") AND CompanyID IN(" + companyIDInQueryPart + ");"
		SNoQueryGbl = s.ReadFromObjectStore(query)
		break

	case "SearchSerial_Like":
		visitHeader = true
		tenantIDInQueryPart := s.GetINQueryPart(user.TenantIDs.Int, "int")
		companyIDInQueryPart := s.GetINQueryPart(user.CompanyIDs.Int, "int")
		query := "select * from [MBV5DBLive].[dbo].[INV_SerialNumberHeader] where SerialNumber LIKE '" + "%" + Value + "%" + "' AND TenantID IN (" + tenantIDInQueryPart + ") AND CompanyID IN(" + companyIDInQueryPart + ");"
		SNoQueryGbl = s.ReadFromObjectStore(query)
		break

	case "storeCode":
		visitHeader = true
		tenantIDInQueryPart := s.GetINQueryPart(user.TenantIDs.Int, "int")
		companyIDInQueryPart := s.GetINQueryPart(user.CompanyIDs.Int, "int")
		query := "select * from [MBV5DBLive].[dbo].[INV_SerialNumberHeader] where StoreCode='" + Value + "' AND TenantID IN (" + tenantIDInQueryPart + ") AND CompanyID IN(" + companyIDInQueryPart + ");"
		SNoQueryGbl = s.ReadFromObjectStore(query)
		break

	case "Status":
		visitHeader = true
		tenantIDInQueryPart := s.GetINQueryPart(user.TenantIDs.Int, "int")
		companyIDInQueryPart := s.GetINQueryPart(user.CompanyIDs.Int, "int")
		query := "select * from [MBV5DBLive].[dbo].[INV_SerialNumberHeader] where Status='" + Value + "' AND TenantID IN (" + tenantIDInQueryPart + ") AND CompanyID IN(" + companyIDInQueryPart + ");"
		SNoQueryGbl = s.ReadFromObjectStore(query)
		break

	case "ChipID":
		//  //visitHeader = true;
		// //SNoQuery = InvContext.Where(row => row.chipID == Value && user.TenantIDs.Contains(row.tenantID) && user.CompanyIDs.Contains(row.companyID));
		//  //SNoQueryGbl = SNoQuery.ToList();
		break

	case "GUStoreID":
		visitHeader = true
		tenantIDInQueryPart := s.GetINQueryPart(user.TenantIDs.Int, "int")
		companyIDInQueryPart := s.GetINQueryPart(user.CompanyIDs.Int, "int")
		query := "select * from [MBV5DBLive].[dbo].[INV_SerialNumberHeader] where GUStoreID=" + Value + " AND TenantID IN (" + tenantIDInQueryPart + ") AND CompanyID IN(" + companyIDInQueryPart + ");"
		SNoQueryGbl = s.ReadFromObjectStore(query)
		break

	case "GUSerialNoID":
		visitHeader = true
		tenantIDInQueryPart := s.GetINQueryPart(user.TenantIDs.Int, "int")
		companyIDInQueryPart := s.GetINQueryPart(user.CompanyIDs.Int, "int")
		query := "select * from [MBV5DBLive].[dbo].[INV_SerialNumberHeader] where GUSerialNoID=" + Value + " AND TenantID IN (" + tenantIDInQueryPart + ") AND CompanyID IN(" + companyIDInQueryPart + ");"
		SNoQueryGbl = s.ReadFromObjectStore(query)
		break

	case "AllSerialNo":
		visitHeader = true
		tenantIDInQueryPart := s.GetINQueryPart(user.TenantIDs.Int, "int")
		companyIDInQueryPart := s.GetINQueryPart(user.CompanyIDs.Int, "int")
		query := "select * from [MBV5DBLive].[dbo].[INV_SerialNumberHeader] where TenantID IN (" + tenantIDInQueryPart + ") AND CompanyID IN(" + companyIDInQueryPart + ");"
		SNoQueryGbl = s.ReadFromObjectStore(query)
		break

		//case SerialNoContract.SNoSearchBy.MacAddress:
		//    visitHeader = true;
		//    SNoQuery = InvContext.INVSerialNumberHeaders.Where(row => row.macAddress == Value && user.TenantIDs.Contains(row.tenantID) && user.CompanyIDs.Contains(row.companyID));
		//    SNoQueryGbl = SNoQuery.ToList();
		//    break;

	case "GUGRNID":
		visitDetail = true
		tenantIDInQueryPart := s.GetINQueryPart(user.TenantIDs.Int, "int")
		companyIDInQueryPart := s.GetINQueryPart(user.CompanyIDs.Int, "int")
		query := "select * from [MBV5DBLive].[dbo].[INV_SerialNumberDetail] where GURefID=" + Value + " AND TenantID IN (" + tenantIDInQueryPart + ") AND CompanyID IN(" + companyIDInQueryPart + ");"
		SNoQueryGbldetail = s.ReadFromObjectStore(query)
		break

	case "GUAODID":
		visitDetail = true
		tenantIDInQueryPart := s.GetINQueryPart(user.TenantIDs.Int, "int")
		companyIDInQueryPart := s.GetINQueryPart(user.CompanyIDs.Int, "int")
		query := "select * from [MBV5DBLive].[dbo].[INV_SerialNumberDetail] where GURefID=" + Value + " AND TenantID IN (" + tenantIDInQueryPart + ") AND CompanyID IN(" + companyIDInQueryPart + ");"
		SNoQueryGbldetail = s.ReadFromObjectStore(query)
		break

	case "GUItemID":
		visitHeader = true
		tenantIDInQueryPart := s.GetINQueryPart(user.TenantIDs.Int, "int")
		companyIDInQueryPart := s.GetINQueryPart(user.CompanyIDs.Int, "int")
		query := "select * from [MBV5DBLive].[dbo].[INV_SerialNumberHeader] where GUItemID=" + Value + " AND TenantID IN (" + tenantIDInQueryPart + ") AND CompanyID IN(" + companyIDInQueryPart + ");"
		SNoQueryGbl = s.ReadFromObjectStore(query)
		break

	case "Issued":
		tenantIDInQueryPart := s.GetINQueryPart(user.TenantIDs.Int, "int")
		companyIDInQueryPart := s.GetINQueryPart(user.CompanyIDs.Int, "int")
		query := "select * from [MBV5DBLive].[dbo].[INV_SerialNumberHeader] where Issued=" + Value + " AND TenantID IN (" + tenantIDInQueryPart + ") AND CompanyID IN(" + companyIDInQueryPart + ");"
		SNoQueryGbl = s.ReadFromObjectStore(query)
		break

	case "ItemCode":
		visitHeader = true
		tenantIDInQueryPart := s.GetINQueryPart(user.TenantIDs.Int, "int")
		companyIDInQueryPart := s.GetINQueryPart(user.CompanyIDs.Int, "int")
		query := "select * from [MBV5DBLive].[dbo].[INV_SerialNumberHeader] where ItemCode='" + Value + "' AND TenantID IN (" + tenantIDInQueryPart + ") AND CompanyID IN(" + companyIDInQueryPart + ");"
		SNoQueryGbl = s.ReadFromObjectStore(query)
		break

	case "GUTranID":
		visitHeader = true
		tenantIDInQueryPart := s.GetINQueryPart(user.TenantIDs.Int, "int")
		companyIDInQueryPart := s.GetINQueryPart(user.CompanyIDs.Int, "int")
		query := "select * from [MBV5DBLive].[dbo].[INV_SerialNumberHeader] where GUTranID=" + Value + " AND TenantID IN (" + tenantIDInQueryPart + ") AND CompanyID IN(" + companyIDInQueryPart + ");"
		SNoQueryGbl = s.ReadFromObjectStore(query)
		break

	case "CommitStatus":
		visitHeader = true
		tenantIDInQueryPart := s.GetINQueryPart(user.TenantIDs.Int, "int")
		companyIDInQueryPart := s.GetINQueryPart(user.CompanyIDs.Int, "int")
		query := "select * from [MBV5DBLive].[dbo].[INV_SerialNumberHeader] where CommitStatus=" + Value + " AND TenantID IN (" + tenantIDInQueryPart + ") AND CompanyID IN(" + companyIDInQueryPart + ");"
		SNoQueryGbl = s.ReadFromObjectStore(query)
		break

	}

	if visitHeader {
		for _, object := range SNoQueryGbl {
			outData = append(outData, object)
		}
	}

	if visitDetail {
		for _, object := range SNoQueryGbldetail {
			outData = append(outData, object)
		}
	}

	return outData
}

// public enum SNoSearchBy
// {
//     [EnumMember(Value = "SerialNo")]
//     SerialNo,
//     [EnumMember(Value = "GUStoreID")]
//     GUStoreID,
//     [EnumMember(Value = "Status")]
//     Status,
//     [EnumMember(Value = "GUSerialNoID")]
//     GUSerialNoID,

//     [EnumMember(Value = "ChipID")]
//     ChipID,
//     [EnumMember(Value = "MacAddress")]
//     MacAddress,
//     [EnumMember(Value = "GUGRNID")]
//     GUGRNID,
//     [EnumMember(Value = "GUAODID")]
//     GUAODID,
//     [EnumMember(Value = "GUItemID")]
//     GUItemID,
//     [EnumMember(Value = "Issued")]
//     Issued,
//     [EnumMember(Value = "storeCode")]
//     storeCode,
//     [EnumMember(Value = "ItemCode")]
//     ItemCode,
//     [EnumMember(Value = "AllSerialNo")]
//     AllSerialNo,
//     [EnumMember(Value = "SearchSerial_Like")]
//     SearchSerial_Like,

//     [EnumMember(Value = "GUTranID")]
//     GUTranID,
//     [EnumMember(Value = "CommitStatus")]
//     CommitStatus,

//     [EnumMember(Value = "EchoDic_By_GUDirectPurchaseID")]
//     EchoDic_By_GUDirectPurchaseID,

// }

func (s *SerialNoHandler) ReadFromObjectStore(query string) []map[string]interface{} {
	arr := make([]map[string]interface{}, 0)

	settings := make(map[string]interface{})
	settings["DB_Type"] = "MSSQL"
	settings["Username"] = "smsuser"
	settings["Password"] = "sms"
	settings["Server"] = "203.147.88.139"
	settings["Port"] = "1433"

	queryData, _ := client.GoSmoothFlow("securityToken", "ignore", "ignore", settings).GetMany().ByQuerying(query).Ok()

	erdd := json.Unmarshal(queryData, &arr)
	if erdd != nil {
		fmt.Println(erdd.Error())
	}
	return arr
}

func (s *SerialNoHandler) GetINQueryPart(data interface{}, Type string) (queryPart string) {
	if Type == "string" {
		rowData := data.([]string)
		for _, value := range rowData {
			queryPart += "'" + value + "',"
		}
		queryPart = strings.TrimSuffix(queryPart, ",")
	} else if Type == "int" {
		rowData := data.([]int)
		for _, value := range rowData {
			queryPart += strconv.Itoa(value) + ","
		}
		queryPart = strings.TrimSuffix(queryPart, ",")
	}
	return
}
