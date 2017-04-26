package logger

import (
	"duov6.com/term"
	"github.com/fatih/color"
	"os"
	"processengine/Common"
	"time"
)

var IsLoggerInitiated bool
var IsLogstashEnabled bool

const (
	Information = 0
	Error       = 1
	Debug       = 2
	Splash      = 3
	Blank       = 4
	Warning     = 5

	ProcessEngine = 6
	Activity      = 7
	WorkFlow      = 8
	Default       = 9
)

func InitiateLogger() {
	IsLogstashEnabled = true
	IsLoggerInitiated = true
	//Initiate duo term
	term.GetConfig()
	//Load default logstash URL

	config := Common.VerifyCEBAgentConfig()

	if config["objUrl"] != nil {
		DefaultObjectStoreUrl = config["objUrl"].(string)
	}
	if config["authUrl"] != nil {
		DefaultAuthUrl = config["authUrl"].(string)
	}
	if config["logstashUrl"] != nil {
		DefaultLogstashUrl = config["logstashUrl"].(string)
	}

	fileHandlers = make(map[string]*os.File)
	domainConfig = make(map[string]DomainLogConfig)
	plans = make(map[string]UserPlan)
}

func ToggleLogs() string {
	msg := term.ToggleConfig()
	return msg
}

func ToggleLogstash() (msg string) {
	if IsLogstashEnabled {
		IsLogstashEnabled = false
		msg = "Disabled Logstash logging."
	} else {
		IsLogstashEnabled = true
		msg = "Enabled Logstash logging."
	}
	return
}

func Log(message interface{}, logType int, sessionID string, category int) {
	if category < 5 || category > 9 {
		category = Default
	}

	if !IsLoggerInitiated {
		//log to disk and return
		logNames := GetLogNamesByCategory(category, sessionID)
		for _, name := range logNames {
			PublishToDisk(name, GetMessageInString(message))
		}
		return
	}

	//Log to Terminal Output. filename is optional
	term.Write(message, logType)
	//Check if Logs are from Engine or others.. if engine print disk logs anyway
	if VerifyIsSmoothFlowEngine() {
		//Do disk logs
		logNames := GetLogNamesByCategory(category, sessionID)
		for _, name := range logNames {
			PublishToDisk(name, GetMessageInString(message))
		}
		//Do Logstashlogs
		if DefaultLogstashUrl != "" && IsLogstashEnabled == true {
			logstashUrl := "http://" + DefaultLogstashUrl + "/Logstash"
			PublishToLogstash(logstashUrl, sessionID, GetMessageInString(message), logType, category)
		}
	} else {
		//Work according to domain's specified settings
		config := GetDomainLogConfig(sessionID)
		if config.DiskLogs {
			//Write to disk
			logNames := GetLogNamesByCategory(category, sessionID)
			for _, name := range logNames {
				PublishToDisk(name, GetMessageInString(message))
			}
		}
		if config.LogStash {
			//Check for Plan Support
			if CheckForPlanSupport(sessionID) {
				if config.LogStashUrl != "" && IsLogstashEnabled != true {
					PublishToLogstash(config.LogStashUrl, sessionID, GetMessageInString(message), logType, category)
				} else {
					if DefaultLogstashUrl != "" && IsLogstashEnabled == true {
						PublishToLogstash(("http://" + DefaultLogstashUrl + "/Logstash"), sessionID, GetMessageInString(message), logType, category)
					}
				}
			}
		}
	}
}

func TermLog(message interface{}, logType int) {
	term.Write(message, logType)
}

func EventLog(message interface{}, logType int) {
	term.Write(message, logType)
	PublishToDisk("ProcessEngineLog.log", GetMessageInString(message))
}

func Log_PE(message interface{}, logType int, sessionID string) {
	Log(message, logType, sessionID, ProcessEngine)
}

func Log_WF(message interface{}, logType int, sessionID string) {
	Log(message, logType, sessionID, WorkFlow)
}

func Log_ACT(message interface{}, logType int, sessionID string) {
	Log(message, logType, sessionID, Activity)
}

func Log_Default(message interface{}, logType int, sessionID string) {
	Log(message, logType, sessionID, Default)
}

func HighLight(message, sessionID string) {
	Log((time.Now().Format("2006-01-02 15:04:05") + ":" + message), Debug, sessionID, ProcessEngine)
	font := color.New(color.FgBlack)
	colorscheme := font.Add(color.BgWhite)
	colorscheme.Println(message)
}
