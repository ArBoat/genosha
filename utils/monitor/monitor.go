package monitor

import (
  "encoding/json"
  "expvar"
  "fmt"
  "github.com/gin-gonic/gin"
  "net/http"
  "runtime"
  "time"
)

//var CuMemoryPtr *map[string]models.Kline
//var BTCMemoryPtr *map[string]models.Kline


//
var start = time.Now()

// calculateUptime
func calculateUptime() interface{} {
  return time.Since(start).String()
}

// currentGoVersion
func currentGoVersion() interface{} {
  return runtime.Version()
}

// getNumCPUs
func getNumCPUs() interface{} {
  return runtime.NumCPU()
}

// getGoOS
func getGoOS() interface{} {
  return runtime.GOOS
}

// getNumGoroutins
func getNumGoroutins() interface{} {
  return runtime.NumGoroutine()
}

// getNumCgoCall
func getNumCgoCall() interface{} {
  return runtime.NumCgoCall()
}


var lastPause uint32

// getLastGCPauseTime
func getLastGCPauseTime() interface{} {
  var gcPause uint64
  ms := new(runtime.MemStats)

  statString := expvar.Get("memstats").String()
  if statString != "" {
    json.Unmarshal([]byte(statString), ms)

    if lastPause == 0 || lastPause != ms.NumGC {
      gcPause = ms.PauseNs[(ms.NumGC+255)%256]
      lastPause = ms.NumGC
    }
  }

  return gcPause
}

// GetCurrentRunningStats
func GetCurrentRunningStats(c *gin.Context) {
  c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
  first := true
  report := func(key string, value interface{}) {
    if !first {
      fmt.Fprintf(c.Writer, ",\n")
    }
    first = false
    if str, ok := value.(string); ok {
      fmt.Fprintf(c.Writer, "%q: %q", key, str)
    } else {
      fmt.Fprintf(c.Writer, "%q: %v", key, value)
    }
  }

  fmt.Fprintf(c.Writer, "{\n")
  expvar.Do(func(kv expvar.KeyValue) {
    report(kv.Key, kv.Value)
  })
  fmt.Fprintf(c.Writer, "\n}\n")

  c.String(http.StatusOK, "")
}

func init() {
  expvar.Publish("upTime", expvar.Func(calculateUptime))
  expvar.Publish("version", expvar.Func(currentGoVersion))
  expvar.Publish("cores", expvar.Func(getNumCPUs))
  expvar.Publish("os", expvar.Func(getGoOS))
  expvar.Publish("cgo", expvar.Func(getNumCgoCall))
  expvar.Publish("goroutine", expvar.Func(getNumGoroutins))
  expvar.Publish("gcpause", expvar.Func(getLastGCPauseTime))
}
