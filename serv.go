package restless

import (
    "fmt"
    "github.com/gorilla/mux"
    "labix.org/v2/mgo"
    "net/http"
)

func ListenAndServe() {
    http.ListenAndServe(APIServAddr, nil)
}

func Register(dbname, collection string, constructor Constructor) error {

    session, err := getsess()
    if err != nil {
        return err
    }

    colRoot := fmt.Sprintf("/%s", collection)
    colIdRoot := fmt.Sprintf("%s/id", colRoot)

    r := mux.NewRouter()
    r.HandleFunc(colRoot, GetGenHandler(session, dbname, collection, constructor))
    r.HandleFunc(colIdRoot, GetIdHandler(session, dbname, collection, constructor))
    http.Handle(colRoot, r)
    return nil
}

func getsess() (*mgo.Session, error) {
    var err error
    if Session != nil {
        return Session, nil
    }

    if err = getmogurl(); err != nil {
        return nil, err
    }
    return mgo.Dial(MongoUrl)
}

func getmogurl() error {
    if MongoUrl == "NULL" {
        return ErrInvalidMongoUrl
    }
    return nil
}
