package sevrest

import (
	"fmt"

	"github.com/mcronce/gorest"
)

type PluginObjectType struct {
	ID uint `json:"id,omitempty"`
	PluginID uint `json:"pluginId"`
	ParentObjectTypeId uint `json:"parentObjectTypeId"`
	Name string `json:"name"`
	IsEnabled bool `json:"isEnabled"`
	IsEditable bool `json:"isEditable"`
	ExtendedInfo map[string]interface{} `json:"extendedInfo"`
}

// Sane defaults:  include_extended_info = false, filter = nil
func (this *SevRest) GetPluginObjectTypes(include_extended_info bool, filter map[string]interface{}) ([]PluginObjectType, error) {
	// TODO:  Loop through pages
	page := 0
	size := 50

	var response *gorest.Response
	var err error
	if(filter == nil) {
		response, err = this.Rest.Get(fmt.Sprintf("plugins/objecttypes?page=%d&size=%d&includeExtendedInfo=%t", page, size, include_extended_info))
	} else {
		response, err = this.Rest.Post(fmt.Sprintf("plugins/objecttypes/filter?page=%d&size=%d&includeExtendedInfo=%t", page, size, include_extended_info), filter)
	}

	if(err != nil) {
		return nil, err
	}

	var response_data map[string]interface{}
	err = response.Decode(&response_data)
	if(err != nil) {
		return nil, err
	}

	PrettyPrint(response_data)
	return nil, nil
}
