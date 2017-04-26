package Components

type ExternalPlugins struct {
	EmailAddress string //PK
	Jira         []Jira
}

type Jira struct {
	Name          string
	DisplayName   string
	TimeZone      string
	Locale        string
	IsDefault     bool
	JiraDomain    string
	SelfUrl       string
	SecurityToken string
	TokenType     string //Basic, Bearer
}
