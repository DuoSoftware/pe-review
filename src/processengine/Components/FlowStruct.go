package Components

import "processengine/context"

// SourceEndpoint struct
type SourceEndpoint struct {
	ID       string `json:"id"`
	Location string `json:"location"`
}

// connection struct
type connection struct {
	ID       string `json:"id"`
	SourceId string `json:"sourceId"`
	TargetId string `json:"targetId"`
}

// ifcondition struct
type ifcondition struct {
	False string `json:"false"`
	ID    string `json:"id"`
	True  string `json:"true"`
}

//switchStuct struct
type switchStuct struct {
	ID          string `json:"id"`
	SwitchState string `json:"switchState"`
}

//foreachStuct struct
type foreachStuct struct {
	ID           string `json:"id"`
	ForloopState string `json:"forloopState"`
}

//caseStuct struct
type caseStuct struct {
	ID        string `json:"id"`
	CaseState string `json:"caseState"`
}

// while struct for while controller
type whileStuct struct {
	ID         string `json:"id"`
	WhileState string `json:"whileState"`
}

//sub variable struct
type subvariable struct {
	DataType string `json:"datatype"`
	Type     string `json:"type"`
	Valid    bool   `json:"valid"`
	Value    string `json:"value"`
}

//variable struct
type variable struct {
	Key       string        `json:"Key"`
	Category  string        `json:"Category"`
	Type      string        `json:"Type"`
	Value     string        `json:"Value"`
	ValueType string        `json:"ValueType"`
	ValueList []subvariable `json:"ValueList"`
	Priority  string        `json:"Priority"`
	Group     string        `json:"Group"`
	DataType  string        `json:"DataType"`
	ConvertTo string        `json:"ConvertTo"`
	IsValid   bool          `json:"IsValid"`
}

// varable struct
type OtherData struct {
	TrueStateUUID  string `json:"TrueStateUUID"`
	FalseStateUUID string `json:"FalseStateUUID"`
	ForeachUUID    string `json:"TrueStateUUID"`
	SwitchUUID     string `json:"SwitchUUID"`
	CaseUUID       string `json:"CaseUUID"`
	DefaultUUID    string `json:"DefaultUUID"`
	WhileUUID      string `json:"WhileUUID"`
	Name           string `json:"Name"`
	Email          string `json:"Email"`
	MobileNo       string `json:"MobileNo"`
	Company        string `json:"Company"`
}

// nodedata struct
type nodedata struct {
	SourceEndpoints     []SourceEndpoint `json:"SourceEndpoints"`
	BlockId             string           `json:"blockId"`
	Category            string           `json:"Category"`
	ControlEditDisabled bool             `json:"ControlEditDisabled"`
	Description         string           `json:"Description"`
	Icon                string           `json:"Icon"`
	Name                string           `json:"Name"`
	TargetEndpoints     []interface{}    `json:"TargetEndpoints"`
	Type                string           `json:"Type"`
	Variables           []variable       `json:"Variables"`
	LibraryID           string           `json:"library_id"`
	SchemaID            string           `json:"schema_id"`
	ParentView          string           `json:"parentView"`
	DisplayName         string           `json:"DisplayName"`
	OtherData           OtherData        `json:"OtherData"`
}

// jsonflow stuct
type JsonFlow struct {
	FlowName        string `json:"FlowName"`
	SessionID       string `json:"SessionID"`
	ResponseMessage chan *context.FlowResult
	Connections     []connection   `json:"connections"`
	Ifconditions    []ifcondition  `json:"ifconditions"`
	Nodes           []nodedata     `json:"nodes"`
	Arguments       []variable     `json:"arguments"`
	Switches        []switchStuct  `json:"switchs"`
	Forloops        []foreachStuct `json:"forloops"`
	Cases           []caseStuct    `json:"cases"`
	Views           []string       `json:"views"`
	WhileLoops      []whileStuct   `json:"whileloops"`
	Port            string         `json:"Port"`
	OSCode          string         `json:"OSCode"`
	SysArch         string         `json:"SysArch"`
}

// hibernate wf struct
type HibernatedWF struct {
	DateTime       string                 `json:"DateTime"`
	ExecutionLevel string                 `json:"ExecutionLevel"`
	FlowData       map[string]interface{} `json:"FlowData"`
	SessionID      string                 `json:"SessionID"`
	WFName         string                 `json:"WFName"`
}

/*type VersionInfo struct {
Version   string   `json:"Version"`
Date      string   `json:"Date"`
Changelog []string `json:"Changelog"`
}*/

// version response struct
type VersionResponse struct {
	IsUptodate       bool   `json:"IsUptodate"`
	PEVersionDetails string `json:"PEVersionDetails"`
	SFVersionDetails string `json:"SFVersionDetails"`
	Engine_Details   map[string]interface{}
}

// sessiontrandedetails struct
type SessionTranDetails struct {
	SessionID      string `json:"SessionID"`
	SessionType    string `json:"SessionType"`
	SessionDetails string `json:"SessionDetails"`
	Message        string `json:"Message"`
}

// test workflow invoke struct
type TestWorkflowInvoke struct {
	SessionID   string `json:"SessionID"`
	InArguments string `json:"InArguments"`
}

// package details to insall on the server will be send with this
type PackageInstaller struct {
	Content []string `json:"Content"`
}

// download executable struct
type ExecutableFlow struct {
	FlowName  string `json:"FlowName"`
	SessionID string `json:"SessionID"`
	Port      string `json:"Port"`
}

// publish the wf to a docker with this file
type DockerDeployement struct {
	SessionID string `json:"SessionID"`
	WFName    string `json:"WFName"`
	Port      string `json:"Port"`
	RAM       string `json:"RAM"`
	CPU       string `json:"CPU"`
}
