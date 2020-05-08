package controllers

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"genosha/dao"
	"genosha/models"
	"genosha/utils/myLogger"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

func parseUserGuidFromRequest(c *gin.Context) string {
	return jwt.ExtractClaims(c)["guid"].(string)
}

func keySign(c *gin.Context, method string, body []byte, sUrl string) (string, error) {
	var content string
	userGuid := parseUserGuidFromRequest(c)
	userInfo := dao.GetUserByGuid(userGuid)
	myLogger.Log.Info(string(body))
	if body == nil || string(body) == "" {
		content = ""
	} else {
		has := md5.Sum(body)
		content = fmt.Sprintf("%x", has)
	}
	sign := method + content + sUrl
	myLogger.Log.Info("sign", zap.Any("sign", sign))
	secret := []byte(userInfo.Name)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(sign))
	sha := hex.EncodeToString(h.Sum(nil))
	encodeSign := base64.StdEncoding.EncodeToString([]byte(sha))
	myLogger.Log.Info("encodeSign", zap.Any("encodeSign", encodeSign))
	params := url.Values{}
	params.Add("token", userInfo.Token)
	params.Add("sign", encodeSign)
	return params.Encode(), nil
}

func handleResp(resp *http.Response, c *gin.Context, body []byte, header http.Header) {
	if resp.StatusCode == 400 {
		var BackendError models.BackendError
		err := json.Unmarshal(body, &BackendError)
		if err != nil {
			myLogger.Log.Error("error", zap.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":  "internal json err",
				"code": 5000,
			})
			return
		}
		switch BackendError.Code {
		default:
			myLogger.Log.Error("bad request")
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":      BackendError.Msg,
				"backcode": BackendError.Code,
				"code":     4100,
			})
			return
		}
	}
	if resp.StatusCode == 401 {
		myLogger.Log.Error("signature failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "internal err401",
			"code": 5000,
		})
		return
	}
	if resp.StatusCode == 403 {
		myLogger.Log.Error("authentication failed")
		c.JSON(http.StatusForbidden, gin.H{
			"msg":  "authentication failed",
			"code": 4103,
		})
		return
	}
	if resp.StatusCode == 404 {
		myLogger.Log.Error("Not Found")
		c.JSON(http.StatusForbidden, gin.H{
			"msg":  "Not Found",
			"code": 4104,
		})
		return
	}
	if resp.StatusCode == 500 {
		myLogger.Log.Error("authentication failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "internal err500",
			"code": 5000,
		})
		return
	}
	for _, v := range header["Set-Cookie"] {
		c.Writer.Header().Add("Set-Cookie", v)
	}
	c.Writer.Write(body)
}

