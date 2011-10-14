// Copyright 2011 Aalok Shah. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonhelper

import (
    "json"
    "strconv"
    "strings"
    "time"
)

func JSONValueToString(value interface{}) string {
    switch v := value.(type) {
    case nil:
        return ""
    case string:
        return v
    case int:
        return strconv.Itoa(v)
    case int64:
        return strconv.Itoa64(v)
    case float64:
        return strconv.Ftoa64(v, 'g', -1)
    case bool:
        if v {
            return "true"
        }
        return "false"
    }
    bytes, _ := json.Marshal(value)
    return string(bytes)
}

func JSONValueToInt(value interface{}) int {
    switch v := value.(type) {
    case nil:
        return 0
    case int:
        return v
    case float64:
        return int(v)
    case int64:
        return int(v)
    case string:
        i, _ := strconv.Atoi(v)
        return i
    case bool:
        if v {
            return 1
        }
        return 0
    case JSONObject:
        return len(v)
    case JSONArray:
        return len(v)
    case map[string]interface{}:
        return len(v)
    case []interface{}:
        return len(v)
    }
    return 0
}

func JSONValueToInt64(value interface{}) int64 {
    switch v := value.(type) {
    case nil:
        return 0
    case int:
        return int64(v)
    case float64:
        return int64(v)
    case int64:
        return v
    case string:
        i, _ := strconv.Atoi64(v)
        return i
    case bool:
        if v {
            return 1
        }
        return 0
    case JSONObject:
        return int64(len(v))
    case JSONArray:
        return int64(len(v))
    case map[string]interface{}:
        return int64(len(v))
    case []interface{}:
        return int64(len(v))
    }
    return 0
}

func JSONValueToFloat64(value interface{}) float64 {
    switch v := value.(type) {
    case nil:
        return 0
    case int:
        return float64(v)
    case float64:
        return v
    case int64:
        return float64(v)
    case string:
        i, _ := strconv.Atof64(v)
        return i
    case bool:
        if v {
            return 1
        }
        return 0
    case JSONObject:
        return float64(len(v))
    case JSONArray:
        return float64(len(v))
    case map[string]interface{}:
        return float64(len(v))
    case []interface{}:
        return float64(len(v))
    }
    return 0
}

func JSONValueToBool(value interface{}) bool {
    switch v := value.(type) {
    case nil:
        return false
    case bool:
        return v
    case int:
        return v != 0
    case float64:
        return v != 0.0
    case int64:
        return v != 0
    case string:
        s := strings.ToLower(v)
        return s == "true" || s == "1" || s == "yes"
    case JSONObject:
        return len(v) > 0
    case JSONArray:
        return len(v) > 0
    case map[string]interface{}:
        return len(v) > 0
    case []interface{}:
        return len(v) > 0
    }
    return false
}

func JSONValueToObject(value interface{}) JSONObject {
    switch v := value.(type) {
    case nil, bool, int, float64, int64, string, JSONArray, []interface{}:
        return NewJSONObject()
    case JSONObject:
        return v
    case map[string]interface{}:
        return NewJSONObjectFromMap(v)
    }
    return NewJSONObject()
}

func JSONValueToArray(value interface{}) JSONArray {
    switch v := value.(type) {
    case nil, bool, int, float64, int64, string, JSONObject, map[string]interface{}:
        return NewJSONArray()
    case JSONArray:
        return v
    case []interface{}:
        return NewJSONArrayFromArray(v)
    }
    return NewJSONArray()
}

func JSONValueToTime(value interface{}, format string) *time.Time {
    switch v := value.(type) {
    case nil, bool, JSONArray, JSONObject, []interface{}, map[string]interface{}:
        return nil
    case string:
        t, _ := time.Parse(format, v)
        return t
    case int64:
        return time.SecondsToUTC(v)
    case int:
        return time.SecondsToUTC(int64(v))
    case float64:
        return time.SecondsToUTC(int64(v))
    case time.Time:
        return &v
    case *time.Time:
        return v
    }
    return nil
}
