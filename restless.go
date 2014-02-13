package restless

import (
   "github.com/gorilla/mux"
   "labix.org/v2/mgo"
   "labix.org/v2/mgo/bson"
   "net/http"
   "fmt"
   "encoding/json"
   log "github.com/zdannar/flogger"
)

func GetAll(c *mgo.Collection, ip interface{}) {
   c.Find(nil).All(ip)
}

func Insert(c *mgo.Collection, i interface{}) (string, error) {
    info, err := c.Upsert(bson.M{"_id" : nil}, i)
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

func GetGenHandler(s *mgo.Session, dbName string, colName string, cns Constructor) http.HandlerFunc {
   return func(w http.ResponseWriter, r *http.Request) {

        var jdata []byte
        var err error

        ns := s.Clone()
        defer ns.Close()

        col := ns.DB(dbName).C(colName)
        
        switch r.Method {
            //TODO: Add ability to queary specifics
            case "GET":
                i := cns.Slice()
                GetAll(col, i)
                jdata, err = json.Marshal(i)

                w.Header().Add("Content-Type", "application/json")
                fmt.Fprintf(w, "%s", jdata)

            case "PUT":
                var lastId string 

                i := cns.Single()
                if err = r.ParseForm(); err != nil {
                    http.Error(w, "Unable to parse form", http.StatusBadRequest)
                    log.Panicf("%s", err)
                }

                jString := []byte(r.PostForm.Get("json"))
                if err = json.Unmarshal(jString, i); err != nil {
                    http.Error(w, "Unable to unmarshal data", http.StatusBadRequest)
                    log.Panicf("UnMarshal error : %s", err)
                }

                if lastId, err = Insert(col, i); err != nil {
                    log.Panicf("Insert Error : %#v", err)
                }

                if jdata, err = json.Marshal(i); err != nil {
                    http.Error(w, "Marshal error", http.StatusInternalServerError)
                }
                w.Header().Add("Location", fmt.Sprintf("%s/%s", r.URL, lastId))
                w.WriteHeader(http.StatusCreated)
        }
        return
    }
}

func GetIdHandler(s *mgo.Session, dbName string, colName string, cns Constructor) http.HandlerFunc {
   return func (w http.ResponseWriter, r *http.Request) {

        var jdata []byte
        var err error
        var ids string

        ns := s.Clone()
        defer ns.Close()

        vars := mux.Vars(r)
        ids = vars["id"]

        col := ns.DB(dbName).C(colName)

        if !bson.IsObjectIdHex(ids) {
            http.Error(w, "Provided ID is unknown", http.StatusNotFound)
            return
        }

        id := bson.ObjectIdHex(ids)
        i := cns.Single()

        if err = GetId(col, i, id); err != nil {
            http.Error(w, "Provided ID is unknown", http.StatusNotFound)
            return 
        }

        if jdata, err = json.Marshal(i); err != nil {
            http.Error(w, "", http.StatusBadRequest)
        }


        switch r.Method {
            case "GET":
                w.Header().Add("Content-Type", "application/json")
                fmt.Fprintf(w, "%s", jdata)

            case "PUT":
                if r.ParseForm(); err != nil {
                    http.Error(w, "", http.StatusBadRequest)
                }

                if err = json.Unmarshal([]byte(r.PostForm.Get("json")), i); err != nil {
                    http.Error(w, "", http.StatusBadRequest)
                    log.Panicf("UnMarshal error : %s", err)
                }

                if err = UpdateId(col, i, id); err != nil {
                    http.Error(w, "Failed to update provided ID", http.StatusInternalServerError)
                    log.Panicf("UnMarshal error : %s", err)
                }

            case "DELETE":
                if err = RemoveId(col, id); err != nil { 
                    http.Error(w, "Failed to remove provided ID", http.StatusInternalServerError)
                    log.Panicf("Failed to remove id %s; error : %s", id, err)
                }
        }
        return 
    }
}
