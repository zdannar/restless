package restless

import (
   "github.com/gorilla/mux"
   "labix.org/v2/mgo"
   "net/http"
   "fmt"
   "errors"
)

func ListenAndServe(netaddr string) {
    http.ListenAndServe(APIServAddr, nil)
}

func Register(dbname, collection string, construcutor Constructor) error {

    session, err := getsess(Session)
    if err != nil { 
        return err
    }

    colRoot := fmt.Sprintf("/%s", collection)
    colIdRoot := fmt.Sprintf("%s/id", colroot)

    r := mux.NewRouter()
    r.HandleFunc( colRoot, restless.GetGenHandler(Session, dbname, collection, constructor) )
    r.HandleFunc( colIdRoot, restless.GetIdHandler(Session, dbname, collection, constructor) )
    http.Handle(colRoot, r)

}

func getsess() (*mgo.Session, error) {
    var err error
    if !Session == nil {
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
    return MongoUrl, nil
}
