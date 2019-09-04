package main

import (
	"fmt"
	"genosha/controllers"
	"genosha/utils/confs"
	"genosha/utils/monitor"
	"genosha/utils/myLogger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"time"
)

var ginFile *rotatelogs.RotateLogs

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
}

func main() {
	maxAgeTime, err := strconv.Atoi(confs.ConfigMap["MaxAge"])
	if err != nil {
		myLogger.Log.Error("covert maxAgeTime string to num fail")
		return
	}
	l := gin.LoggerConfig{
		Output: ginFile,
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf(" %s |%s | %s | %s | %s | %s | %d | %s | %s | [%s] \n",
				"[GIN]",
				params.TimeStamp.Format(time.RFC3339),
				params.ClientIP,
				params.Method,
				params.Path,
				params.Request.Proto,
				params.StatusCode,
				params.Latency,
				params.Request.UserAgent(),
				params.ErrorMessage,
			)
		},
	}
	router := gin.New()
	router.Use(gin.LoggerWithConfig(l))
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           time.Duration(maxAgeTime) * time.Hour,
	}))

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{"genosha": "genosha"}))
	authorized.GET("/ping", ping)
	router.GET("/genoshaStatus", monitor.GetCurrentRunningStats)

	usersController := controllers.NewUsersController()
	forwardController := controllers.NewForwardController()

	router.POST("/signup", usersController.Signup)
	router.POST("/login", usersController.Login)
	router.POST("/users/send_passwordreset_email", usersController.SendResetPassWordEmail)
	router.POST("/users/reset_password", usersController.ResetPassWord)

	userGroup := router.Group("/users").Use(usersController.UserAgentAuth(), controllers.Auth.MiddlewareFunc())
	{
		userGroup.GET("/refreshToken", usersController.RefreshToken)
		userGroup.GET("/userInfo", usersController.GetUserInfo)
		userGroup.POST("/passWord", usersController.ChangePassWord)
	}

	forwardGroup := router.Group("/forward").Use(usersController.UserAgentAuth(), controllers.Auth.MiddlewareFunc())
	{

		forwardGroup.GET("/devices/:id", forwardController.GetSimpleForward)
		forwardGroup.POST("/devices/:id", forwardController.PostSimpleForward)
		forwardGroup.DELETE("/devices/:id", forwardController.DeleteSimpleForward)
	}

	myLogger.Log.Info("====== server to listen")
	err = router.Run(confs.FlagSericePort)
	if err != nil {
		myLogger.Log.Info("====== server fail to run")
	}
	myLogger.Log.Info("====== server to listen done")
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
