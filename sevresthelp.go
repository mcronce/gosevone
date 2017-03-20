package sevrest

import (
	"fmt"
)

type SevRestResponse struct {
	Description string `json: description`
	Schema map[string]string `json: schema`
}

type SevRestParameter struct {
	Name string `json: name`
	Description string `json: description`
	Type string `json: type`
	Format string `json: format`
	Required string `json: required`
	Enum []string `json: enum`
	Schema map[string] string `json: schema`
}

type SevRestPath struct {
	Description string `json: description`
	Parameters []SevRestParameter `json: parameters`
	Responses map[string]SevRestResponse
	Summary string `json: summary`
	Tags []string `json: tags`
}

type SevRestDefinitionProperties struct {
	Type string `json: type`
	Enum []string `json: enum`
	Ref string `json: $ref`
	Items map[string]string `json: items`
}

type SevRestDefinition struct {
	Properties map[string]SevRestDefinitionProperties `json: properties`
}

type SevRestApiDocs struct {
	Paths map[string]map[string]SevRestPath `json: paths`
	Definitions map[string]SevRestDefinition `json: definitions`
}


func (apiDocs SevRestApiDocs) PrintSchemaDefinition(ref string, prepend string) {
	const definitionsRef = "#/definitions/"

	if(ref == "") { return }
	ref = ref[len(definitionsRef):]

	fmt.Printf("%s%s :\n", prepend, ref)
	for k, v := range apiDocs.Definitions[ref].Properties {
		if v.Ref != "" {
			apiDocs.PrintSchemaDefinition(v.Ref, prepend+"    ")
		} else if v.Type == "array" {
			//PrettyPrint(v)
			fmt.Printf("  %s%s: [\n", prepend, k)
			apiDocs.PrintSchemaDefinition(v.Items["$ref"], prepend+"    ")
			fmt.Printf("  %s]\n", prepend)
		} else {
			fmt.Printf("  %s%s: %s\n", prepend, k, v.Type)
		}
	}
	//fmt.Printf("%s]\n", prepend)
}

func (apiDocs SevRestApiDocs) PrintSchema(schema map[string]string, prepend string) {
	for k, v := range schema {
		switch k {
		case "$ref":
			apiDocs.PrintSchemaDefinition(v, prepend)
		default:
			fmt.Printf("%sUNHANDLED: %s\n", prepend, k)
		}
	}
}
