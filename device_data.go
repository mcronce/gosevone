package sevrest

type DeviceData struct {
	Name string `json:"name"`
	Type string `json:"type"`
	OldestTimestamp uint `json:"oldTs"`
	LatestTimestamp uint `json:"newTs"`
	IP string `json:"ip"`
	CreateAutomatically bool `json:"automaticCreation,omitempty"`
	SourceID uint `json:"sourceId"`
	Objects []DeviceDataObject `json:"objects"`
	// Map of object names to indices
	ObjectMap map[string]int
}

type DeviceDataObject struct {
	Name string `json:"name"`
	Type string `json:"type"`
	PluginID uint `json:"pluginId,omitempty"`
	PluginName string `json:"pluginName,omitempty"`
	Description string `json:"description.omitempty"`
	CreateAutomatically bool `json:"automaticCreation,omitempty"`
	Timestamps []DeviceDataTimestamp `json:"timestamps"`
	// Map of times to indices
	TimestampMap map[uint]int
}

type DeviceDataTimestamp struct {
	Time uint `json:"timestamp"`
	Indicators []DeviceDataIndicator `json:"indicators"`
	// Map of indicator names to indices
	IndicatorMap map[string]int
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

// TODO:  More args?
func NewDeviceData(name string, initial_timestamp uint, source_id uint) DeviceData {
	return DeviceData{
		Name : name,
		Type : "Generic",
		OldestTimestamp : initial_timestamp,
		LatestTimestamp : initial_timestamp,
		SourceID : source_id,
		IP : "0.0.0.0",
		Objects : make([]DeviceDataObject, 0),
		ObjectMap : make(map[string]int),
	}
}

// TODO:  More args?
func (this *DeviceData) NewObject(name string, type_name string) (uint, *DeviceDataObject) {
	object := DeviceDataObject{
		Name : name,
		Type : type_name,
		PluginID : 17,
		PluginName : "BULKDATA",
		Timestamps : make([]DeviceDataTimestamp, 0),
	}
	id := len(this.Objects)
	this.ObjectMap[name] = id
	this.Objects = append(this.Objects, object)
	return id, &object
}

func (this *DeviceData) AddIndicator(object_name string, object_type string, time uint, indicator_name string, value float64) {
	var object *DeviceDataObject
	var object_id uint

	id, exists := this.ObjectMap[object_name]
	if(exists) {
		object = &this.Objects[id]
		object_id = id
	} else {
		object_id, object = this.NewObject(object_name, object_type)
	}

	object.AddIndicator(time, indicator_name, value)
}

// TODO:  More args?
func (this *DeviceDataObject) NewTimestamp(time uint) (uint, *DeviceDataTimestamp) {
	timestamp := DeviceDataTimestamp{
		Time : time,
		Indicators : make([]DeviceDataIndicator, 0),
		IndicatorMap : make(map[string]int),
	}
	id := len(this.Timestamps)
	this.TimestampMap[time] = id
	this.Timestamps = append(this.Timestamps, timestamp)
	return id, &timestamp
}

func (this *DeviceDataObject) AddIndicator(time uint, name string, value float64) {
	var timestamp *DeviceDataTimestamp
	var timestamp_id uint

	id, exists := this.TimestampMap[time]
	if(exists) {
		timestamp = &this.Timestamps[id]
		timestamp_id = id
	} else {
		timestamp_id, timestamp = this.NewTimestamp(time)
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
	id := len(this.Indicators)
	this.IndicatorMap[name] = id
	this.Indicators = append(this.Indicators, indicator)
	return id, &indicator
}

func (this *DeviceDataTimestamp) AddIndicator(name string, value float64) {
	var indicator *DeviceDataIndicator
	var indicator_id uint

	id, exists := this.IndicatorMap[name]
	if(exists) {
		indicator = &this.Indicators[id]
		indicator_id = id
	} else {
		indicator_id, indicator = this.NewIndicator(name, value)
	}
}

