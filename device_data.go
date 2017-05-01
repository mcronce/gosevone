package sevrest

import (
	"errors"
	"fmt"
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
	CreateTypesAutomatically bool `json:"-"`
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
	Description string `json:"-"`
}

func (this *SevRest) PostDeviceData(device *DeviceData) (*string, error) {
	response, err := this.Rest.Post("devices/data", device)
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
func (this *DeviceData) NewObject(name string, type_name string, plugin_name string, create_automatically bool) (uint, *DeviceDataObject) {
	object := DeviceDataObject{
		Name : name,
		Type : type_name,
		PluginID : 0,
		PluginName : plugin_name,
		CreateAutomatically : create_automatically,
		Timestamps : make([]DeviceDataTimestamp, 0),
		TimestampMap : make(map[uint]uint),
	}
	id := uint(len(this.Objects))
	this.ObjectMap[name] = id
	this.Objects = append(this.Objects, object)
	return id, &this.Objects[id]
}

func (this *DeviceData) AddIndicator(object_name string, object_type string, plugin_name string, time uint, indicator_name string, value float64) (uint, uint, uint, *DeviceDataIndicator) {
	var object *DeviceDataObject

	id, exists := this.ObjectMap[object_name]
	if(exists) {
		object = &this.Objects[id]
	} else {
		id, object = this.NewObject(object_name, object_type, plugin_name, this.CreateAutomatically)
	}

	timestamp_id, indicator_id, indicator := object.AddIndicator(time, indicator_name, value)
	return id, timestamp_id, indicator_id, indicator
}

func (this *DeviceData) ResolveTimestamps() {
	for _, o := range this.Objects {
		for _, t := range o.Timestamps {
			if(t.Time < this.OldestTimestamp || this.OldestTimestamp == 0) {
				this.OldestTimestamp = t.Time
			}
			if(t.Time > this.LatestTimestamp) {
				this.LatestTimestamp = t.Time
			}
		}
	}
}

func (this *DeviceData) ResolvePluginIDs(api *SevRest) error {
	ids := make(map[string]uint)

	for _, object := range this.Objects {
		if(object.PluginID != 0) {
			continue
		}

		_, exists := ids[object.PluginName]
		if(!exists) {
			ids[object.PluginName] = 0
		}
	}

	for name, _ := range ids {
		filter := map[string]string{"objectName" : name}
		plugins, err := api.GetPlugins(filter)
		if(err != nil) {
			return err
		}
		if(len(plugins) == 0) {
			return errors.New(fmt.Sprintf("Plugin \"%s\" not found", name))
		}
		ids[name] = plugins[0].ID
	}

	for i, object := range this.Objects {
		if(object.PluginID != 0) {
			continue
		}

		id, _ := ids[object.PluginName]
		this.Objects[i].PluginID = id
	}

	return nil
}

func (this *DeviceData) CreateMissingTypes(api *SevRest) error {
	type_tree := make(map[string]*ObjectType)
	for _, object := range this.Objects {
		object_type, exists := type_tree[object.Type]
		if(!exists) {
			object_type = &ObjectType{
				PluginID : object.PluginID,
				// TODO:  ParentObjectTypeID?
				Name : object.Type,
				IsEnabled : true,
				IsEditable : false,
				// TODO:  ExtendedInfo?
				IndicatorTypes : make([]IndicatorType, 0),
				IndicatorTypeMap : make(map[string]uint),
			}
			type_tree[object.Type] = object_type
		}
		for _, timestamp := range object.Timestamps {
			for _, indicator := range timestamp.Indicators {
				object_type.AddIndicatorType(indicator.Name, true, true, indicator.Format, indicator.Units, indicator.Units, indicator.Description, true)
			}
		}
	}

	for name, object_type := range type_tree {
		filter := map[string]interface{}{
			"name" : name,
			"pluginId" : object_type.PluginID,
		}
		existing_types, err := api.GetObjectTypes(false, filter)
		if(err != nil) {
			return err
		}
		if(len(existing_types) == 0) {
			api.CreateObjectType(object_type)
		} else {
			// The object type already exists; make sure all the indicator
			//    types exist as well.
			existing_object_type := &existing_types[0]
			for _, indicator_type := range object_type.IndicatorTypes {
				_, exists := object_type.IndicatorTypeMap[indicator_type.Name]
				if(!exists) {
					indicator_type.ObjectTypeID = existing_object_type.ID
					api.CreateIndicatorType(&indicator_type)
				}
			}
		}
	}

	return nil
}

func (this *DeviceData) Post(api *SevRest) (*string, error) {
	this.ResolveTimestamps()
	if(this.CreateTypesAutomatically) {
		err := this.ResolvePluginIDs(api)
		if(err != nil) {
			return nil, err
		}

		err = this.CreateMissingTypes(api)
		if(err != nil) {
			return nil, err
		}
	}
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

func (this *DeviceDataObject) AddIndicator(time uint, name string, value float64) (uint, uint, *DeviceDataIndicator) {
	var timestamp *DeviceDataTimestamp

	id, exists := this.TimestampMap[time]
	if(exists) {
		timestamp = &this.Timestamps[id]
	} else {
		id, timestamp = this.NewTimestamp(time)
	}

	indicator_id, indicator := timestamp.AddIndicator(name, value)
	return id, indicator_id, indicator
}

// TODO:  More args?
func (this *DeviceDataTimestamp) NewIndicator(name string, value float64) (uint, *DeviceDataIndicator) {
	indicator := DeviceDataIndicator{
		Name : name,
		Description : name,
		Value : value,
		Format : "GAUGE",
	}
	id := uint(len(this.Indicators))
	this.IndicatorMap[name] = id
	this.Indicators = append(this.Indicators, indicator)
	return id, &this.Indicators[id]
}

func (this *DeviceDataTimestamp) AddIndicator(name string, value float64) (uint, *DeviceDataIndicator) {
	id, exists := this.IndicatorMap[name]
	if(exists) {
		indicator := &this.Indicators[id]
		indicator.Value = value
		return id, indicator
	}
	return this.NewIndicator(name, value)
}

