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

func GetShapeHandler(s *mgo.Session, dbName string, colName string, cns Constructor) http.HandlerFunc {
   return func(w http.ResponseWriter, r *http.Request) {
       var jdata []byte

       var err error
       ns := s.Clone()
       defer ns.Close()

       col := ns.DB(dbName).C(colName)

       if r.Method == "GET" {
           i := cns.Slice()
           //TODO: Add ability to queary specifics
           GetAll(col, i)
           jdata, _ = json.Marshal(i)
       }

       if r.Method == "POST" {
           i := cns.Single()
           err = r.ParseForm()
           if err != nil {
               http.Error(w, "Unable to parse form", http.StatusNoContent)
           }

           jString := []byte(r.PostForm.Get("json"))
           err = json.Unmarshal(jString, i)
           if err != nil {
               log.Panicf("UnMarshal error : %s", err)
           }

           lastId, err := Insert(col, i)
           if err != nil  {
               log.Panicf("Insert Error : %#v", err)
           }
           jdata, _ = json.Marshal(i)

           w.Header().Add("Location", fmt.Sprintf("%s/%s", r.URL, lastId))
           w.WriteHeader(http.StatusCreated)

       }


       w.Header().Add("Content-Type", "application/json")

       fmt.Fprintf(w,"%s", jdata)
    }
}

func GetShapeByIdHandler(s *mgo.Session, dbName string, colName string, cns Constructor) http.HandlerFunc {
   return func (w http.ResponseWriter, r *http.Request) {

       var jdata []byte
       var err error

       ns := s.Clone()
       defer ns.Close()
       vars := mux.Vars(r)
       ids := vars["id"]

       col := ns.DB(dbName).C(colName)

       if !bson.IsObjectIdHex(ids) {
            log.Panicf("ObjectId(%s) is not hex", ids)
       }

       id := bson.ObjectIdHex(ids)

       i := cns.Single()
       //Get Object regardless, Reconsider later
       err = GetId(col, i, id)
       if err != nil {
            log.Panic("Unknown object")
       }

       //TODO: Another error to catch
       jdata, _ = json.Marshal(i)

       if r.Method == "PUT" {
           err := r.ParseForm()
           if err != nil {
               http.Error(w, "Unable to parse form", http.StatusNoContent)
           }
           //TODO: Need error handling
           ijson :=  r.PostForm.Get("json")
           err = json.Unmarshal([]byte(ijson), i)
           if err != nil {
               log.Panicf("UnMarshal error : %s", err)
           }

           //TODO : Need to handle false insert and hand back method
           UpdateId(col, i, id)
       }

       if r.Method == "DELETE" {
            //TODO: Catch Error
            err = RemoveId(col, id) 
            if err != nil {
               log.Panicf("Failed to remove id %s; error : %s", id, err)
            }
       }
       w.Header().Add("Content-Type", "application/json")
       fmt.Fprintf(w, "%s", jdata)
       return
    }
}
