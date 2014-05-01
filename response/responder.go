package response

import ()

type Responder struct {}
func (r *Responder) Response(d interface{}) (string, error) {
    return NewRespStr(d)
}

func GetResponder() *Responder {
    return &Responder{}
}
