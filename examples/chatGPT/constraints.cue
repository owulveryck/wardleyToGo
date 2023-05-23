Host: string | *"AUTO"

#AIPlugin: {
	// SchemaVersion Manifest schema version - required - v1
	schema_version: string & "v1"

	// NameForHuman Human-readable name, such as the full company name. 20 character max. - required
	name_for_human: string & =~"^.{1,20}$"

	// NameForModel Name the model will use to target the plugin (no spaces allowed, only letters and numbers). 50 character max. - required
	name_for_model: string & =~"^[a-zA-Z0-9]{1,50}$"

	// DescriptionForHuman Human-readable description of the plugin. 100 character max.
	description_for_human: string & =~"^.{1,100}$"

	// DescriptionForModel Description better tailored to the model, such as token context length considerations or keyword usage for improved plugin prompting. 8,000 character max. - required
	description_for_model: string & =~"^.{20,1000}$"
	auth:                  #Auth
	api:                   #API
	logo_url:              string | *"\(Host)/logo.png"
	contact_email:         string & =~"^.*@.*$"
	legal_info_url:        string | *"\(Host)/legal"
}

#Auth: {
	type: string | *"none"
}

#API: {
	type:                  string | *"openapi"
	url:                   string | *"\(Host)/openapi.yaml"
	is_user_authenticated: bool | *false
}
