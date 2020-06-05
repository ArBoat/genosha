package ginUnit

import (
  "genosha/ginUnit/controllers"
  "github.com/gin-gonic/gin"
  "net/http"
)

//func RouterInit(path io.Writer) *gin.Engine {
//	maxAgeTime, err := strconv.Atoi(confs.ConfigMap["MaxAge"])
//	if err != nil {
//		myLogger.Log.Error("covert maxAgeTime string to num fail")
//	}
//	l := gin.LoggerConfig{
//		Output: path,
//		Formatter: func(params gin.LogFormatterParams) string {
//			return fmt.Sprintf(" %s |%s | %s | %s | %s | %s | %d | %s | %s | [%s] \n",
//				"[GIN]",
//				params.TimeStamp.Format(time.RFC3339),
//				params.ClientIP,
//				params.Method,
//				params.Path,
//				params.Request.Proto,
//				params.StatusCode,
//				params.Latency,
//				params.Request.UserAgent(),
//				params.ErrorMessage,
//			)
//		},
//	}
//	gin.SetMode(gin.DebugMode)
//	router := gin.New()
//	router.Use(gin.LoggerWithConfig(l))
//	router.Use(gin.Recovery())
//	router.Use(cors.New(cors.Config{
//		AllowAllOrigins:  true,
//		AllowMethods:     []string{"PUT", "POST", "PATCH", "DELETE"},
//		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
//		ExposeHeaders:    []string{"Content-Length"},
//		AllowCredentials: true,
//		MaxAge:           time.Duration(maxAgeTime) * time.Hour,
//	}))
//
//	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{"genosha": "genosha"}))
//	authorized.GET("/ping", ping)
//	router.GET("/genoshaStatus", monitor.GetCurrentRunningStats)
//
//	usersController := controllers.NewUsersController()
//	forwardController := controllers.NewForwardController()
//
//	router.POST("/signup", usersController.Signup)
//	router.POST("/login", usersController.Login)
//	router.POST("/users/send_passwordreset_email", usersController.SendResetPassWordEmail)
//	router.POST("/users/reset_password", usersController.ResetPassWord)
//
//	userGroup := router.Group("/users").Use(usersController.UserAgentAuth(), controllers.Auth.MiddlewareFunc())
//	{
//		userGroup.GET("/refreshToken", usersController.RefreshToken)
//		userGroup.GET("/userInfo", usersController.GetUserInfo)
//		userGroup.POST("/passWord", usersController.ChangePassWord)
//	}
//
//	forwardGroup := router.Group("/forward").Use(usersController.UserAgentAuth(), controllers.Auth.MiddlewareFunc())
//	{
//
//		forwardGroup.GET("/devices/:id", forwardController.GetSimpleForward)
//		forwardGroup.POST("/devices/:id", forwardController.PostSimpleForward)
//		forwardGroup.DELETE("/devices/:id", forwardController.DeleteSimpleForward)
//	}
//	return router
//}
func ping(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{
    "result": true,
  })
}

func setupRouter() *gin.Engine {
  r := gin.Default()
  r.GET("/ping", ping)
  forwardController := controllers.NewForwardController()
  forwardGroup := r.Group("/forward")
  {

    forwardGroup.GET("/devices/:id", forwardController.GetSimpleForward)
    forwardGroup.POST("/devices/:id", forwardController.PostSimpleForward)
    forwardGroup.DELETE("/devices/:id", forwardController.DeleteSimpleForward)
  }
  return r
}
