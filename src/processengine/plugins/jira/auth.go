package jira

import (
	"duov6.com/common"
	"encoding/json"
	"fmt"
	"processengine/Common"
	"processengine/Components"
	"strings"
)

type AuthData struct {
	Domain string
	Email  string
}

func GetExternalPlugins(email string) (con Components.ExternalPlugins) {
	con = Components.ExternalPlugins{}
	url := "http://" + Common.OBJECTSTORE_URL + "/com.smoothflow.io/connectors/" + email + "?securityToken=ignore"
	err, bodyBytes := common.HTTP_GET(url, nil, false)
	if err == nil && len(bodyBytes) > 4 {
		_ = json.Unmarshal(bodyBytes, &con)
	}
	return
}

func JiraAuthorize(jiraDomain string, jiraTokens map[string]string) (authRes JiraResponse) {
	authRes = JiraResponse{}

	isAdmin, profile := VerifyJiraAccess(jiraDomain, jiraTokens)
	_ = profile
	_ = isAdmin
	if isAdmin {
		fmt.Println("Verified Jira Admin.")
	} else {
		fmt.Println("Not a Jira Admin.")
		authRes.IsSuccess = false
		authRes.Message = "Access Denied. Jira administrator rights not found for the user."
		return
	}

	//Get addonlink object
	jiraCon := GetAddonLink(jiraDomain)

	fmt.Print("Jira Addon Link : ")
	fmt.Println(jiraCon)

	if Common.AUTH_URL == "" {
		authRes.IsSuccess = false
		authRes.Message = "Server Error. Authentication server access not found. Please contact smoothflow technical staff."
		return
	}

	securityToken := VerifySFAdminUser()
	if securityToken == "" {
		authRes.IsSuccess = false
		authRes.Message = "Server Error. Server administration error. Please contact smoothflow technical staff."
		return
	}

	if jiraCon.JiraDomain == "" {
		//This is not even initialized.
		//Check for domain availability
		//If YES. check if tenant admin. if tenant admin add to Jira integrations. If no. Send the tenant admin emails and request for integration access.
		//If NO. return with registration needed error

		//Get Tenant admins.
		url := "http://" + Common.AUTH_URL + "/tenant/GetTenantAdmin/" + jiraDomain + ".dev.smoothflow.io"
		headers := make(map[string]string)
		headers["SecurityToken"] = securityToken
		err, tenantBody := common.HTTP_GET(url, headers, false)
		if err != nil {
			authRes.IsSuccess = false
			authRes.Message = "Server Error. Unable to contact authentication server. Please contact smoothflow technical staff."
		} else if len(tenantBody) > 4 && strings.Contains(string(tenantBody), profile.EmailAddress) {
			//Domain has been created earlier and this dude in this as a tenant admin. proceed.
			jiraCon.Email = profile.EmailAddress
			jiraCon.JiraDomain = jiraDomain
			err := SetAddonLink(jiraCon)
			if err != nil {
				authRes.IsSuccess = false
				authRes.Message = err.Error()
			} else {
				authRes.IsSuccess = true
				authRes.Message = "Successfully Authorized Jira Domain : " + jiraDomain + " from SmoothFlow via Account belongs to : " + profile.EmailAddress + "."
			}
		} else {
			authRes.IsSuccess = false
			authRes.Message = "Not authenticated for JIRA integrations. Registration required."
		}
	} else if jiraCon.JiraDomain != "" && jiraCon.Email == "" {
		//This account has been deauthrized and not authrized to an email yet.
		//check if tenant admin. if not fail. return a message saying continue to smoothflow and make assignments.
		url := "http://" + Common.AUTH_URL + "/tenant/GetTenantAdmin/" + jiraDomain + ".dev.smoothflow.io"
		headers := make(map[string]string)
		headers["SecurityToken"] = securityToken
		err, tenantBody := common.HTTP_GET(url, headers, false)
		if err != nil {
			authRes.IsSuccess = false
			authRes.Message = err.Error()
		} else if strings.Contains(string(tenantBody), profile.EmailAddress) {
			//This request came from a tenant admin
			jiraCon.Email = profile.EmailAddress
			//save to ObjecStore
			err := SetAddonLink(jiraCon)
			if err != nil {
				authRes.IsSuccess = false
				authRes.Message = err.Error()
			} else {
				authRes.IsSuccess = true
				authRes.Message = "Successfully Authorized Jira Domain : " + jiraDomain + " from SmoothFlow via Account belongs to : " + profile.EmailAddress + "."
			}
		} else {
			authRes.IsSuccess = false
			authRes.Message = "Access Denied. Request admin access for : " + jiraDomain + ".dev.smoothflow.io tenant to change ownership."
		}
	} else {
		//all okay
		authRes.IsSuccess = true
		authRes.Message = "Already Authorized."
		ad := AuthData{}
		ad.Domain = jiraCon.JiraDomain
		ad.Email = jiraCon.Email
		authRes.Data = ad
	}
	return
}

