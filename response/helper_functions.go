package response

import (
    "reflect"
    "strings"
    "encoding/json"
    //log "github.com/zdannar/flogger"
)

func isSlice(i interface{}) (bool, bool) {
    if reflect.TypeOf(i).Kind() == reflect.Ptr {
        // Get Indirect and check for slice
        return reflect.Indirect(reflect.ValueOf(i)).Kind() == reflect.Slice, true
    }

    return reflect.TypeOf(i).Kind() == reflect.Slice, false
}

func hasLength(i interface{}) int {
    is, indirectly := isSlice(i)
    if !is {
        return -1
    }
    if indirectly { 
        return reflect.Indirect(reflect.ValueOf(i)).Len()
    }
    return reflect.ValueOf(i).Len()
}

func marshal(i interface{}) (string, error) {
    var jdata []byte
    var err error

    if respIndent < 1 {
        jdata, err = json.Marshal(i)
    } else {
        jdata, err = json.MarshalIndent(i, "", strings.Repeat(" ", respIndent))
    }

    if err != nil {
        return "ERROR", err
    }
    return string(jdata), nil
}
