package restless

var Log Logger
type Logger interface {
    Debugf(string, ...interface{})
    Errorf(string, ...interface{})
    Infof(string, ...interface{})
    Panicf(string, ...interface{})
    Warningf(string, ...interface{})
}

var (
    ErrInvalidMongoUrl = errors.New("restless.MongoUrl is not set and a valid MongoDb session was not provided to restless.Session")
    Session *mgo.Session = nil
    MongoUrl = "NULL"
    APIServAddr = "localhost:8080"
)
