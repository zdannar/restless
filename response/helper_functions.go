package response

import (
    "reflect"
    "strings"
    "encoding/json"
)

func isSlice(i interface{}) bool {
    return reflect.TypeOf(i).Kind() == reflect.Slice
}

func hasLength(i interface{}) int {
    if !isSlice(i) {
        return -1
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
