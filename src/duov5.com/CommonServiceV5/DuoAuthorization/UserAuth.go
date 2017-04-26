package DuoAuthorization

import (
	"duov6.com/common"
	"encoding/xml"
	"fmt"
	"strconv"
)

type UserAuthBody struct {
	Body Body `xml:"Body"`
}

type Body struct {
	GetAccessResponse GetAccessResponse `xml:"GetAccessResponse"`
}

type GetAccessResponse struct {
	GetAccessResult UserAuth `xml:"GetAccessResult"`
}

type UserAuth struct {
	AccountContracts AccountContracts `xml:"AccountContracts"`
	AccountRelated   bool             `xml:"AccountRelated"`
	Application      string           `xml:"Application"`
	CompanyID        int              `xml:"CompanyID"`
	CompanyIDs       CompanyIDs       `xml:"CompanyIDs"`
	Data             Data             `xml:"Data"`
	IgnoreViweObj    bool             `xml:"IgnoreViweObj"`
	ObjectID         int              `xml:"ObjectID"`
	SecurityToken    string           `xml:"SecurityToken"`
	TenantID         int              `xml:"TenantID"`
	TenantIDs        TenantIDs        `xml:"TenantIDs"`
	TokenExpireOn    string           `xml:"TokenExpireOn"`
	Type             string           `xml:"Type"`
	UserName         string           `xml:"UserName"`
	Write            bool             `xml:"Write"`
	GuUserGrpID      GuUserGrpID      `xml:"guUserGrpID"`
	GuUserId         int64            `xml:"guUserId"`
	ViweObjectIDs    ViweObjectIDs    `xml:"viweObjectIDs"`
}

type CompanyIDs struct {
	Int []int `xml:"int"`
}

type Data struct {
	KeyValueOfstringstring []KeyValueOfstringstring `xml:"KeyValueOfstringstring"`
}

type KeyValueOfstringstring struct {
	Key   string
	Value string
}

type TenantIDs struct {
	Int []int `xml:"int"`
}

type ViweObjectIDs struct {
	Int []int `xml:"int"`
}

type AccountContracts struct {
	GUAccountID   float64
	AccountNo     string
	GUPromotionID float64
}

type GuUserGrpID struct {
	Decimal []float64 `xml:"decimal"`
}

func (u *UserAuth) GetUserAuthObjectFromXML(document string) (authObj UserAuth, err error) {
	authObj = UserAuth{}
	res := UserAuthBody{}

	if err = xml.Unmarshal([]byte(document), &res); err == nil {
		authObj = res.Body.GetAccessResponse.GetAccessResult
	}

	return
}

type DuoAuthorization struct {
}

func (d *DuoAuthorization) GetAccess(securityToken, accessCode string) (authObject UserAuth, err error) {
	authObject = UserAuth{}

	//Call to Auth Service
	messageBody := `<x:Envelope xmlns:x="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tem="http://tempuri.org/">
    <x:Header/>
    <x:Body>
        <tem:GetAccess>
            <tem:SecurityToken>` + securityToken + `</tem:SecurityToken>
            <tem:AccessCode>` + accessCode + `</tem:AccessCode>
        </tem:GetAccess>
    </x:Body>
</x:Envelope>`

	url := "http://192.168.1.194/DuoSubscribe5/CommonServices/Authorization/auth.svc"
	headers := make(map[string]string)
	headers["Content-Type"] = "text/xml"
	headers["SOAPAction"] = "http://tempuri.org/Iauth/GetAccess"

	if err, bodyByteArray := common.HTTP_POST(url, headers, []byte(messageBody), false); err == nil {
		authObject, err = (&authObject).GetUserAuthObjectFromXML(string(bodyByteArray))
	} else {
		fmt.Println(err.Error())
	}
	return
}

func (d *DuoAuthorization) GetAccessCode(ServiceCode, FunctionCode int, ServiceType string) string {
	return strconv.Itoa(ServiceCode) + "#" + strconv.Itoa(FunctionCode) + "#" + ServiceType
}
