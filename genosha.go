package main

import (
	"genosha/collyUnit"
	"genosha/utils/myLogger"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var logFile *rotatelogs.RotateLogs

//var router  *gin.Engine

func init() {
	logFile, _ = rotatelogs.New(
		"./log/log.%Y-%m-%d-%H-%M",
		rotatelogs.WithMaxAge(30*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	gin.DisableConsoleColor()
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	myLogger.MyLogInit(logFile)
	//router = ginUnit.RouterInit(logFile)
}

func main() {
	myLogger.Log.Info("====== server to listen")
	//err := router.Run(confs.FlagSericePort)
	//if err != nil {
	//	myLogger.Log.Info("====== server fail to run")
	//}
	//collyUnit.CollyInit()
	//crawler.Douban250()
	collyUnit.CollyRun()
	myLogger.Log.Info("====== server to listen done")
}
