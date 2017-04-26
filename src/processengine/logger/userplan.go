package logger

import (
	"duov6.com/common"
	"encoding/json"
	"sync"
	"time"
)

//cebadapter.GetAgent().Client.GetAgentName()

type UserPlan struct {
	Codename                string `json:"codename"`
	PlanName                string `json:"planName"`
	Type                    string `json:"type"`
	Status                  string `json:"status"`
	SubscriptionperiodEnd   string `json:"subscriptionperiodEnd"`
	SubscriptionPeriodStart string `json:"subscriptionPeriodStart"`
	SubscriptionCreated     string `json:"subscriptionCreated"`
}

var plans map[string]UserPlan
var plansLock = sync.RWMutex{}

func ReadPlanMap(sessionID string) (plan UserPlan) {
	//check from map
	plansLock.RLock()
	defer plansLock.RUnlock()
	plan = plans[sessionID]
	return
}

func SetPlanMap(sessionID string, plan UserPlan) {
	plansLock.Lock()
	defer plansLock.Unlock()
	plans[sessionID] = plan
}

func FetchUserPlan(sessionID string) (plan UserPlan) {
	domain := GetDomainBySessionID(sessionID)
	url := "http://" + domain + "/apis/plan/current"
	err, bodyBytes := common.HTTP_GET(url, nil, false)
	if err == nil {
		data := make(map[string]interface{})
		_ = json.Unmarshal(bodyBytes, &data)
		if data["success"].(bool) {
			plan = UserPlan{}
			dataByte, _ := json.Marshal(data["data"])
			_ = json.Unmarshal(dataByte, &plan)
		}
	}
	return
}

func CheckForPlanSupport(sessionID string) bool {
	plan := UserPlan{}
	//check map
	plan = ReadPlanMap(sessionID)
	if plan != (UserPlan{}) && plan.Codename != "Free" {
		if !CheckPlanTimeValidity(plan) {
			//Subscribtion ended it seems. Check one more time if they have renewed.
			newPlan := FetchUserPlan(sessionID)
			if newPlan != (UserPlan{}) && plan.Codename != "Free" {
				if CheckPlanTimeValidity(newPlan) {
					//still invalid. clear the map too
					SetPlanMap(sessionID, UserPlan{})
					return false
				} else {
					//nice.. it has been updated. update the map and allow.
					SetPlanMap(sessionID, newPlan)
					return true
				}
			} else {
				//something went wrong with plan or its been downgraded to FREE, remove entry in map
				SetPlanMap(sessionID, UserPlan{})
				return false
			}
		} else {
			//All okay..
			return true
		}

	} else {
		//map is empty. get from URL
		newPlan := FetchUserPlan(sessionID)
		if newPlan != (UserPlan{}) && newPlan.Codename != "Free" {
			if CheckPlanTimeValidity(newPlan) {
				//all okay set the map entry and go
				SetPlanMap(sessionID, newPlan)
				return true
			} else {
				// havent updated plan subscription.
				return false
			}
		} else {
			//no paid plan is found.
			return false
		}
	}

	return false
}

func CheckPlanTimeValidity(plan UserPlan) bool {
	planEndTime, _ := time.Parse("2006-01-02 15:04:05", plan.SubscriptionperiodEnd)
	if !(time.Now().Before(planEndTime)) {
		return false
	}
	return true
}
