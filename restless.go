package restless

import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
    "net/http"
    "reflect"
)

func GetAll(c *mgo.Collection, ip interface{}) error {
    return c.Find(nil).All(ip)
}

func Insert(c *mgo.Collection, i interface{}) (string, error) {
    info, err := c.Upsert(bson.M{"_id": nil}, i)
    id := info.UpsertedId.(bson.ObjectId)
    return id.Hex(), err
}

func RemoveId(c *mgo.Collection, id bson.ObjectId) error {
    return c.RemoveId(id)
}

func GetId(c *mgo.Collection, i interface{}, id bson.ObjectId) error {
    return c.FindId(id).One(i)
}

func UpdateId(c *mgo.Collection, i interface{}, id bson.ObjectId) error {
    return c.UpdateId(id, i)
}

func GetGenHandler(s *mgo.Session, dbName string, colName string, n interface{}) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        var err error

        ns := s.Clone()
        defer ns.Close()

        col := ns.DB(dbName).C(colName)

        switch r.Method {
        //TODO: Add ability to queary specifics
        case "GET":
            i := reflect.New(reflect.SliceOf(reflect.TypeOf(n))).Interface()
            GetAll(col, i)

            r, err := resp.Response(i)
            if err != nil {
                http.Error(w, "Unable to parse form", http.StatusBadRequest)
            } else {
                w.Header().Add("Content-Type", "application/json")
                fmt.Fprintf(w, "%s", r)
            }

        case "POST":
            var lastId string

            i := reflect.New(reflect.TypeOf(n)).Interface()

            if err = r.ParseForm(); err != nil {
                http.Error(w, "Unable to parse form", http.StatusBadRequest)
                log.Errorf("Parsing form : %s", err)
            }

            jString := []byte(r.PostForm.Get(colName))
            if err = json.Unmarshal(jString, i); err != nil {
                http.Error(w, "Unable to unmarshal data", http.StatusBadRequest)
                log.Errorf("UnMarshal error : %s", err)
                return
            }

            if lastId, err = Insert(col, i); err != nil {
                http.Error(w, "Unable to unmarshal data", http.StatusInternalServerError)
                log.Errorf("Insert Error : %#v", err)
                return
            }

            w.Header().Add("Location", fmt.Sprintf("%s/%s", r.URL, lastId))
            w.WriteHeader(http.StatusCreated)
        }
        return
    }
}

func GetIdHandler(s *mgo.Session, dbName string, colName string, n interface{}) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        var err error
        var ids string

        ns := s.Clone()
        defer ns.Close()

        vars := mux.Vars(r)
        ids = vars["id"]
        log.Debugf("Provided ID = %s", ids)

        col := ns.DB(dbName).C(colName)

        if !bson.IsObjectIdHex(ids) {
            http.Error(w, "Provided ID is unknown", http.StatusNotFound)
            return
        }

        id := bson.ObjectIdHex(ids)
        i := reflect.New(reflect.TypeOf(n)).Interface()

        if err = GetId(col, i, id); err != nil {
            http.Error(w, "Provided ID is unknown", http.StatusNotFound)
            return
        }

        switch r.Method {
        case "GET":

            r, err := resp.Response(i)
            if err != nil {
                http.Error(w, "", http.StatusBadRequest)
            } else {
                w.Header().Add("Content-Type", "application/json")
                fmt.Fprintf(w, "%s", r)
            }

        case "PUT":
            if r.ParseForm(); err != nil {
                http.Error(w, "", http.StatusBadRequest)
            }

            if err = json.Unmarshal([]byte(r.PostForm.Get(colName)), i); err != nil {
                http.Error(w, "", http.StatusBadRequest)
                log.Errorf("UnMarshal error : %s", err)
            }

            if err = UpdateId(col, i, id); err != nil {
                http.Error(w, "Failed to update provided ID", http.StatusInternalServerError)
                log.Errorf("UnMarshal error : %s", err)
            }

        case "DELETE":
            if err = RemoveId(col, id); err != nil {
                http.Error(w, "Failed to remove provided ID", http.StatusInternalServerError)
                log.Errorf("Failed to remove id %s; error : %s", id, err)
            }
        }
        return
    }
}
