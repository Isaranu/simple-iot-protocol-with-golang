package main

import (
  "log"
  "fmt"
  "net/http"
  "strconv"
  "time"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

func main(){

  http.HandleFunc("/", readychk)
  http.HandleFunc("/write/", writedata_iot)

  fmt.Printf("GO server port 9090")
  err := http.ListenAndServe(":9090", nil)
  if err != nil{
    log.Fatal(err)
    return
  }

}

func readychk(w http.ResponseWriter, r *http.Request)  {
  //fmt.Printf("Hi !")  //same as console log
  fmt.Fprintf(w, "ready to go")
}


/* - Write data from IoT device into MongoDB database - */
func writedata_iot(w http.ResponseWriter, r *http.Request)  {
  devid := r.URL.Query()["devid"]
  slot := r.URL.Query()["slot"]
  val := r.URL.Query()["val"]
  val_float, err := strconv.ParseFloat(string(val[0]), 64)

  t := int64(time.Now().Unix())
  now := time.Now()
  tm := now.Format(time.RFC3339)

  //fmt.Printf("devid " + string(devid[0]))
  //fmt.Printf("slot " + string(slot[0]))
  //fmt.Printf("val " + string(val[0]))

  session, err := mgo.Dial("127.0.0.1:27017")
  if err != nil{
    log.Fatal(err)
    return
  }
  defer session.Close()

  collection := session.DB("iotdatabygo").C(string(devid[0]))
  err = collection.Insert(bson.M{"devid":string(devid[0]),"slot":string(slot[0]),"val":val_float,"ts":t})
  if err != nil{
    log.Fatal(err)
    return
  }
  fmt.Fprintf(w, "record completed on " + tm)
}
