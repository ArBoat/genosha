package main

import (
	"genosha/ginUnit"
	"genosha/utils/confs"
	"genosha/utils/myLogger"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var ginFile *rotatelogs.RotateLogs
var router  *gin.Engine

func init() {
	ginFile, _ = rotatelogs.New(
		"./log/log.%Y-%m-%d-%H-%M",
		rotatelogs.WithMaxAge(30*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	gin.DisableConsoleColor()
	log.SetOutput(ginFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	myLogger.MyLogInit(ginFile)
	router = ginUnit.RouterInit(ginFile)
}

func main() {
	myLogger.Log.Info("====== server to listen")
	err := router.Run(confs.FlagSericePort)
	if err != nil {
		myLogger.Log.Info("====== server fail to run")
	}
	myLogger.Log.Info("====== server to listen done")
}
