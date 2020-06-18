package main

import (
  "bytes"
  "encoding/json"
  "io/ioutil"
  "log"
  "net/http"
)

type RedisLastPush struct {
  Message string `json:"message"`
}

func main() {

  resp, _ := http.Get("http://localhost:8600/ping/test")
  result := RedisLastPush{}
  data, err := ioutil.ReadAll(resp.Body)
  log.Println(string(data))
  if err == nil && data != nil {
  err = json.Unmarshal(data, &result)
  if err != nil {
    log.Printf("fail to Unmarshal:%v", err)
  }
  }
  log.Printf("existValue:%+v", result)

  //data1, err := ioutil.ReadAll(resp.Body)
  //log.Println(string(data1))
  //if err == nil && data1 != nil {
  //  err = json.Unmarshal(data1, &result1)
  //  if err != nil {
  //    log.Printf("fail to Unmarshal:%v", err)
  //  }
  //}
  resp.Body.Close()
  resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
  result1 := RedisLastPush{}
  err = json.NewDecoder(resp.Body).Decode(&result1)
  if err != nil {
    log.Printf("fail to Unmarshal:%v", err)
  }
  log.Printf("existValue1:%+v", result1)


}
