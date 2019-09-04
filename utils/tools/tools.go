package tools

import (
  "crypto/md5"
  crand "crypto/rand"
  "encoding/base64"
  "encoding/hex"
  "io"
  "math/rand"
  "time"
)
//generate random password
func  GetRandomString(l int) string {
  str := "0123456789QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm"
  bytes := []byte(str)
  var result []byte
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  for i := 0; i < l; i++ {
    result = append(result, bytes[r.Intn(len(bytes))])
  }
  return string(result)
}


func GetMd5String(s string) string {
  h := md5.New()
  h.Write([]byte(s))
  return hex.EncodeToString(h.Sum(nil))
}

//generate
func GetGuid() string {
  b := make([]byte, 48)

  if _, err := io.ReadFull(crand.Reader, b); err != nil {
    return ""
  }
  return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

func RemoveDuplicateSliceInt64(s []int64) []int64 {
  result := make([]int64, 0, len(s))
  temp := map[int64]struct{}{}
  for _, item := range s {
    if _, ok := temp[item]; !ok {
      temp[item] = struct{}{}
      result = append(result, item)
    }
  }
  return result
}

func sevenDayString(from string) string {
  tm, _ := time.Parse("2006-01-02", from)
  return tm.AddDate(0, 0, 6).Format("2006-01-02")
}