func JiraUnAuthorize(jiraDomain string, jiraTokens map[string]string) (authRes JiraResponse) {
	authRes = JiraResponse{}

	isAdmin, profile := VerifyJiraAccess(jiraDomain, jiraTokens)

	if isAdmin {
		fmt.Println("Verified Jira Admin.")
	} else {
		fmt.Println("Not a Jira Admin.")
		authRes.IsSuccess = false
		authRes.Message = "Access Denied. Jira administrator rights not found for the user."
		return
	}

	//Get addonlink object
	jiraCon := GetAddonLink(jiraDomain)

	fmt.Print("Jira Addon Link : ")
	fmt.Println(jiraCon)

	if Common.AUTH_URL == "" {
		authRes.IsSuccess = false
		authRes.Message = "Server Error. Authentication server access not found. Please contact smoothflow technical staff."
		return
	}

	securityToken := VerifySFAdminUser()
	if securityToken == "" {
		authRes.IsSuccess = false
		authRes.Message = "Server Error. Server administration error. Please contact smoothflow technical staff."
		return
	} else {
		//securityToken is there. Check if domain admin.
		url := "http://" + Common.AUTH_URL + "/tenant/GetTenantAdmin/" + jiraDomain + ".dev.smoothflow.io"
		headers := make(map[string]string)
		headers["SecurityToken"] = securityToken
		err, tenantBody := common.HTTP_GET(url, headers, false)
		if err != nil {
			authRes.IsSuccess = false
			authRes.Message = "Server Error. Unable to contact authentication server. Please contact smoothflow technical staff."
			return
		} else if len(tenantBody) > 4 && !strings.Contains(string(tenantBody), profile.EmailAddress) {
			//Not a tenant admin.. No access
			authRes.IsSuccess = false
			authRes.Message = "Access Denied. No domain admin access found in smoothflow."
			return
		}
	}

	if jiraCon.JiraDomain != "" && jiraCon.Email != "" {
		//check if token email is equal to email in addonlink and proceed.
		//authObj := Authenticate(jiraDomain, email)
		authObj := JiraUser{}
		if authObj.EmailAddress != "" {
			//Authenticated.. So remove the email!
			jiraCon.Email = ""
			//save to ObjecStore
			err := SetAddonLink(jiraCon)

			if err != nil {
				authRes.IsSuccess = false
				authRes.Message = err.Error()
			} else {
				authRes.IsSuccess = true
				authRes.Message = "Successfully Deauthorized Jira Domain : " + jiraDomain + " from SmoothFlow."
			}

		} else {
			authRes.IsSuccess = false
			authRes.Message = "This user not authorzed in SmoothFlow to make changes. Create account under " + jiraDomain + " tenant in smoothflow and retry. Please contact your Jira Administrator for furthur instructions."
		}
	} else if jiraCon.JiraDomain != "" && jiraCon.Email == "" {
		authRes.IsSuccess = true
		authRes.Message = "Already Unauthorized. Nothing to do"
	} else if jiraCon.JiraDomain == "" && jiraCon.Email == "" {
		authRes.IsSuccess = false
		authRes.Message = "No Jira domain entry under :" + jiraDomain + " found. Nothing to be unauthorized."
	}

	return
}

func VerifyJiraAccess(domain string, tokens map[string]string) (isAdmin bool, user JiraUser) {
	//Get Self
	url := "https://" + domain + "/rest/api/2/myself"
	err, selfBody := Jira_HTTP_GET(url, tokens)
	if err != nil {
		isAdmin = false
	} else {
		//Object available
		selfObject := make(map[string]interface{})
		_ = json.Unmarshal(selfBody, &selfObject)

		//Get the self url append groups request and resend
		url = selfObject["self"].(string) + "&expand=groups"
		err, groupBody := Jira_HTTP_GET(url, tokens)
		if err != nil {
			//couldnt retrieve group properties..
			isAdmin = false
		} else {
			//recieved group properties
			_ = json.Unmarshal(groupBody, &user)
			items := user.Groups["items"].([]interface{})
			itemsArray := Common.ConvertInterfaceArrayToObjectArray(items)
			fmt.Println("Jira Groups for User : ")
			for _, arrayObj := range itemsArray {
				fmt.Println(arrayObj["name"].(string))
				if strings.Contains(arrayObj["name"].(string), "admin") {
					isAdmin = true
					break
				}
			}
			fmt.Println()
		}
	}
	return
}

func VerifySFAdminUser() (securityToken string) {
	//Try to login
	url := "http://" + Common.AUTH_URL + "/Login/prasadacicts@gmail.com/123/smoothflow.io"
	err, loginBody := common.HTTP_GET(url, nil, false)
	if err != nil && strings.Contains(err.Error(), "The username or password is incorrect.") {
		//Create User
		url = "http://" + Common.AUTH_URL + "/InvitedUserRegistration/"
		userBody := `{"EmailAddress":"prasadacicts@gmail.com","Name":"SmoothFlow Admin Test Account","Password":"123","ConfirmPassword":"123"}`
		err, signBody := common.HTTP_POST(url, nil, []byte(userBody), false)
		if err != nil {
			fmt.Println(err.Error())
			return
		} else {
			fmt.Println(string(signBody))
			//Login Again
			url = "http://" + Common.AUTH_URL + "/Login/prasadacicts@gmail.com/123/smoothflow.io"
			err, loginBody = common.HTTP_GET(url, nil, false)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}

	loginMap := make(map[string]interface{})
	_ = json.Unmarshal(loginBody, &loginMap)
	securityToken = loginMap["SecurityToken"].(string)

	return
}
