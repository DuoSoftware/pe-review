package Components

import "processengine/context"

type InsParameters struct {
	KeyProperty string `json:"KeyProperty"`
}
type InsertTemplate struct {
	Object     map[string]interface{}   `json:"Object"`
	Objects    []map[string]interface{} `json:"Objects"`
	Parameters InsParameters
}

type ActivityStruct struct {
	ID           string `json:"ID"`
	ActivityName string `json:"ActivityName"`
	Description  string `json:"Description"`
	GoCode       string `json:"GoCode"`
}

type InvokeStruct struct {
	AppCode         string `json:"AppCode"`
	ProcessCode     string `json:"ProcessCode"`
	SessionID       string `json:"SessionID"`
	SecurityToken   string `json:"SecurityToken"`
	Log             string `json:"Log"`
	Namespace       string `json:"Namespace"`
	JSONData        string `json:"JSONData"`
	ResponseMessage chan *context.FlowResult
}

type ProcessMapping struct {
	ID          string `json:"ID"`
	ProcessCode string `json:"ProcessCode"`
	WorkflowID  string `json:"WorkflowID"`
	Name        string `json:"Name"`
}

/*type HibernadedProcess struct {
	AppCode     string                 `json:"AppCode"`
	ProcessCode string                 `json:"ProcessCode"`
	SessionData map[string]interface{} `json:"SessionData"`
	SessionID   string                 `json:"SessionID"`
}
*/
