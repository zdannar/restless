package restless

import (
    "errors"
    "labix.org/v2/mgo"
    "restless/response"
)

var (
    log         Logger       = nil
    resp        Responser    = response.GetResponder()
    Session     *mgo.Session = nil
    MongoUrl                 = "NULL"
    APIServAddr              = ":8080"
)
var ErrInvalidMongoUrl = errors.New("restless.MongoUrl is not set and a valid MongoDb session was not provided to restless.Session")

func init() {
    response.SetApiVersion(".7")
}

type Logger interface {
    Debugf(string, ...interface{})
    Errorf(string, ...interface{})
    Infof(string, ...interface{})
    Panicf(string, ...interface{})
    Warningf(string, ...interface{})
}

func SetLogger(l Logger) {
    log = l
}


type Responser interface {
    Response(interface{}) (string, error)
}

func SetResponser(r Responser) {
    resp = r
}
