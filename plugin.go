package sevrest

import (
	"encoding/json"
	"fmt"

	"github.com/sevone/gorest"
)

type Plugin struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	ShortName string `json:"objectName"`
	Dir string `json:"dir"`
	Plottable int `json:"plottable"`
}

type ObjectType struct {
	ID uint `json:"id,omitempty"`
	PluginID uint `json:"pluginId"`
	ParentObjectTypeId uint `json:"parentObjectTypeId"`
	Name string `json:"name"`
	IsEnabled bool `json:"isEnabled"`
	IsEditable bool `json:"isEditable"`
	ExtendedInfo json.RawMessage `json:"extendedInfo,omitempty"`
	IndicatorTypes []IndicatorType `json:"-"`
	IndicatorTypeMap map[string]uint `json:"-"`
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
	ExtendedInfo json.RawMessage `json:"extendedInfo,omitempty"`
}

func (this *SevRest) GetPlugins(filter map[string]string) ([]Plugin, error) /* {{{ */ {
	page := 0
	size := 50

	filter_str := fmt.Sprintf("page=%d&size=%d", page, size)
	value, exists := filter["objectName"]
	if(exists) {
		filter_str += fmt.Sprintf("&objectName=%s", value)
	}
	value, exists = filter["name"]
	if(exists) {
		filter_str += fmt.Sprintf("&name=%s", value)
	}

	response, err := this.Rest.Get(fmt.Sprintf("plugins?%s", filter_str))
	if(err != nil) {
		return nil, err
	}

	var response_data SearchResponse
	err = response.Decode(&response_data)
	if(err != nil) {
		return nil, err
	}

	var array []Plugin
	err = json.Unmarshal(response_data.Content, &array)
	if(err != nil) {
		return nil, err
	}

	return array, nil
} // }}}

func (this *SevRest) GetIndicatorTypes(include_extended_info bool, filter map[string]interface{}) ([]IndicatorType, error) /* {{{ */ {
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
} // }}}

// Sane defaults:  include_extended_info = false, filter = nil
func (this *SevRest) GetObjectTypes(include_extended_info bool, filter map[string]interface{}) ([]ObjectType, error) /* {{{ */ {
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
} // }}}

func (this *SevRest) GetObjectTypeExtendedInfo(plugin uint) (json.RawMessage, error) /* {{{ */ {
	response, err := this.Rest.Get(fmt.Sprintf("plugins/objecttypes/schema/%d", plugin))
	if(err != nil) {
		return nil, err
	}

	var response_data map[string]string
	err = response.Decode(&response_data)
	if(err != nil) {
		return nil, err
	}

	output_map := make(map[string]string)
	for k, _ := range response_data {
		output_map[k] = ""
	}

	return json.Marshal(output_map)
} // }}}

func (this *SevRest) GetIndicatorTypeExtendedInfo(plugin uint) (json.RawMessage, error) /* {{{ */ {
	response, err := this.Rest.Get(fmt.Sprintf("plugins/indicatortypes/schema/%d", plugin))
	if(err != nil) {
		return nil, err
	}

	var response_data map[string]string
	err = response.Decode(&response_data)
	if(err != nil) {
		return nil, err
	}

	output_map := make(map[string]string)
	for k, _ := range response_data {
		output_map[k] = ""
	}

	return json.Marshal(output_map)
} // }}}

func (this *SevRest) CreateIndicatorType(payload *IndicatorType) (uint, error) /* {{{ */ {
	ext, err := this.GetIndicatorTypeExtendedInfo(payload.PluginID)
	if(err != nil) {
		return 0, err
	}
	payload.ExtendedInfo = ext

	response, err := this.Rest.Post("plugins/indicatortypes", payload)
	if(err != nil) {
		return 0, err
	}

	var body IndicatorType
	err = response.Decode(&body)
	if(err != nil) {
		return 0, err
	}

	payload.ID = body.ID
	return payload.ID, nil
} // }}}

func (this *SevRest) CreateObjectType(payload *ObjectType) (uint, []uint, error) /* {{{ */ {
	ext, err := this.GetObjectTypeExtendedInfo(payload.PluginID)
	if(err != nil) {
		return 0, nil, err
	}
	payload.ExtendedInfo = ext

	response, err := this.Rest.Post("plugins/objecttypes", payload)
	if(err != nil) {
		return 0, nil, err
	}

	var body ObjectType
	err = response.Decode(&body)
	if(err != nil) {
		return 0, nil, err
	}

	payload.ID = body.ID
	created_indicator_type_ids := make([]uint, 0)
	for i, _ := range payload.IndicatorTypes {
		payload.IndicatorTypes[i].ObjectTypeID = payload.ID
		id, err := this.CreateIndicatorType(&payload.IndicatorTypes[i])
		if(err != nil) {
			return 0, nil, err
		}
		created_indicator_type_ids = append(created_indicator_type_ids, id)
	}

	return payload.ID, created_indicator_type_ids, nil
} // }}}

func (this *ObjectType) NewIndicatorType(name string, is_enabled bool, is_default bool, format string, data_units string, display_units string, description string, allow_max bool) (uint, *IndicatorType) /* {{{ */ {
	indicator_type := IndicatorType{
		PluginID : this.PluginID,
		ObjectTypeID : this.ID,
		Name : name,
		IsEnabled : is_enabled,
		IsDefault : is_default,
		Format : format,
		DataUnits : data_units,
		DisplayUnits : display_units,
		Description : description,
		AllowMaxValue : allow_max,
	}
	id := uint(len(this.IndicatorTypes))
	this.IndicatorTypeMap[name] = id
	this.IndicatorTypes = append(this.IndicatorTypes, indicator_type)
	return id, &this.IndicatorTypes[id]
} // }}}

func (this *ObjectType) AddIndicatorType(name string, is_enabled bool, is_default bool, format string, data_units string, display_units string, description string, allow_max bool) /* {{{ */ {
	id, exists := this.IndicatorTypeMap[name]
	if(exists) {
		indicator_type := &this.IndicatorTypes[id]
		indicator_type.PluginID = this.PluginID
		indicator_type.ObjectTypeID = this.ID
		indicator_type.Name = name
		indicator_type.IsEnabled = is_enabled
		indicator_type.IsDefault = is_default
		indicator_type.Format = format
		indicator_type.DataUnits = data_units
		indicator_type.DisplayUnits = display_units
		indicator_type.Description = description
		indicator_type.AllowMaxValue = allow_max
	} else {
		this.NewIndicatorType(name, is_enabled, is_default, format, data_units, display_units, description, allow_max)
	}
} // }}}

