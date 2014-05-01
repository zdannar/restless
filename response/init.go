package response

import (
)

const (
    RESTLESS_RESP_VER = "0.3"
    INDENT_PRETTY = 4
)

var (
    apiVersion string = "UNKNOWN"
    respIndent int    = INDENT_PRETTY
)

func SetApiVersion(vstr string) {
    apiVersion = vstr
}

func SetRespIndent(i int) {
    respIndent = i
}

type Resp struct{}
func(r *Resp) Response(i interface{}) (string, error) {
    return NewRespStr(i)
}

