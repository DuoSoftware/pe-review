package logger

import (
	"duov6.com/common"
	"encoding/json"
	"strconv"
	"sync"
)

var DefaultLogstashUrl string
var DefaultObjectStoreUrl string
var DefaultAuthUrl string

type DomainLogConfig struct {
	Domain      string
	DiskLogs    bool
	LogStash    bool
	LogStashUrl string
}

var domainConfig map[string]DomainLogConfig
var domainConfigLock = sync.RWMutex{}

func ReadDomainConfigMap(domain string) (config DomainLogConfig) {
	domainConfigLock.RLock()
	defer domainConfigLock.RUnlock()
	config = domainConfig[domain]
	return
}

func SetDomainConfigMap(domain string, config DomainLogConfig) {
	domainConfigLock.Lock()
	defer domainConfigLock.Unlock()
	domainConfig[domain] = config
}

func FetchDomainLogConfigurations(domain string) DomainLogConfig {
	config := DomainLogConfig{}
	//url := "http://dev.smoothflow.io:3000/com.duosoftware.logs/logsettings/" + domain + "?securityToken=ignore"
	url := "http://" + DefaultObjectStoreUrl + "/com.duosoftware.logs/logsettings/" + domain + "?securityToken=ignore"
	err, bodyBytes := common.HTTP_GET(url, nil, false)
	if err == nil && len(bodyBytes) > 4 {
		_ = json.Unmarshal(bodyBytes, &config)
	}
	return config
}

func UploadDomainLogConfigurations(domain string, config DomainLogConfig) (err error) {
	//url := "http://dev.smoothflow.io:3000/com.duosoftware.logs/logsettings?securityToken=ignore"
	url := "http://" + DefaultObjectStoreUrl + "/com.duosoftware.logs/logsettings?securityToken=ignore"

	payload := `{"Object":{"Domain":"` + config.Domain + `", 
	"DiskLogs":` + strconv.FormatBool(config.DiskLogs) + `, "LogStash":` + strconv.FormatBool(config.LogStash) + `,
	 "LogStashUrl":"` + config.LogStashUrl + `"}, "Parameters":{"KeyProperty":"Domain"}}`

	err, _ = common.HTTP_POST(url, nil, []byte(payload), false)
	return
}

func GetDomainLogConfig(sessionID string) DomainLogConfig {
	domain := GetDomainBySessionID(sessionID)
	config := DomainLogConfig{}

	//first read from the Map
	config = ReadDomainConfigMap(domain)
	if config == (DomainLogConfig{}) {
		//if the config is empty read from the objectstore
		config = FetchDomainLogConfigurations(domain)
		//set config to map
		SetDomainConfigMap(domain, config)
	}
	return config
}
