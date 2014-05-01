restless
========

###SUMMARY
A simple implementation of building a REST API on top of MongoDB.    

More documentation and functionality to come...

Current TODOS:
- Align more closely with google JSON style guide.
- Contemplate abstracting BSON.Id from being required in structure.
- Add query format to restless.
- Cleanup code, specifically GET/POST functionality into cleaner functions.

###INSTALL
go get github.com/zdannar/restless

###EXAMPLE
```go
package main 

import (
    "labix.org/v2/mgo/bson"
    "github.com/zdannar/restless"
    "github.com/zdannar/flogger"
    "encoding/json"

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

func main() {

    var err error
    var log *flogger.Flogger = flogger.New(flogger.DEBUG, flogger.FLOG_FORMAT, flogger.FLOG_LEVELS)
    defer log.Close()

    s := Square{Tlc: Point{X:0, Y:1}, Trc: Point{X:1, Y:1}, Blc: Point{X:0, Y:0}, Brc: Point{X:1, Y:0}}
    jdata, err := json.Marshal(&s)
    log.Infof("JDATA : %s", jdata)

    restless.Log = log
    restless.MongoUrl = "mongodb://localhost"

    err = restless.AddHandler("restless_example", "square", Square{})
    if err != nil {
        log.Fatalf("Somthing went wrong : %s", err)
    }

    restless.Register()
    restless.ListenAndServe()
}
```

Curl based example.

```
[zdannar@bsd-org ~/go_include/src/restless]$ curl -k -v -X POST -d 'square={"tlc":{"x":0,"y":1},"trc":{"x":1,"y":1},"blc":{"x":0,"y":0},"brc":{"x":1,"y":0}}' 'http://localhost:8080/square'
< HTTP/1.1 201 Created
< Location: /square/5361b64a029f8b60d93d1ee5
< Date: Thu, 01 May 2014 02:49:46 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8


[zdannar@bsd-org ~/go_include/src/restless]$ curl -k -v -X POST -d 'square={"tlc":{"x":0,"y":9},"trc":{"x":9,"y":9},"blc":{"x":0,"y":0},"brc":{"x":9,"y":0}}' 'http://localhost:8080/square'
< HTTP/1.1 201 Created
< Location: /square/5361b657029f8b60d93d1ee6
< Date: Thu, 01 May 2014 02:49:59 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8
< 


[zdannar@bsd-org ~/go_include/src/restless]$ curl -k -v -X GET 'http://localhost:8080/square'
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Thu, 01 May 2014 02:50:33 GMT
< Transfer-Encoding: chunked
< 
{
    "apiVersion": ".7",
    "rstlsRespApiVer": "0.3",
    "data": {
        "totalItems": 1,
        "itemsPerPage": 1,
        "items": [
            {
                "_id": "535aa21919057eb16cad555b",
                "tlc": {
                    "x": 0,
                    "y": 9
                },
                "trc": {
                    "x": 9,
                    "y": 9
                },
                "blc": {
                    "x": 0,
                    "y": 0
                },
                "brc": {
                    "x": 9,
                    "y": 0
                }
            },
            {
                "_id": "5361b64a029f8b60d93d1ee5",
                "tlc": {
                    "x": 0,
                    "y": 1
                },
                "trc": {
                    "x": 1,
                    "y": 1
                },
                "blc": {
                    "x": 0,
                    "y": 0
                },
                "brc": {
                    "x": 1,
                    "y": 0
                }
            },
            {
                "_id": "5361b657029f8b60d93d1ee6",
                "tlc": {
                    "x": 0,
                    "y": 9
                },
                "trc": {
                    "x": 9,
                    "y": 9
                },
                "blc": {
                    "x": 0,
                    "y": 0
                },
                "brc": {
                    "x": 9,
                    "y": 0
                }
            }
        ]
    }

[zdannar@bsd-org ~/go_include/src/restless]$ curl -k -v -X GET 'http://localhost:8080/square/5361b657029f8b60d93d1ee6'
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Thu, 01 May 2014 02:51:31 GMT
< Content-Length: 536
< 
{
    "apiVersion": ".7",
    "rstlsRespApiVer": "0.3",
    "data": {
        "totalItems": 1,
        "itemsPerPage": 1,
        "items": {
            "_id": "5361b657029f8b60d93d1ee6",
            "tlc": {
                "x": 0,
                "y": 9
            },
            "trc": {
                "x": 9,
                "y": 9
            },
            "blc": {
                "x": 0,
                "y": 0
            },
            "brc": {
                "x": 9,
                "y": 0
            }
        }
    }

[zdannar@bsd-org ~/go_include/src/restless]$ curl -k -v -X PUT -d 'square={"tlc":{"x":0,"y":2},"trc":{"x":2,"y":2},"blc":{"x":0,"y":0},"brc":{"x":2,"y":0}}' 'http://localhost:8080/square/5361b657029f8b60d93d1ee6'
< HTTP/1.1 200 OK
< Date: Thu, 01 May 2014 02:53:06 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8


[zdannar@bsd-org ~/go_include/src/restless]$ curl -k -v -X GET 'http://localhost:8080/square/5361b657029f8b60d93d1ee6'
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Thu, 01 May 2014 02:53:28 GMT
< Content-Length: 536
< 
{
    "apiVersion": ".7",
    "rstlsRespApiVer": "0.3",
    "data": {
        "totalItems": 1,
        "itemsPerPage": 1,
        "items": {
            "_id": "5361b657029f8b60d93d1ee6",
            "tlc": {
                "x": 0,
                "y": 2
            },
            "trc": {
                "x": 2,
                "y": 2
            },
            "blc": {
                "x": 0,
                "y": 0
            },
            "brc": {
                "x": 2,
                "y": 0
            }
        }
    }
}

[zdannar@bsd-org ~/go_include/src/restless]$ curl -k -v -X DELETE 'http://localhost:8080/square/5361b657029f8b60d93d1ee6'
< HTTP/1.1 200 OK
< Date: Thu, 01 May 2014 02:54:10 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8
```
