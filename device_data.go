package sevrest

import (
	"encoding/json"
	"io/ioutil"
)

type DeviceData struct {
	Name string `json:"name"`
	Type string `json:"type"`
	OldestTimestamp uint `json:"oldTs"`
	LatestTimestamp uint `json:"newTs"`
	IP string `json:"ip"`
	CreateAutomatically bool `json:"automaticCreation,omitempty"`
	SourceID uint `json:"sourceId,omitempty"`
	Objects []DeviceDataObject `json:"objects"`
	// Map of object names to indices
	ObjectMap map[string]uint `json:"-"`
}

type DeviceDataObject struct {
	Name string `json:"name"`
	Type string `json:"type"`
	PluginID uint `json:"pluginId,omitempty"`
	PluginName string `json:"pluginName,omitempty"`
	Description string `json:"description,omitempty"`
	CreateAutomatically bool `json:"automaticCreation,omitempty"`
	Timestamps []DeviceDataTimestamp `json:"timestamps"`
	// Map of times to indices
	TimestampMap map[uint]uint `json:"-"`
}

type DeviceDataTimestamp struct {
	Time uint `json:"timestamp"`
	Indicators []DeviceDataIndicator `json:"indicators"`
	// Map of indicator names to indices
	IndicatorMap map[string]uint `json:"-"`
}

type DeviceDataIndicator struct {
	Name string `json:"name"`
	Value float64 `json:"value"`
	// TODO:  This can only be COUNTER32, COUNTER64, or GAUGE; it should
	//    probably be an enum
	Format string `json:"format"`
	Units string `json:"units,omitempty"`
	MaxValue float64 `json:"maxValue,omitempty"`
}

func (this *SevRest) PostDeviceData(device *DeviceData) (*string, error) {
	data, err := json.Marshal(device)
	if(err != nil) {
		return nil, err
	}

	response, err := this.Rest.Post("devices/data", data)
	if(err != nil) {
		return nil, err
	}

	body_raw, err := ioutil.ReadAll(response.Body)
	if(err != nil) {
		return nil, err
	}
	body := string(body_raw)

	return &body, nil
}

// TODO:  More args?
func NewDeviceData(name string, initial_timestamp uint, source_id uint) DeviceData {
	device := DeviceData{
		Name : name,
		Type : "Generic",
		OldestTimestamp : initial_timestamp,
		LatestTimestamp : initial_timestamp,
		SourceID : source_id,
		IP : "0.0.0.0",
		Objects : make([]DeviceDataObject, 0),
		ObjectMap : make(map[string]uint),
	}
	if(source_id == 0) {
		device.CreateAutomatically = true
	}
	return device
}

// TODO:  More args?
func (this *DeviceData) NewObject(name string, type_name string, create_automatically bool) (uint, *DeviceDataObject) {
	object := DeviceDataObject{
		Name : name,
		Type : type_name,
		PluginID : 17,
		PluginName : "BULKDATA",
		CreateAutomatically : create_automatically,
		Timestamps : make([]DeviceDataTimestamp, 0),
		TimestampMap : make(map[uint]uint),
	}
	id := uint(len(this.Objects))
	this.ObjectMap[name] = id
	this.Objects = append(this.Objects, object)
	return id, &this.Objects[id]
}

func (this *DeviceData) AddIndicator(object_name string, object_type string, time uint, indicator_name string, value float64) {
	var object *DeviceDataObject

	id, exists := this.ObjectMap[object_name]
	if(exists) {
		object = &this.Objects[id]
	} else {
		_, object = this.NewObject(object_name, object_type, this.CreateAutomatically)
	}

	object.AddIndicator(time, indicator_name, value)
}

func (this *DeviceData) ResolveTimestamps() {
	for _, o := range this.Objects {
		for _, t := range o.Timestamps {
			if(t.Time < this.OldestTimestamp) {
				this.OldestTimestamp = t.Time
			}
			if(t.Time > this.LatestTimestamp) {
				this.LatestTimestamp = t.Time
			}
		}
	}
}

func (this *DeviceData) Post(api *SevRest) (*string, error) {
	this.ResolveTimestamps()
	return api.PostDeviceData(this)
}

// TODO:  More args?
func (this *DeviceDataObject) NewTimestamp(time uint) (uint, *DeviceDataTimestamp) {
	timestamp := DeviceDataTimestamp{
		Time : time,
		Indicators : make([]DeviceDataIndicator, 0),
		IndicatorMap : make(map[string]uint),
	}
	id := uint(len(this.Timestamps))
	this.TimestampMap[time] = id
	this.Timestamps = append(this.Timestamps, timestamp)
	return id, &this.Timestamps[id]
}

func (this *DeviceDataObject) AddIndicator(time uint, name string, value float64) {
	var timestamp *DeviceDataTimestamp

	id, exists := this.TimestampMap[time]
	if(exists) {
		timestamp = &this.Timestamps[id]
	} else {
		_, timestamp = this.NewTimestamp(time)
	}

	timestamp.AddIndicator(name, value)
}

// TODO:  More args?
func (this *DeviceDataTimestamp) NewIndicator(name string, value float64) (uint, *DeviceDataIndicator) {
	indicator := DeviceDataIndicator{
		Name : name,
		Value : value,
		Format : "GAUGE",
	}
	id := uint(len(this.Indicators))
	this.IndicatorMap[name] = id
	this.Indicators = append(this.Indicators, indicator)
	return id, &this.Indicators[id]
}

func (this *DeviceDataTimestamp) AddIndicator(name string, value float64) {
	id, exists := this.IndicatorMap[name]
	if(exists) {
		indicator := &this.Indicators[id]
		indicator.Value = value
	} else {
		this.NewIndicator(name, value)
	}
}

