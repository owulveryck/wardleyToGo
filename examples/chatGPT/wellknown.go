package main

import (
	"io/ioutil"

	"cuelang.org/go/cue/cuecontext"
)

type AIPlugin struct {
	// SchemaVersion Manifest schema version - required - v1
	SchemaVersion string `json:"schema_version"`
	// NameForHuman Human-readable name, such as the full company name. 20 character max. - required
	NameForHuman string `json:"name_for_human"`
	// NameForModel Name the model will use to target the plugin (no spaces allowed, only letters and numbers). 50 character max. - required
	NameForModel string `json:"name_for_model"`
	// DescriptionForHuman Human-readable description of the plugin. 100 character max.
	DescriptionForHuman string `json:"description_for_human"`
	// DescriptionForModel Description better tailored to the model, such as token context length considerations or keyword usage for improved plugin prompting. 8,000 character max. - required
	DescriptionForModel string `json:"description_for_model"`
	Auth                Auth   `json:"auth"`
	API                 API    `json:"api"`
	LogoURL             string `json:"logo_url"`
	ContactEmail        string `json:"contact_email"`
	LegalInfoURL        string `json:"legal_info_url"`
	BaseURL             string `json:"-"`
}

func (a *AIPlugin) Check() error {
	return nil
}

type Auth struct {
	Type string `json:"type"`
}

type API struct {
	Type                string `json:"type"`
	URL                 string `json:"url"`
	IsUserAuthenticated bool   `json:"is_user_authenticated"`
}

// NewAIPlugin generates an AI plugin from the cue files
// Host is substituated by Host
func NewAIPlugin(host string) (*AIPlugin, error) {
	host = `Host: "` + host + `"`
	constraints, err := ioutil.ReadFile("constraints.cue")
	if err != nil {
		return nil, err
	}
	configuration, err := ioutil.ReadFile("wellknown.cue")
	if err != nil {
		return nil, err
	}
	content := append(constraints, configuration...)
	content = append([]byte(host+"\n"), content...)
	ctx := cuecontext.New()
	v := ctx.CompileBytes(content)
	v = v.Lookup("configuration")
	var aiplugin AIPlugin
	err = v.Decode(&aiplugin)
	if err != nil {
		return nil, err
	}
	return &aiplugin, nil
}
