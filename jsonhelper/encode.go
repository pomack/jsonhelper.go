// Copyright 2011 Aalok Shah. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// heavily borrowed from json/encode.go in the official Go source code

package jsonhelper

import (
    "encoding/base64"
    "json"
    "reflect"
    "runtime"
    "os"
    "sort"
    "strconv"
    "time"
    "unicode"
)

var byteSliceType = reflect.TypeOf([]byte(nil))
var timeType = reflect.TypeOf(time.Time{})

func Marshal(v interface{}) (retval interface{}, err os.Error) {
    defer func() {
        if r := recover(); r != nil {
            if _, ok := r.(runtime.Error); ok {
                panic(r)
            }
            err = r.(os.Error)
        }
    }()
    if v == nil {
        return nil, nil
    }
    e := &encodeState{}
    retval = e.reflectValue(reflect.ValueOf(v), false)
    return
}

func MarshalWithOptions(v interface{}, timeFormat string) (retval interface{}, err os.Error) {
    defer func() {
        if r := recover(); r != nil {
            if _, ok := r.(runtime.Error); ok {
                panic(r)
            }
            err = r.(os.Error)
        }
    }()
    if v == nil {
        return nil, nil
    }
    e := &encodeState{timeFormat:timeFormat}
    retval = e.reflectValue(reflect.ValueOf(v), false)
    return
}


// stringValues is a slice of reflect.Value holding *reflect.StringValue.
// It implements the methods to sort by string.
type stringValues []reflect.Value

func (sv stringValues) Len() int           { return len(sv) }
func (sv stringValues) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv stringValues) Less(i, j int) bool { return sv.get(i) < sv.get(j) }
func (sv stringValues) get(i int) string   { return sv[i].String() }

type encodeState struct {
    isObject bool
    isArray bool
    isFloat bool
    isInt bool
    isUint bool
    isString bool
    isBool bool
    isNull bool
    obj JSONObject
    arr JSONArray
    fValue float64
    iValue int64
    uValue uint64
    sValue string
    bValue bool
    timeFormat string
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c != '$' && c != '-' && c != '_' && !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func (e *encodeState) newWithSameOptions() *encodeState {
    o := new(encodeState)
    o.timeFormat = e.timeFormat
    return o
}

func (e *encodeState) Value() interface{} {
    if e.isObject { return e.obj }
    if e.isArray { return e.arr }
    if e.isFloat { return e.fValue }
    if e.isInt { return e.iValue }
    if e.isUint { return e.uValue }
    if e.isString { return e.sValue }
    if e.isBool { return e.bValue }
    if e.isNull { return nil }
    return nil
}

func (e *encodeState) error(err os.Error) {
    panic(err)
}

func (e *encodeState) reflectValue(v reflect.Value, stringify bool) (retval interface{}) {
    if !v.IsValid() {
        e.isNull = true
        return
    }
    if j, ok := v.Interface().(json.Marshaler); ok {
        b, err := j.MarshalJSON()
        if err == nil {
            var value interface{}
            err = json.Unmarshal(b, &value)
        }
        if err != nil {
            e.error(&json.MarshalerError{v.Type(), err})
        }
        return retval
    }
    switch v.Kind() {
    case reflect.Bool:
        if stringify {
            x := v.Bool()
            if x {
                e.sValue = "true"
            } else {
                e.sValue = "false"
            }
            e.isString = true
            retval = e.sValue
        } else {
            e.bValue = v.Bool()
            e.isBool = true
            retval = e.bValue
        }
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        if stringify {
            e.sValue = strconv.Itoa64(v.Int())
            e.isString = true
            retval = e.sValue
        } else {
            e.iValue = v.Int()
            e.isInt = true
            retval = e.iValue
        }
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
        if stringify {
            e.sValue = strconv.Uitoa64(v.Uint())
            e.isString = true
            retval = e.sValue
        } else {
            e.uValue = v.Uint()
            e.isUint = true
            retval = e.uValue
        }
    case reflect.Float32, reflect.Float64:
        if stringify {
            e.sValue = strconv.FtoaN(v.Float(), 'g', -1, v.Type().Bits())
            e.isString = true
            retval = e.sValue
        } else {
            e.fValue = v.Float()
            e.isFloat = true
            retval = e.fValue
        }
    case reflect.String:
        e.sValue = v.String()
        e.isString = true
        retval = e.sValue
    case reflect.Struct:
        t := v.Type()
        if t == timeType && e.timeFormat != "" {
            s := v.Interface().(time.Time)
            e.sValue = s.Format(e.timeFormat)
            e.isString = true
            retval = e.sValue
            break
        }
        n := v.NumField()
        e.obj = NewJSONObject()
        e.isObject = true
        for i := 0; i< n; i++ {
            f := t.Field(i)
            if f.PkgPath != "" {
                continue
            }
            tag, omitEmpty, collapse, stringify := f.Name, false, false, false
            if tv := f.Tag.Get("json"); tv != "" {
                name, opts := parseTag(tv)
                if isValidTag(name) {
                    tag = name
                }
                omitEmpty = opts.Contains("omitempty")
                stringify = opts.Contains("string")
                collapse = opts.Contains("collapse")
            }
            fieldValue := v.Field(i)
            if omitEmpty && isEmptyValue(fieldValue) {
                continue
            }
            subvalue := e.newWithSameOptions().reflectValue(fieldValue, stringify)
            if subvalue != nil {
                if subobj, ok := subvalue.(JSONObject); ok {
                    if omitEmpty && subobj.Len() == 0 {
                        continue
                    }
                    if collapse {
                        for k, v := range subobj {
                            e.obj[k] = v
                        }
                        continue
                    }
                }
            } else if omitEmpty {
                continue
            }
            e.obj[tag] = subvalue
        }
        retval = e.obj
    case reflect.Map:
        if v.Type().Key().Kind() != reflect.String {
            e.error(&json.UnsupportedTypeError{v.Type()})
        }
        if v.IsNil() {
            e.isNull = true
            break
        }
        e.isObject = true
        e.obj = NewJSONObject()
        var sv stringValues = v.MapKeys()
        sort.Sort(sv)
        for _, k := range sv {
            e.obj[k.String()] = e.newWithSameOptions().reflectValue(v.MapIndex(k), false)
        }
        retval = e.obj
    case reflect.Array, reflect.Slice:
        if v.Type() == byteSliceType {
            s := v.Interface().([]byte)
            dst := make([]byte, base64.StdEncoding.EncodedLen(len(s)))
            base64.StdEncoding.Encode(dst, s)
            e.isString = true
            e.sValue = string(dst)
            break
        }
        n := v.Len()
        arr := make([]interface{}, n)
        for i := 0; i < n; i++ {
            arr[i] = e.newWithSameOptions().reflectValue(v.Index(i), false)
        }
        e.arr = NewJSONArrayFromArray(arr)
        e.isArray = true
        retval = e.arr
    case reflect.Interface, reflect.Ptr:
        if v.IsNil() {
            e.isNull = true
            retval = nil
            return
        }
        retval = e.reflectValue(v.Elem(), false)
    default:
        e.error(&json.UnsupportedTypeError{v.Type()})
    }
    return
}

