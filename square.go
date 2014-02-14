package restless

import (
    "labix.org/v2/mgo/bson"
)

type Point struct {
    X   int `json:"x" bson:"x"`
    Y   int `json:"y" bson:"y"`
}

type Square struct {
    Id  bson.ObjectId `json:"_id,omitempty"    bson:"_id,omitempty"`
    Tlc Point         `json:"tlc"              bson:"tlc"`
    Trc Point         `json:"trc"              bson:"trc"`
    Blc Point         `json:"blc"              bson:"blc"`
    Brc Point         `json:"brc"              bson:"brc"`
}

type SquareConst struct{}

func (_ *Square) Single() interface{} {
    return &Square{}
}

func (_ *SquareConst) Slice() interface{} {
    items := make([]Square, 0)
    return &items
}
