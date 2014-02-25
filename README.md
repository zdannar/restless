restless
========

###SUMMARY
A simple implementation of building a REST API on top of MongoDB.    

More documentation and functionality to come...

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
[zdannar@pcbsd-4462 ~/go_include]$ curl -k -v -X POST -d 'square={"tlc":{"x":0,"y":1},"trc":{"x":1,"y":1},"blc":{"x":0,"y":0},"brc":{"x":1,"y":0}}' 'http://localhost:8080/square'
* Adding handle: conn: 0x8038e4600
* Adding handle: send: 0
* Adding handle: recv: 0
* Curl_addHandleToPipeline: length: 1
* - Conn 0 (0x8038e4600) send_pipe: 1, recv_pipe: 0
* About to connect() to localhost port 8080 (#0)
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> POST /square HTTP/1.1
> User-Agent: curl/7.33.0
> Host: localhost:8080
> Accept: */*
> Content-Length: 88
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 88 out of 88 bytes
< HTTP/1.1 201 Created
< Location: /square/530c2bb59d2fbb8476c89f79
< Date: Tue, 25 Feb 2014 05:35:49 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8
< 
* Connection #0 to host localhost left intact


[zdannar@pcbsd-4462 ~/go_include]$ curl -k -v -X POST -d 'square={"tlc":{"x":0,"y":9},"trc":{"x":9,"y":9},"blc":{"x":0,"y":0},"brc":{"x":9,"y":0}}' 'http://localhost:8080/square'
* Adding handle: conn: 0x8038e4600
* Adding handle: send: 0
* Adding handle: recv: 0
* Curl_addHandleToPipeline: length: 1
* - Conn 0 (0x8038e4600) send_pipe: 1, recv_pipe: 0
* About to connect() to localhost port 8080 (#0)
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> POST /square HTTP/1.1
> User-Agent: curl/7.33.0
> Host: localhost:8080
> Accept: */*
> Content-Length: 88
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 88 out of 88 bytes
< HTTP/1.1 201 Created
< Location: /square/530c2e839d2fbb8476c89f7a
< Date: Tue, 25 Feb 2014 05:47:47 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8
< 
* Connection #0 to host localhost left intact


[zdannar@pcbsd-4462 ~/go_include]$ curl -k -v -X GET 'http://localhost:8080/square'
* Adding handle: conn: 0x8038e4600
* Adding handle: send: 0
* Adding handle: recv: 0
* Curl_addHandleToPipeline: length: 1
* - Conn 0 (0x8038e4600) send_pipe: 1, recv_pipe: 0
* About to connect() to localhost port 8080 (#0)
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> GET /square HTTP/1.1
> User-Agent: curl/7.33.0
> Host: localhost:8080
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Tue, 25 Feb 2014 05:48:35 GMT
< Content-Length: 240
< 
* Connection #0 to host localhost left intact
[]square=[{"_id":"530c2bb59d2fbb8476c89f79","tlc":{"x":0,"y":1},"trc":{"x":1,"y":1},"blc":{"x":0,"y":0},"brc":{"x":1,"y":0}},{"_id":"530c2e839d2fbb8476c89f7a","tlc":{"x":0,"y":9},"trc":{"x":9,"y":9},"blc":{"x":0,"y":0},"brc":{"x":9,"y":0}}]


[zdannar@pcbsd-4462 ~/go_include]$ curl -k -v -X GET 'http://localhost:8080/square/530c2e839d2fbb8476c89f7a'
* Adding handle: conn: 0x8038e4600
* Adding handle: send: 0
* Adding handle: recv: 0
* Curl_addHandleToPipeline: length: 1
* - Conn 0 (0x8038e4600) send_pipe: 1, recv_pipe: 0
* About to connect() to localhost port 8080 (#0)
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> GET /square/530c2e839d2fbb8476c89f7a HTTP/1.1
> User-Agent: curl/7.33.0
> Host: localhost:8080
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Tue, 25 Feb 2014 05:48:48 GMT
< Content-Length: 121
< 
* Connection #0 to host localhost left intact
square={"_id":"530c2e839d2fbb8476c89f7a","tlc":{"x":0,"y":9},"trc":{"x":9,"y":9},"blc":{"x":0,"y":0},"brc":{"x":9,"y":0}}[zdannar@pcbsd-4462 ~/go_include]$ 


[zdannar@pcbsd-4462 ~/go_include]$ curl -k -v -X PUT -d 'square={"tlc":{"x":0,"y":2},"trc":{"x":2,"y":2},"blc":{"x":0,"y":0},"brc":{"x":2,"y":0}}' 'http://localhost:8080/square/530c2bb59d2fbb8476c89f79'
* Adding handle: conn: 0x8038e4600
* Adding handle: send: 0
* Adding handle: recv: 0
* Curl_addHandleToPipeline: length: 1
* - Conn 0 (0x8038e4600) send_pipe: 1, recv_pipe: 0
* About to connect() to localhost port 8080 (#0)
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> PUT /square/530c2bb59d2fbb8476c89f79 HTTP/1.1
> User-Agent: curl/7.33.0
> Host: localhost:8080
> Accept: */*
> Content-Length: 88
> Content-Type: application/x-www-form-urlencoded
>
* upload completely sent off: 88 out of 88 bytes
< HTTP/1.1 200 OK
< Date: Tue, 25 Feb 2014 05:50:23 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host localhost left intact



[zdannar@pcbsd-4462 ~/go_include]$ curl -k -v -X GET 'http://localhost:8080/square'
* Adding handle: conn: 0x8038e4600
* Adding handle: send: 0
* Adding handle: recv: 0
* Curl_addHandleToPipeline: length: 1
* - Conn 0 (0x8038e4600) send_pipe: 1, recv_pipe: 0
* About to connect() to localhost port 8080 (#0)
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> GET /square HTTP/1.1
> User-Agent: curl/7.33.0
> Host: localhost:8080
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Tue, 25 Feb 2014 05:50:38 GMT
< Content-Length: 240
<
* Connection #0 to host localhost left intact
[]square=[{"_id":"530c2bb59d2fbb8476c89f79","tlc":{"x":0,"y":2},"trc":{"x":2,"y":2},"blc":{"x":0,"y":0},"brc":{"x":2,"y":0}},{"_id":"530c2e839d2fbb8476c89f7a","tlc":{"x":0,"y":9},"trc":{"x":9,"y":9},"blc":{"x":0,"y":0},"brc":{"x":9,"y":0}}]


[zdannar@pcbsd-4462 ~/go_include/src/restless]$ curl -k -v -X DELETE 'http://localhost:8080/square/530c2bb59d2fbb8476c89f79'
* Adding handle: conn: 0x8038e4600
* Adding handle: send: 0
* Adding handle: recv: 0
* Curl_addHandleToPipeline: length: 1
* - Conn 0 (0x8038e4600) send_pipe: 1, recv_pipe: 0
* About to connect() to localhost port 8080 (#0)
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> DELETE /square/530c2bb59d2fbb8476c89f79 HTTP/1.1
> User-Agent: curl/7.33.0
> Host: localhost:8080
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Tue, 25 Feb 2014 05:54:48 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host localhost left intact


[zdannar@pcbsd-4462 ~/go_include/src/restless]$ curl -k -v -X GET 'http://localhost:8080/square'
* Adding handle: conn: 0x8038e4600
* Adding handle: send: 0
* Adding handle: recv: 0
* Curl_addHandleToPipeline: length: 1
* - Conn 0 (0x8038e4600) send_pipe: 1, recv_pipe: 0
* About to connect() to localhost port 8080 (#0)
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> GET /square HTTP/1.1
> User-Agent: curl/7.33.0
> Host: localhost:8080
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Tue, 25 Feb 2014 05:55:06 GMT
< Content-Length: 125
<
* Connection #0 to host localhost left intact
[]square=[{"_id":"530c2e839d2fbb8476c89f7a","tlc":{"x":0,"y":9},"trc":{"x":9,"y":9},"blc":{"x":0,"y":0},"brc":{"x":9,"y":0}}]
```
