package jsonhelper

import (
    "time"
)

type JSONArray []interface{}

func NewJSONArray() JSONArray {
    return make([]interface{}, 0)
}

func NewJSONArrayFromArray(value []interface{}) JSONArray {
    return JSONArray(value)
}

func (p JSONArray) Len() int {
    return len(p)
}

func (p JSONArray) GetAsString(index int) string {
    value := p[index]
    return JSONValueToString(value)
}

func (p JSONArray) GetAsInt(index int) int {
    value := p[index]
    return JSONValueToInt(value)
}

func (p JSONArray) GetAsInt64(index int) int64 {
    value := p[index]
    return JSONValueToInt64(value)
}

func (p JSONArray) GetAsFloat64(index int) float64 {
    value := p[index]
    return JSONValueToFloat64(value)
}

func (p JSONArray) GetAsObject(index int) JSONObject {
    value := p[index]
    return JSONValueToObject(value)
}

func (p JSONArray) GetAsArray(index int) JSONArray {
    value := p[index]
    return JSONValueToArray(value)
}

func (p JSONArray) GetAsTime(index int, format string) *time.Time {
    value := p[index]
    return JSONValueToTime(value, format)
}
