package sevrest

import (
	"encoding/json"
	"fmt"

	"github.com/sevone/gorest"
)

type ObjectType struct {
	ID uint `json:"id,omitempty"`
	PluginID uint `json:"pluginId"`
	ParentObjectTypeId uint `json:"parentObjectTypeId"`
	Name string `json:"name"`
	IsEnabled bool `json:"isEnabled"`
	IsEditable bool `json:"isEditable"`
	ExtendedInfo json.RawMessage `json:"extendedInfo,omitempty"`
}

type IndicatorType struct {
	ID uint `json:"id,omitempty"`
	PluginID uint `json:"pluginId"`
	ObjectTypeID uint `json:"pluginObjectTypeId"`
	Name string `json:"name"`
	IsEnabled bool `json:"isEnabled"`
	IsDefault bool `json:"isDefault"`
	// TODO:  This can only be COUNTER32, COUNTER64, or GAUGE; it should
	//    probably be an enum
	Format string `json:"format"`
	DataUnits string `json:"dataUnits"`
	DisplayUnits string `json:"displayUnits"`
	Description string `json:"description"`
	AllowMaxValue bool `json:"allowMaximumValue"`
	SyntheticExpression string `json:"syntheticExpression"`
	SyntheticMaximumExpression string `json:"syntheticMaximumExpression"`
	ExtendedInfo json.RawMessage `json:"extendedInfo.omitempty"`
}

func (this *SevRest) GetIndicatorTypes(include_extended_info bool, filter map[string]interface{}) ([]IndicatorType, error) {
	// TODO:  Loop through pages
	page := 0
	size := 50

	var response *gorest.Response
	var err error
	if(filter == nil) {
		response, err = this.Rest.Get(fmt.Sprintf("plugins/indicatortypes?page=%d&size=%d&includeExtendedInfo=%t", page, size, include_extended_info))
	} else {
		response, err = this.Rest.Post(fmt.Sprintf("plugins/indicatortypes/filter?page=%d&size=%d&includeExtendedInfo=%t", page, size, include_extended_info), filter)
	}

	if(err != nil) {
		return nil, err
	}

	var response_data SearchResponse
	err = response.Decode(&response_data)
	if(err != nil) {
		return nil, err
	}

	var array []IndicatorType
	err = json.Unmarshal(response_data.Content, &array)
	if(err != nil) {
		return nil, err
	}

	return array, nil
}

// Sane defaults:  include_extended_info = false, filter = nil
func (this *SevRest) GetObjectTypes(include_extended_info bool, filter map[string]interface{}) ([]ObjectType, error) {
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

	var response_data SearchResponse
	err = response.Decode(&response_data)
	if(err != nil) {
		return nil, err
	}

	var array []ObjectType
	err = json.Unmarshal(response_data.Content, &array)
	if(err != nil) {
		return nil, err
	}

	return array, nil
}

