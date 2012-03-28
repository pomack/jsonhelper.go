// Copyright 2011 Aalok Shah. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonhelper

import (
    "container/list"
    "encoding/json"
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

func (p JSONArray) String() string {
    b, _ := json.Marshal(&p)
    return string(b)
}

func (p JSONArray) GetAsString(index int) string {
    value := p[index]
    return JSONValueToString(value)
}

func (p JSONArray) GetAsInt(index int) int {
    value := p[index]
    return JSONValueToInt(value)
}

func (p JSONArray) GetAsInt32(index int) int32 {
    value := p[index]
    return JSONValueToInt32(value)
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

func (p JSONArray) GetAsTime(index int, format string) time.Time {
    value := p[index]
    return JSONValueToTime(value, format)
}

func (p JSONArray) Compact(removeFalse bool, removeEmptyStrings bool, removeZero bool, removeEmptyArrays bool, removeEmptyObjects bool) JSONArray {
    if len(p) == 0 {
        if removeEmptyArrays {
            return nil
        }
        return p
    }
    l := list.New()
    for _, v := range p {
        var value interface{}
        value = v
        switch t := v.(type) {
        case nil:
            continue
        case string:
            if removeEmptyStrings && len(t) == 0 {
                continue
            }
        case JSONObject:
            value = t.Compact(removeFalse, removeEmptyStrings, removeZero, removeEmptyArrays, removeEmptyObjects)
        case JSONArray:
            value = t.Compact(removeFalse, removeEmptyStrings, removeZero, removeEmptyArrays, removeEmptyObjects)
        case map[string]interface{}:
            value = NewJSONObjectFromMap(t).Compact(removeFalse, removeEmptyStrings, removeZero, removeEmptyArrays, removeEmptyObjects)
        case []interface{}:
            value = NewJSONArrayFromArray(t).Compact(removeFalse, removeEmptyStrings, removeZero, removeEmptyArrays, removeEmptyObjects)
        case float64:
            if removeZero && t == 0.0 {
                continue
            }
        case float32:
            if removeZero && t == 0.0 {
                continue
            }
        case int64:
            if removeZero && t == 0 {
                continue
            }
        case int32:
            if removeZero && t == 0 {
                continue
            }
        case int:
            if removeZero && t == 0 {
                continue
            }
        case int16:
            if removeZero && t == 0 {
                continue
            }
        case int8:
            if removeZero && t == 0 {
                continue
            }
        case byte:
            if removeZero && t == 0 {
                continue
            }
        case bool:
            if removeFalse && t == false {
                continue
            }
        }
        if value == nil {
            continue
        }
        l.PushBack(value)
    }
    lLen := l.Len()
    if removeEmptyArrays && lLen == 0 {
        return nil
    }
    arr := make([]interface{}, lLen)
    for e, i := l.Front(), 0; e != nil; e, i = e.Next(), i+1 {
        arr[i] = e.Value
    }
    return NewJSONArrayFromArray(arr)
}
