package restless

import (
    "fmt"
    "github.com/gorilla/mux"
    "labix.org/v2/mgo"
    "net/http"
)

var router *mux.Router

func ListenAndServe() {
    Log.Infof("RESTLESS started on %s", APIServAddr)
    http.ListenAndServe(APIServAddr, nil)
}

func AddHandler(dbname, collection string, ytype interface{}) error {

    session, err := getsess()
    if err != nil {
        return err
    }

    colRoot := fmt.Sprintf("/%s", collection)
    colIdRoot := fmt.Sprintf("%s/{id}", colRoot)

    Log.Debugf("Adding general handler (%s)", colRoot)
    Log.Debugf("Adding ID based handler (%s)", colIdRoot)

    if router == nil {
        router = mux.NewRouter()
    }

    router.HandleFunc(colRoot, GetGenHandler(session, dbname, collection, ytype))
    router.HandleFunc(colIdRoot, GetIdHandler(session, dbname, collection, ytype))
    return nil
}

func Register() {
    rpath := "/"
    Log.Debugf("Setting root handler (%s)", rpath)
    http.Handle(rpath, router)
}


func getsess() (*mgo.Session, error) {
    var err error
    if Session != nil {
        return Session, nil
    }

    if err = getmogurl(); err != nil {
        return nil, err
    }
    Log.Infof("Establishing MongoDB connection for cloning(%s)", MongoUrl)
    return mgo.Dial(MongoUrl)
}

func getmogurl() error {
    if MongoUrl == "NULL" {
        return ErrInvalidMongoUrl
    }
    return nil
}
