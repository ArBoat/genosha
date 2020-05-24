package controllers

import (
	"fmt"
	"genosha/dao"
	"genosha/models"
	"genosha/utils/confs"
	"genosha/utils/myLogger"
	"genosha/utils/tools"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type UsersController struct{}

func NewUsersController() *UsersController {
	return &UsersController{}
}

var usersController = NewUsersController()
var acl = casbin.NewEnforcer("./auth/authz_model.conf", "./auth/authz_policy.csv")
var timeOutTime, _ = strconv.Atoi(confs.ConfigMap["Timeout"])
var timeOut = time.Duration(timeOutTime) * time.Hour
var timeOutTimeMobile, _ = strconv.Atoi(confs.ConfigMap["TimeoutMobile"])
var timeOutMobile = time.Duration(timeOutTimeMobile) * time.Hour
var maxRefreshTime, _ = strconv.Atoi(confs.ConfigMap["MaxRefresh"])
var Auth = &jwt.GinJWTMiddleware{
	Realm:      "genosha zone",
	Key:        []byte(confs.ConfigMap["jwtKey"]),
	Timeout:    timeOut,
	MaxRefresh: time.Duration(maxRefreshTime) * time.Hour,
	PayloadFunc: func(data interface{}) jwt.MapClaims {
		guid := usersController.GetUserGUID(data.(string))
		userInfo := dao.GetUserByGuid(guid)
		return jwt.MapClaims{"version": userInfo.Version, "guid": guid}
	},
	LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(code, gin.H{
			"code":   code,
			"expire": (expire.UnixNano()) / 1e6,
			"token":  token,
		})
	},
	RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(code, gin.H{
			"code":   code,
			"expire": (expire.UnixNano()) / 1e6,
			"token":  token,
		})
	},
	// user validatoin control
	Authenticator: func(c *gin.Context) (interface{}, error) {
		return nil, nil
	},
	// permission control
	Authorizator: func(userId interface{}, c *gin.Context) bool {
		myLogger.Log.Info("in the authorizator..")
		guidClaim := jwt.ExtractClaims(c)["guid"]
		if guidClaim == nil {
			myLogger.Log.Error("guidClaim nil")
			return false
		}
		guid := guidClaim.(string)
		version := jwt.ExtractClaims(c)["version"].(string)
		userInfo := dao.GetUserByGuid(guid)
		if userInfo == nil {
			myLogger.Log.Warn("wrong guid")
			return false
		}
		if version != userInfo.Version {
			myLogger.Log.Warn("wrong version")
			return false
		}
		roles := usersController.GetUserRoles(guid)
		for _, role := range roles {
			// use casbin's ACL to check user's permission
			result, err := acl.EnforceSafe(role, c.Request.URL.Path, c.Request.Method)
			if err == nil && result {
				myLogger.Log.Info("acl validation done!")
				return true
			}
			myLogger.Log.Error("error", zap.Error(err))
		}
		myLogger.Log.Warn("failed to pass acl..")
		return false
	},
	Unauthorized: func(c *gin.Context, code int, message string) {
		c.JSON(401, gin.H{
			"code": 4101,
			"msg":  message,
		})
		myLogger.Log.Warn("acl validation fail!")
	},
	TokenLookup:   "header:Authorization",
	TokenHeadName: "Bearer",
	TimeFunc:      time.Now,
}

func (uc UsersController) Login(c *gin.Context) {
	userAgent := userAgent(c)
	userAgent.LoginHandler(c)
}
func (uc UsersController) RefreshToken(c *gin.Context) {
	userAgent := userAgent(c)
	userAgent.RefreshHandler(c)
}
func (uc UsersController) UserAgentAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		Auth = userAgent(c)
	}
}
func userAgent(c *gin.Context) *jwt.GinJWTMiddleware {
	useragent := c.Request.UserAgent()
	myLogger.Log.Info("useragent" + useragent)
	useragentSlice := strings.Split(useragent, ";")
	if len(useragentSlice) >= 3 {
		match, _ := regexp.MatchString("App/[^/]+/[^/]+", useragentSlice[0])
		iosIndex := strings.Index(useragentSlice[1], "iOS")
		androidIndex := strings.Index(useragentSlice[1], "Android")
		if match == true && (iosIndex != -1 || androidIndex != -1) {
			Auth.Timeout = timeOutMobile
		} else {
			Auth.Timeout = timeOut
		}
	}
	myLogger.Log.Info("Auth", zap.Any("Auth", Auth.Timeout))
	return Auth
}

