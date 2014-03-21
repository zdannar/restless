package response

import (
)

func New(d interface{}) *Response {
    rd := &ResponseData{}
    length := hasLength(d)
    if length >= 0 {
        rd.ItemsPerPage = length
        rd.TotalItems = length
    } else {
        rd.ItemsPerPage = 1
        rd.TotalItems = 1
    }
    rd.Items = d
    return &Response{
        ApiVersion: apiVersion,
        Data: rd,
        RestlessVersion: RESTLESS_RESP_VER}
}

type Response struct {
    ApiVersion      string          `json:"apiVersion"`
    RestlessVersion string          `json:"rstlsRespApiVer"`
    Data            *ResponseData   `json:"data"`
}

type ResponseData struct {
    TotalItems      int             `json:"totalItems,"`
    ItemsPerPage    int             `json:"itemsPerPage,omitempty"`
    Items           interface{}     `json:"items,omitempty"`
}

func NewRespStr(d interface{}) (string, error) {
    r := New(d)
    return marshal(r)
}




