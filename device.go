package sevrest

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sevone/gorest"
)

type Device struct {
	ID uint `json:"id,omitempty"`
	PeerID uint `json:"peerId,omitempty"`
	PluginManagerID uint `json:"pluginManagerId,omitempty"`
	WorkHoursGroupId uint `json:"workhoursGroupId,omitempty"`
	Name string `json:"name"`
	IP string `json:"ipAddress"`
	AltName string `json:"alternateName,omitempty"`
	NumElements uint `json:"numElements,omitempty"`
	DateAdded uint `json:"dateAdded,omitempty"`
	LastDiscovery uint `json:"lastDiscovery,omitempty"`
	PollFrequency uint `json:"pollFrequency"`
	IsNew bool `json:"isNew,omitempty"`
	IsDeleted bool `json:"isDeleted,omitempty"`
	Description string `json:"description"`
	Timezone string `json:"timezone,omitempty"`
	AllowDelete bool `json:"allowDelete,omitempty"`
	DisablePolling bool `json:"disablePolling,omitempty"`
	DisableConcurrentPolling bool `json:"disableConcurrentPolling,omitempty"`
	DisableThresholding bool `json:"disableThresholding,omitempty"`
	Objects []DeviceObject `json:"objects,omitempty"`
	// TODO:  Implement this
	PluginInfo map[string]interface{} `json:"pluginInfo,omitempty"`
}

type DeviceObject struct {
	ID uint `json:"id,omitempty"`
	PluginID uint `json:"pluginId"`
	ObjectTypeID uint `json:"pluginObjectTypeId"`
	SubtypeID uint `json:"subtypeId"`
	DeviceID uint `json:"deviceId,omitempty"`
	Name string `json:"name"`
	IsEnabled bool `json:"isEnabled"`
	IsVisible bool `json:"isVisible,omitempty"`
	IsDeleted bool `json:"isDeleted,omitempty"`
	Description string `json:"description"`
	Indicators []DeviceIndicator `json:"indicators"`
	ExtendedInfo json.RawMessage `json:"extendedInfo,omitempty"`
}

type DeviceIndicator struct {
	ID uint `json:"id,omitempty"`
	PluginID uint `json:"pluginId"`
	IndicatorTypeID uint `json:"pluginIndicatorTypeId"`
	DeviceID uint `json:"deviceId,omitempty"`
	ObjectID uint `json:"objectId,omitempty"`
	IsEnabled bool `json:"isEnabled"`
	IsDeleted bool `json:"isDeleted,omitempty"`
	IsBaselining bool `json:"isBaselining,omitempty"`
	// TODO:  This can only be COUNTER32, COUNTER64, or GAUGE; it should
	//    probably be an enum
	Format string `json:"format"`
	MaxValue int `json:"maxValue"`
	LastInvalidationTime uint `json:"lastInvalidationTime,omitempty"`
	EvaluationOrder int `json:"evaluationOrder,omitempty"`
	SyntheticExpression string `json:"syntheticExpression,omitempty"`
	ExtendedInfo json.RawMessage `json:"extendedInfo,omitempty"`
}

func (this *SevRest) GetDevices(filter map[string]interface{}) ([]Device, error) {
	// TODO:  Loop through pages
	page := 0
	size := 50

	var response *gorest.Response
	var err error
	if(filter == nil) {
		response, err = this.Rest.Get(fmt.Sprintf("devices?page=%d&size=%d", page, size))
	} else {
		response, err = this.Rest.Post(fmt.Sprintf("devices/filter?page=%d&size=%d", page, size), filter)
	}

	if(err != nil) {
		return nil, err
	}

	var response_data SearchResponse
	err = response.Decode(&response_data)
	if(err != nil) {
		return nil, err
	}

	var array []Device
	err = json.Unmarshal(response_data.Content, &array)
	if(err != nil) {
		return nil, err
	}

	return array, nil
}

func (this *SevRest) GetDeviceObjects(include_indicators bool, include_extended_info bool, filter map[string]interface{}) ([]DeviceObject, error) {
	if(len(filter) == 0) {
		return nil, errors.New("GetDeviceObjects():  filter cannot be empty; recommend passing deviceId")
	}

	response, err := this.Rest.Post(fmt.Sprintf("devices/objects/filter?includeIndicators=%t&includeExtendedInfo=%t", include_indicators, include_extended_info), filter)
	if(err != nil) {
		return nil, err
	}

	var response_data []DeviceObject
	err = response.Decode(&response_data)
	if(err != nil) {
		return nil, err
	}

	return response_data, nil
}