func (uc UsersController) Signup(c *gin.Context) {
	var signup models.Signup
	if err := c.ShouldBindJSON(&signup); err != nil {
		myLogger.Log.Error("error", zap.Error(err))
		return
	}
	myLogger.Log.Info("requestBody", zap.Any("request", signup))
	if signup.Spell != "Because I am batman!" {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "you should tell me the spell"})
		myLogger.Log.Warn("wrong spell", zap.String("spell", signup.Spell))
		return
	}
	randomPassword := tools.GetRandomString(8)
	hashPassword, hashErr := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
	if hashErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "encrypt failed"})
		myLogger.Log.Error("error", zap.Error(hashErr))
		return
	}
	guid := tools.GetGuid()
	if err := dao.CreateNewUser(&models.User{
		Guid:         guid,
		Name:         signup.UserName,
		Email:        signup.Email,
		LowEmail:     strings.ToLower(signup.Email),
		PasswordHash: string(hashPassword),
		Version:      tools.GetRandomString(8),
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		myLogger.Log.Error("error", zap.Error(err))
		return
	}
	roles := strings.Split(signup.Roles, ",")
	for i := 0; i < len(roles); i++ {
		if err := dao.CreateUserToRole(&models.UserToRole{UserGuid: guid, Role: roles[i]}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			myLogger.Log.Error("error", zap.Error(err))
			return
		}
	}
	sendSignUpEmail2User(signup.Email, randomPassword)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (uc UsersController) ValidateUser(userEmail string, password string) bool {
	if existedUser := dao.GetUserByLowEmail(strings.ToLower(userEmail)); existedUser != nil {
		myLogger.Log.Info("existed user:" + existedUser.Email)
		if cmpErr := bcrypt.CompareHashAndPassword([]byte(existedUser.PasswordHash), []byte(password)); cmpErr == nil {
			return true
		}
	}
	return false
}

func (uc UsersController) GetUserInfo(c *gin.Context) {
	userGuid := parseUserGuidFromRequest(c)
	userInfo := dao.GetUserByGuid(userGuid)
	c.JSON(http.StatusOK, gin.H{
		"userName": userInfo.Name,
		"email":    userInfo.Email,
		"roles":    dao.GetUserRolesByGuid(userGuid),
	})
}

func (uc UsersController) ChangePassWord(c *gin.Context) {
	userGuid := parseUserGuidFromRequest(c)
	var changePassWD models.ChangePassWD
	if err := c.ShouldBindJSON(&changePassWD); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "json format err",
			"code": 4000,
		})
		myLogger.Log.Error("error", zap.Error(err))
		return
	}
	myLogger.Log.Info("requestBody", zap.Any("request", changePassWD))
	existedUser := dao.GetUserByGuid(userGuid)
	if existedUser == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "the user does not exist",
			"code": 4001,
		})
		myLogger.Log.Error("the user does not exist")
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(existedUser.PasswordHash), []byte(changePassWD.OldPassword))
	if err != nil {
		myLogger.Log.Error("error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "wrong old password",
			"code": 4002,
		})
		return
	}
	newPassWDHash, err := bcrypt.GenerateFromPassword([]byte(changePassWD.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		myLogger.Log.Error("error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "encrypt newPassWDHash failed",
			"code": 5000,
		})
		return
	}
	err = dao.UpdatePassWDByEmail(existedUser.Email, string(newPassWDHash))
	if err != nil {
		myLogger.Log.Error("error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "failed to update new password",
			"code": 5000,
		})
		return
	}
	err = dao.UpdateVersionByEmail(existedUser.Email, tools.GetRandomString(8))
	if err != nil {
		myLogger.Log.Error("error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "failed to update new version",
			"code": 5000,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": true})
}

func (uc UsersController) SendResetPassWordEmail(c *gin.Context) {
	var requestBody models.SendRestPassWDEmail
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "json format err",
			"code": 4000,
		})
		myLogger.Log.Error("error", zap.Error(err))
		return
	}
	existedUser := dao.GetUserByLowEmail(strings.ToLower(requestBody.To))
	if existedUser == nil || existedUser.Guid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "the user does not exist",
			"code": 4001,
		})
		myLogger.Log.Error("error", zap.Error(err))
		return
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	updateTime1min := existedUser.UpdatedAt.Add(1 * time.Minute).Format("2006-01-02 15:04:05")
	updateTime1day := existedUser.UpdatedAt.Add(24 * time.Hour).Format("2006-01-02 15:04:05")

	if currentTime < updateTime1min {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":       "not reach the min retry interval",
			"code":      4003,
			"extraInfo": "1",
		})
		myLogger.Log.Error("error:" + "not reach the min retry interval")
		return
	}
	if currentTime > updateTime1day {
		err = dao.ResetTokenCountByLowEmail(strings.ToLower(requestBody.To))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":  "fail to reset tokenCount",
				"code": 5000,
			})
			myLogger.Log.Error("error", zap.Error(err))
			return
		}
	}
	if existedUser.TokenCount >= 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":       "retry count exceeds the limit",
			"code":      4004,
			"extraInfo": "10",
		})
		myLogger.Log.Error("error:" + "retry count exceeds the limit")
		return
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	token := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	myLogger.Log.Info("reset", zap.Any("token", token))
	err = dao.UpdateTokenByLowEmail(strings.ToLower(requestBody.To), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "fail to update token",
			"code": 5000,
		})
		myLogger.Log.Error("error", zap.Error(err))
		return
	}
	err = dao.IncreaseTokenCountByLowEmail(strings.ToLower(requestBody.To))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "fail to increase tokenCount",
			"code": 5000,
		})
		myLogger.Log.Error("error", zap.Error(err))
		return
	}
	err = sendResetToken2Email(requestBody.To, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "fail to send email",
			"code": 5000,
		})
		myLogger.Log.Error("error", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": true})
}

