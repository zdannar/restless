package restless

import (
    "errors"
    "labix.org/v2/mgo"
)

var Log Logger

type Logger interface {
    Debugf(string, ...interface{})
    Errorf(string, ...interface{})
    Infof(string, ...interface{})
    Panicf(string, ...interface{})
    Warningf(string, ...interface{})
}

var (
    Session     *mgo.Session = nil
    MongoUrl                 = "NULL"
    APIServAddr              = ":8080"
)

var ErrInvalidMongoUrl = errors.New("restless.MongoUrl is not set and a valid MongoDb session was not provided to restless.Session")
