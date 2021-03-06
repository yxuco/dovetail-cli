package model

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-lib/app"
)

type ModelResources struct {
	AppName      string
	Assets       []string
	Transactions []string
	Schemas      map[string]string
}

type ResourceMetadataModel struct {
	Metadata struct {
		Type         string `json:"type"`
		Parent       string `json:"parent"`
		CordaClass   string `json:"cordaClass"`
		IdentifiedBy string `json:"identifiedBy"`
		IsAbstract   bool   `json:"isAbstract, omitempty"`
		Decorators   []struct {
			Name string   `json:"name"`
			Args []string `json:"args, omitempty"`
		} `json:"decorators, omitempty"`
	} `json:"metadata"`
	Attributes []ResourceAttribute `json:"attributes"`
}

type ResourceAttribute struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	IsRef      bool   `json:"isRef"`
	IsArray    bool   `json:"isArray"`
	IsOptional bool   `json:"isOptional"`
}

func ParseResourceModel(jsonResource string) *ResourceMetadataModel {
	resource := ResourceMetadataModel{}

	json.Unmarshal([]byte(jsonResource), &resource)
	return &resource
}

func ParseFlowApp(jsonFile string) (*ModelResources, error) {
	appCfg, err := ParseApp(jsonFile)
	if err != nil {
		return nil, err
	}

	model := ModelResources{}
	model.AppName = appCfg.Name
	model.Schemas = make(map[string]string)
	json.Unmarshal([]byte(appCfg.Triggers[0].GetSetting("assets")), &model.Assets)
	json.Unmarshal([]byte(appCfg.Triggers[0].GetSetting("transactions")), &model.Transactions)

	schemas := struct{ schemas [][]string }{}
	json.Unmarshal([]byte(appCfg.Triggers[0].GetSetting("schemas")), &schemas.schemas)

	for _, value := range schemas.schemas {
		model.Schemas[value[0]] = value[1]
	}
	return &model, nil
}

// ParseApp parses the model file into an app.Config struct
func ParseApp(modelfile string) (*app.Config, error) {

	flowjson, err := ioutil.ReadFile(modelfile)
	if err != nil {
		return nil, err
	}

	return decodeApp(bytes.NewReader(flowjson))
}

// decodeApp decodes the model file into an app.Config struct
func decodeApp(r io.Reader) (*app.Config, error) {
	appCfg := &app.Config{}

	jsonParser := json.NewDecoder(r)
	err := jsonParser.Decode(&appCfg)
	if err != nil {
		return nil, err
	}

	return appCfg, nil
}