func (uc UsersController) ResetPassWord(c *gin.Context) {
	var requestBody models.RestPassWD
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "json format err",
			"code": 4000,
		})
		myLogger.Log.Error("error", zap.Error(err))
		return
	}
	existedUser := dao.GetUserByLowEmail(strings.ToLower(requestBody.Email))
	if existedUser == nil || existedUser.Guid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "the user does not exist",
			"code": 4001,
		})
		myLogger.Log.Error("error", zap.Error(err))
		return
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	updateTime30min := existedUser.UpdatedAt.Add(30 * time.Minute).Format("2006-01-02 15:04:05")
	if currentTime > updateTime30min {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "verification token is expired",
			"code": 4005,
		})
		myLogger.Log.Error("error:" + "verification token is expired")
		return
	}
	if existedUser.Token != requestBody.Token {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "verification token is incorrect",
			"code": 4006,
		})
		myLogger.Log.Error("error:" + "verification token is incorrect")
		return
	}
	newPassWDHash, err := bcrypt.GenerateFromPassword([]byte(requestBody.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		myLogger.Log.Error("error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "encrypt newPassWDHash failed",
			"code": 5000,
		})
		return
	}
	err = dao.UpdatePassWDByEmail(existedUser.Email, string(newPassWDHash))
	if err != nil {
		myLogger.Log.Error("error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "failed to update new password",
			"code": 5000,
		})
		return
	}
	err = dao.UpdateVersionByEmail(existedUser.Email, tools.GetRandomString(8))
	if err != nil {
		myLogger.Log.Error("error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "failed to update new version",
			"code": 5000,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": true})
}

func (uc UsersController) GetUserRoles(guid string) []string {
	return dao.GetUserRolesByGuid(guid)
}

func (uc UsersController) GetUserGUID(userEmail string) string {
	return dao.GetUserByLowEmail(strings.ToLower(userEmail)).Guid
}

func sendSignUpEmail2User(email, password string) {
	from := mail.NewEmail(confs.ConfigMap["senderName"], confs.ConfigMap["senderAddress"])
	to := mail.NewEmail(email, email)
	subject := "genosha Console"
	plainTextContent := "-"
	htmlContent := `<p><b>Welcome to genosha console !</b></p>` +
		`<p>---------------------</p>` +
		`<p><b>Login Email: ` + email + `</b></p>` +
		`<p><b>Password: ` + password + `</b></p>` +
		`<p>---------------------</p>` +
		confs.ConfigMap["consoleType"]
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(confs.ConfigMap["SENDGRID_API_KEY"])
	response, err := client.Send(message)
	if err != nil {
		myLogger.Log.Error("error", zap.Error(err))
		return
	}
	myLogger.Log.Info("response.StatusCode", zap.Any("response.StatusCode", response.StatusCode))
	myLogger.Log.Info("response.Body", zap.Any("response.Body", response.Body))
}

func sendResetToken2Email(email, token string) error {
	from := mail.NewEmail(confs.ConfigMap["senderName"], confs.ConfigMap["senderAddress"])
	to := mail.NewEmail(email, email)
	subject := "Reset genosha Console Password"
	plainTextContent := "-"
	htmlContent := `<p><b>Reset genosha Console Password Token!</b></p>` +
		`<p>---------------------</p>` +
		`<p><b>Login Email: ` + email + `</b></p>` +
		`<p><b>Token: ` + token + `</b></p>` +
		`<p>---------------------</p>` +
		confs.ConfigMap["consoleType"]
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(confs.ConfigMap["SENDGRID_API_KEY"])
	response, err := client.Send(message)
	if err != nil {
		myLogger.Log.Error("error", zap.Error(err))
		return err
	}
	myLogger.Log.Info("response.StatusCode", zap.Any("response.StatusCode", response.StatusCode))
	myLogger.Log.Info("response.Body", zap.Any("response.Body", response.Body))
	return nil
}
