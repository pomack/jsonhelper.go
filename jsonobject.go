package jsonhelper

import (
	"time"
)

type JSONObject map[string]interface{}

func NewJSONObject() JSONObject {
	return make(JSONObject)
}

func NewJSONObjectFromMap(m map[string]interface{}) JSONObject {
	return JSONObject(m)
}

func (p JSONObject) Del(key string) {
	p[key] = nil, false
}

func (p JSONObject) Set(key string, value interface{}) {
	p[key] = value
}

func (p JSONObject) Get(key string) interface{} {
	value, _ := p[key]
	return value
}

func (p JSONObject) GetAsString(key string) string {
	value, _ := p[key]
	return JSONValueToString(value)
}

func (p JSONObject) GetAsInt(key string) int {
	value, _ := p[key]
	return JSONValueToInt(value)
}

func (p JSONObject) GetAsInt64(key string) int64 {
	value, _ := p[key]
	return JSONValueToInt64(value)
}

func (p JSONObject) GetAsFloat64(key string) float64 {
	value, _ := p[key]
	return JSONValueToFloat64(value)
}

func (p JSONObject) GetAsBool(key string) bool {
	value, _ := p[key]
	return JSONValueToBool(value)
}

func (p JSONObject) GetAsObject(key string) JSONObject {
	value, _ := p[key]
	return JSONValueToObject(value)
}

func (p JSONObject) GetAsArray(key string) JSONArray {
	value, _ := p[key]
	return JSONValueToArray(value)
}

func (p JSONObject) GetAsTime(key string, format string) *time.Time {
	value, _ := p[key]
	return JSONValueToTime(value, format)
}
