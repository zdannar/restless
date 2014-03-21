package restless/response

import (
)


const (
    RESTLESS_RESP_VER "0.0"
    INDENT_PRETTY = 4
)

var (
    apiVersion string = "UNKNOWN"
    respIndent int    = 0
)

func SetApiVersion(vstr string) {
    apiVersion = vstr
}

func SetRespIndent(i int) {
    respIndent = i
}
