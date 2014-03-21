package restless/response

import (
    "reflect"
    "strings"
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

    if indent < 1 {
        jdata, err = json.Marshal(i)
    } else {
        jdata, err = json.MarshalIndent(i, "", strings.Repeat(" ", respIndent))
    }

    if err != nil {
        return "ERROR", err
    }
    return string(jdata), nil
}
