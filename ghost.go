package main

import (
   "github.com/gorilla/mux"
   "labix.org/v2/mgo"
   "labix.org/v2/mgo/bson"
   "net/http"
   "fmt"
   "encoding/json"
   log "github.com/zdannar/flogger"
)


type Point struct {
   X int `json:"x" bson:"x"`
   Y int `json:"y" bson:"y"`
}

type Square struct {
   Id bson.ObjectId    `json:"_id,omitempty"    bson:"_id,omitempty"`
   Tlc Point           `json:"tlc"              bson:"tlc"`
   Trc Point           `json:"trc"              bson:"trc"`
   Blc Point           `json:"blc"              bson:"blc"`
   Brc Point           `json:"brc"              bson:"brc"`
}

type Squares []Square


func GetAll(c *mgo.Collection, ip interface{}) {
   c.Find(nil).All(ip)
}


func Insert(c *mgo.Collection, i interface{}) (string, error) {
    info, err := c.Upsert(bson.M{"_id" : nil}, i)
    log.Debugf("Upserted changeInfo : %#v", info)
    id := info.UpsertedId.(bson.ObjectId)
    return id.Hex(), err

    //id := bson.NewObjectId()
    //info, err := c.UpsertId(id, i)
    //log.Debugf("Upserted changeInfo : %#v", info)
    //return info.UpsertedId.(string), err



/*
    info, err := c.UpsertId(id, i)
    if err != nil {
        return nil, err
    }
    log.Debugf("Upserted changeInfo : %#v", info)
    return &id, nil
*/
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
           log.Info("New iD is : %#v", lastId)

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
       fmt.Fprintf(w,"%s", jdata)
       return 
    }
}


/***********************************************************/
type Constructor interface {
    Single() interface{}
    Slice()  interface{}
}


type SquareConst struct {}
func (s *SquareConst) Single() interface{} {
    return &Square{}
}
func(s *SquareConst) Slice() interface{} {
    items:= make([]Square, 0)
    return &items
}


func main() {
   //var squares = make([]Square, 0)

   session, err := mgo.Dial("mongodb://localhost")
   if err != nil {
       panic(err)
   }

   h := GetShapeHandler(session, "play", "square", &SquareConst{})
   i := GetShapeByIdHandler(session, "play", "square", &SquareConst{})

   r := mux.NewRouter()
   r.HandleFunc("/square", h)
   r.HandleFunc("/square/{id}", i)

   http.Handle("/", r)
   http.ListenAndServe(":3000", nil)

}
