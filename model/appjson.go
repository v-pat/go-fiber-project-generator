package model

type AppJson struct {
	AppName string             `json:"appName"`
	Tables  []StructDefinition `json:"tables"`
}

// StructDefinition represents the data required for generating CRUD methods.
type StructDefinition struct {
	StructName  string                 `json:"name"`
	JSONExample map[string]interface{} `json:"columns"`
	Endpoint    string                 `json:"endpoint"`
}
